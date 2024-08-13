package server

import (
	"e-commerce/internal/api"
	"e-commerce/internal/middleware"
	"e-commerce/internal/ports"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// SetupRouter is where router endpoints are called
func SetupRouter(handler *api.HTTPHandler, repository ports.Repository) *gin.Engine {
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"POST", "GET", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r := router.Group("/")
	{
		r.GET("/", handler.Readiness)
	}

	user := r.Group("/user")
	{
		user.POST("/create", handler.CreateUser)
		user.POST("/login", handler.LoginUser)
	}

	// AuthorizeUser authorizes all the authorized users handlers
	user.Use(middleware.AuthorizeUser(repository.FindUserByEmail, repository.TokenInBlacklist))
	{
		user.POST("/logout", handler.Logout)
		user.GET("/allproducts", handler.GetAllProducts)
		user.GET("/product/:id", handler.GetProductByID)
	}

	seller := r.Group("/seller")
	{
		seller.POST("/create", handler.CreateSeller)
		seller.POST("/login", handler.LoginSeller)
	}

	// AuthorizeSeller authorizes all the authorized sellers hanldlers
	seller.Use(middleware.AuthorizeSeller(repository.FindSellerByEmail, repository.TokenInBlacklist))
	{
		seller.POST("/logout", handler.Logout)
		seller.POST("/createproduct", handler.CreateProduct)
	}

	return router
}
