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
}
