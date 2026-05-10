package model

import (
	"time"
)

type Product struct {
	ID          string     `json:"id" gorm:"unique;not null;index;primary_key"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at" gorm:"index"`
	Code        string     `json:"code" gorm:"uniqueIndex:idx_product_code,not null"`
	Name        string     `json:"name" gorm:"uniqueIndex:idx_product_name,not null"`
	Description string     `json:"description"`
	Price       float64    `json:"price"`
	Active      bool       `json:"active" gorm:"default:true"`
	// StockQuantity / ReservedQuantity mirror the canonical products columns so the order
	// service can compute available stock without crossing into the product domain.
	StockQuantity    int `json:"stock_quantity" gorm:"default:0"`
	ReservedQuantity int `json:"reserved_quantity" gorm:"default:0"`
}
