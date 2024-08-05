package models

import (
	"gorm.io/gorm"
)

type Seller struct {
	gorm.Model
	FirstName     string `json:"first_name" binding:"required"`
	LastName      string `json:"last_name" binding:"required"`
	Password      string `json:"password" binding:"required"`
	DateOfBirth   string `json:"date_of_birth" binding:"required"`
	Email         string `json:"email" binding:"required,email"`
	Phone         string `json:"phone" binding:"required"`
	Address       string `json:"address"`
	StoreName     string `json:"store_name"`
	StoreCategory string `json:"store_category"`
}

type LoginRequestSeller struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}
