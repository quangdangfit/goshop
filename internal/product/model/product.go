package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"goshop/pkg/utils"
)

type Product struct {
	ID            string     `json:"id" gorm:"unique;not null;index;primary_key"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	DeletedAt     *time.Time `json:"deleted_at" gorm:"index"`
	Code          string     `json:"code" gorm:"uniqueIndex:idx_product_code,not null"`
	Name          string     `json:"name" gorm:"uniqueIndex:idx_product_name,not null"`
	Description   string     `json:"description"`
	Price         float64    `json:"price"`
	Active        bool       `json:"active" gorm:"default:true"`
	StockQuantity int        `json:"stock_quantity" gorm:"default:0"`
	AvgRating     float64    `json:"avg_rating" gorm:"default:0"`
	ReviewCount   int        `json:"review_count" gorm:"default:0"`
	Images        []string   `json:"images" gorm:"serializer:json"`
	CategoryID    *string    `json:"category_id"`
	Category      *Category  `json:"category,omitempty"`
}

func (m *Product) BeforeCreate(tx *gorm.DB) error {
	m.ID = uuid.New().String()
	m.Code = utils.GenerateCode("P")
	m.Active = true
	return nil
}
