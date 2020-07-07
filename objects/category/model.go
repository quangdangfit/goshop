package category

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"goshop/objects/product"
	"goshop/utils"
)

type Category struct {
	UUID        string `gorm:"unique;not null;index;primary_key"`
	Code        string `gorm:"unique;not null;index"`
	Name        string
	Description string
	Active      bool               `gorm:"default:true"`
	Products    *[]product.Product `gorm:"foreignkey:UUID;association_foreignkey:CategUUID"`

	gorm.Model
}

func (categ *Category) BeforeCreate(scope *gorm.Scope) error {
	categ.UUID = uuid.New().String()
	categ.Code = utils.GenerateCode("C")
	return nil
}

type CategoryResponse struct {
	UUID        string `json:"uuid"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Active      bool   `json:"active"`
}

type CategoryRequest struct {
	Name        string `json:"name,omitempty" validate:"required"`
	Description string `json:"description,omitempty"`
}
