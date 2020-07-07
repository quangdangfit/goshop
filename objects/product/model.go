package product

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"goshop/utils"
)

type Product struct {
	UUID        string `gorm:"unique;not null;index;primary_key"`
	Code        string `gorm:"unique;not null;index"`
	Name        string
	Description string
	CategUUID   string
	Active      bool `gorm:"default:true"`

	gorm.Model
}

func (product *Product) BeforeCreate(scope *gorm.Scope) error {
	product.UUID = uuid.New().String()
	product.Code = utils.GenerateCode("C")
	return nil
}

type ProductResponse struct {
	UUID        string `json:"uuid"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CategUUID   string `json:"categ_uuid"`
	Active      bool   `json:"active"`
}

type ProductRequest struct {
	Name        string `json:"name,omitempty" validate:"required"`
	Description string `json:"description,omitempty"`
	CategUUID   string `json:"categ_uuid,omitempty" validate:"required"`
}
