package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderLine struct {
	ID        string     `json:"id" gorm:"unique;not null;index;primary_key"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" gorm:"index"`
	OrderID   string     `json:"order_id"`
	ProductID string     `json:"product_id"`
	Product   *Product
	Quantity  uint    `json:"quantity"`
	Price     float64 `json:"price"`
}

func (line *OrderLine) BeforeCreate(tx *gorm.DB) error {
	line.ID = uuid.New().String()
	return nil
}
