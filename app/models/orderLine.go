package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"goshop/dbs"
)

type OrderLine struct {
	Base
	ProductID string  `json:"product_id"`
	OrderID   string  `json:"order_id"`
	Quantity  uint    `json:"quantity"`
	Price     float64 `json:"price"`
}

func (line *OrderLine) BeforeCreate(tx *gorm.DB) error {
	line.ID = uuid.New().String()
	var product Product
	dbs.Database.Where("uuid = ?", line.ProductID).First(&product)
	line.Price = product.Price * float64(line.Quantity)

	return nil
}
