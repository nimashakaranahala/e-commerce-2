package ports

import "e-commerce/internal/models"

type Repository interface {
	FindUserByEmail(email string) (*models.User, error)
	GetUserByID(userID uint) (*models.User, error)
	FindAllUsers() ([]models.User, error)
	FindSellerByEmail(email string) (*models.Seller, error)
	CreateUser(user *models.User) error
	CreateSeller(Seller *models.Seller) error
	UpdateUser(user *models.User) error
	UpdateSeller(user *models.Seller) error
	BlacklistToken(token *models.BlacklistTokens) error
	TokenInBlacklist(token *string) bool
	CreateProduct(product *models.Product) error
	GetProductByID(productID uint) (*models.Product, error)
	GetAllProducts() ([]models.Product, error)

	// Add the AddToCart method
	AddToCart(cart *models.Cart) error

	// Add the AddToCart method
    ViewCart(cart *models.Cart) error
}
