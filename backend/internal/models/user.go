package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	AccountNumber uint   `json:"account_number" gorm:"unique;not null"`
	FirstName     string `json:"first_name" binding:"required"`
	Lastname      string `json:"last_name" binding:"required"`
	Password      string `json:"password" binding:"required"`
	DateOfBirth   string `json:"date_of_birth" binding:"required"`
	Email         string `json:"email" binding:"required,email"`
	Phone         string `json:"phone" binding:"required"`
	Address       string `json:"address" binding:"required"`
}

type LoginRequestUser struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}
