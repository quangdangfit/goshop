package models

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"

	"goshop/pkg/utils"
)

type Product struct {
	Base
	Code        string `json:"code" gorm:"unique;not null;index"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CategID     string `json:"categ_id"`
	Price       uint   `json:"price"`
	Active      bool   `json:"active" gorm:"default:true"`
}

func (product *Product) BeforeCreate(scope *gorm.Scope) error {
	product.ID = uuid.New().String()
	product.Code = utils.GenerateCode("P")
	product.Active = true
	return nil
}
