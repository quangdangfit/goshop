package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderLine struct {
	Base
	OrderID   string `json:"order_id"`
	ProductID string `json:"product_id"`
	Product   *Product
	Quantity  uint    `json:"quantity"`
	Price     float64 `json:"price"`
}

func (line *OrderLine) BeforeCreate(tx *gorm.DB) error {
	line.ID = uuid.New().String()
	return nil
}
