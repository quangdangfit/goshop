package models

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"

	"goshop/dbs"
)

type OrderLine struct {
	Base
	ProductID string `json:"product_id"`
	OrderID   string `json:"order_id"`
	Quantity  uint   `json:"quantity"`
	Price     uint   `json:"price"`
}

func (line *OrderLine) BeforeCreate(scope *gorm.Scope) error {
	line.ID = uuid.New().String()
	var product Product
	dbs.Database.Where("uuid = ?", line.ProductID).First(&product)
	line.Price = product.Price * line.Quantity

	return nil
}
