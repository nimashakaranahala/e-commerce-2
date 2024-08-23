package models

import "gorm.io/gorm"

type Cart struct {
	gorm.Model
	SellerID    uint    `json:"seller_id"`
	ProductID 	uint 	`json:"product_id"`
	Quantity    int     `json:"quantity"`
}
