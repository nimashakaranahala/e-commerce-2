package models

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	SellerID    uint    `json:"seller_id"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	ImageUrl    string  `json:"image_url"`
	Quantity    int     `json:"quantity"`
	Description string  `json:"description"`
	Status      bool    `json:"status"`
}
