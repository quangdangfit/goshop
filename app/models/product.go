package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"goshop/pkg/utils"
)

type Product struct {
	Base
	Code        string  `json:"code" gorm:"uniqueIndex:idx_product_code,not null"`
	Name        string  `json:"name" gorm:"uniqueIndex:idx_product_name,not null"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Active      bool    `json:"active" gorm:"default:true"`
}

func (product *Product) BeforeCreate(tx *gorm.DB) error {
	product.ID = uuid.New().String()
	product.Code = utils.GenerateCode("P")
	product.Active = true
	return nil
}
