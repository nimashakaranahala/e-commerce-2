package api

import (
	"e-commerce/internal/middleware"
	"e-commerce/internal/models"
	"e-commerce/internal/util"
	"os"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// Create User
func (u *HTTPHandler) CreateUser(c *gin.Context) {
	var user *models.User
	if err := c.ShouldBind(&user); err != nil {
		util.Response(c, "invalid request", 400, err.Error(), nil)
		return
	}

	// Hash the password
	hashedPassword, err := util.HashPassword(user.Password)
	if err != nil {
		util.Response(c, "Internal server error", 500, err.Error(), nil)
		return
	}
	user.Password = hashedPassword

	err = u.Repository.CreateUser(user)
	if err != nil {
		util.Response(c, "User not created", 500, err.Error(), nil)
		return
	}
	util.Response(c, "User created", 200, nil, nil)

}

// Login User
func (u *HTTPHandler) LoginUser(c *gin.Context) {
	var loginRequest *models.LoginRequestUser
	err := c.ShouldBind(&loginRequest)
	if err != nil {
		util.Response(c, "invalid request", 400, err.Error(), nil)
		return
	}

	if loginRequest.Email == "" || loginRequest.Password == "" {
		util.Response(c, "Email and Password must not be empty", 400, nil, nil)
		return
	}

	user, err := u.Repository.FindUserByEmail(loginRequest.Email)
	if err != nil {
		util.Response(c, "Email does not exist", 404, err.Error(), nil)
		return
	}

	// Verify the password
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password)); err != nil {
		util.Response(c, "invalid email or password", 400, "invalid email or password", nil)
		return
	}

	accessClaims, refreshClaims := middleware.GenerateClaims(user.Email)

	secret := os.Getenv("JWT_SECRET")

	accessToken, err := middleware.GenerateToken(jwt.SigningMethodHS256, accessClaims, &secret)
	if err != nil {
		util.Response(c, "Error generating access token", 500, err.Error(), nil)
		return
	}

	refreshToken, err := middleware.GenerateToken(jwt.SigningMethodHS256, refreshClaims, &secret)
	if err != nil {
		util.Response(c, "Error generating refresh token", 500, err.Error(), nil)
		return
	}

	c.Header("access_token", *accessToken)
	c.Header("refresh_token", *refreshToken)

	util.Response(c, "Login successful", 200, gin.H{
		"user":          user,
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}, nil)
}

// View Product listing
func (u *HTTPHandler) GetAllProducts(c *gin.Context) {
	_, err := u.GetUserFromContext(c)
	if err != nil {
		util.Response(c, "invalid token", 401, err.Error(), nil)
		return
	}

	products, err := u.Repository.GetAllProducts()
	if err != nil {
		util.Response(c, "Error getting products", 500, err.Error(), nil)
		return
	}
	util.Response(c, "Success", 200, products, nil)
}

// View Product by ID
func (u *HTTPHandler) GetProductByID(c *gin.Context) {
	_, err := u.GetUserFromContext(c)
	if err != nil {
		util.Response(c, "invalid token", 401, err.Error(), nil)
		return
	}

	productID := c.Param("id")
	id, err := strconv.Atoi(productID)
	if err != nil {
		util.Response(c, "Error getting product", 500, err.Error(), nil)
		return
	}
	product, err := u.Repository.GetProductByID(uint(id))
	if err != nil {
		util.Response(c, "Error getting product", 500, err.Error(), nil)
		return
	}
	util.Response(c, "Success", 200, product, nil)
}

// AddToCart adds a product to the user's cart
func (u *HTTPHandler) AddToCart(c *gin.Context) {
	user, err := u.GetUserFromContext(c)
	if err != nil {
		util.Response(c, "invalid token", 401, err.Error(), nil)
		return
	}

	var requestBody struct {
		ProductID uint `json:"product_id"`
		Quantity  int  `json:"quantity"`
	}

	// Bind the request body
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		util.Response(c, "Invalid request", 401, err.Error(), nil)
		return
	}

	// Check if the product exists
	product, err := u.Repository.GetProductByID(requestBody.ProductID)
	if err != nil {
		util.Response(c, "Product not found", 500, err.Error(), nil)
		return
	}

	if product.Quantity < requestBody.Quantity {
		util.Response(c, "Quantity requested is more than available", 400, nil, nil)
		return
	}

	// Create a new cart item
	cart := &models.Cart{
		UserID:    user.ID,
		ProductID: product.ID,
		Quantity:  requestBody.Quantity,
	}

	// Add the item to the cart
	err = u.Repository.AddToCart(cart)
	if err != nil {
		util.Response(c, "Could not add to cart", 500, err.Error(), nil)
		return
	}

	util.Response(c, "Product added to cart", 200, cart, nil)
}

func (u *HTTPHandler) ViewCart(c *gin.Context) {
	user, err := u.GetUserFromContext(c)
	if err != nil {
		util.Response(c, "Invalid token", 401, err.Error(), nil)
		return
	}

	cartItems, err := u.Repository.GetCartsByUserID(user.ID)
	if err != nil {
		util.Response(c, "Error retrieving cart", 500, err.Error(), nil)
		return
	}

	util.Response(c, "Cart retrieved", 200, cartItems, nil)
}

//RemoveFromCart
func (u *HTTPHandler) RemoveFromCart(c *gin.Context) {
    // Get the user from the context
    user, err := u.GetUserFromContext(c)
    if err != nil {
        util.Response(c, "Invalid token", 401, err.Error(), nil)
        return
    }

	productID := c.Param("id")
	productIDInt, err := strconv.Atoi(productID)
	if err != nil {
		util.Response(c, "Invalid product ID", 400, err.Error(), nil)
		return
	}

	// Validate if product exist in cart
	_, err = u.Repository.GetCartItemByProductID(uint(productIDInt))
	if err != nil {
		util.Response(c, "Product not found in cart", 404, err.Error(), nil)
		return
	}
  
    err = u.Repository.RemoveItemFromCart(user.ID, uint(productIDInt))
    if err != nil {
        util.Response(c, "Could not remove item from cart", 500, err.Error(), nil)
        return
    }

    util.Response(c, "Item removed from cart", 200, nil, nil)
}

// edit cart
func (u *HTTPHandler) EditCart(c *gin.Context) {
	// Get user id from context
	user, err := u.GetUserFromContext(c)
	if err != nil {
		util.Response(c, "Error getting user from context", 500, err.Error(), nil)
		return
	}

	// Bind request to struct
	var cart *models.Cart
	if err := c.ShouldBind(&cart); err != nil {
		util.Response(c, "invalid request", 400, err.Error(), nil)
		return
	}

	// Get cart by user id
	shoppingCart, err := u.Repository.GetCartItemByProductID(cart.ProductID)
	if err != nil {
		util.Response(c, "Cart not found", 404, err.Error(), nil)
		return
	}

	// Validate request
	product, err := u.Repository.GetProductByID(cart.ProductID)
	if err != nil {
		util.Response(c, "Product not found", 404, err.Error(), nil)
		return
	}

	// Check if product quantity is less
	if product.Quantity < cart.Quantity {
		util.Response(c, "Product quantity is less", 400, nil, nil)
		return
	}

	// Update cart
	cart.UserID = user.ID
	cart.ID = shoppingCart.ID

	err = u.Repository.AddToCart(cart)
	if err != nil {
		util.Response(c, "Internal server error", 500, err.Error(), nil)
		return
	}
	util.Response(c, "Cart updated", 200, nil, nil)
}