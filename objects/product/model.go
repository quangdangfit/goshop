package product

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"goshop/utils"
)

type Product struct {
	UUID        string `json:"uuid,omitempty" bson:"uuid,omitempty" gorm:"unique;not null;index"`
	Code        string `json:"code,omitempty" bson:"code,omitempty" gorm:"unique;not null;index"`
	Name        string `json:"name,omitempty" bson:"name,omitempty"`
	Description string `json:"description,omitempty" bson:"description,omitempty"`
	CategUUID   string `json:"categ_uuid,omitempty" bson:"categ_uuid,omitempty"`
	Active      bool   `json:"active,omitempty" bson:"active,omitempty" gorm:"default:true"`
	gorm.Model
}

func (product *Product) BeforeCreate(scope *gorm.Scope) error {
	product.UUID = uuid.New().String()
	product.Code = utils.GenerateCode("C")
	return nil
}

type ProductResponse struct {
	UUID        string `json:"uuid" bson:"uuid"`
	Code        string `json:"code" bson:"code"`
	Name        string `json:"name" bson:"name"`
	Description string `json:"description" bson:"description"`
	CategUUID   string `json:"categ_uuid" bson:"code"`
	Active      bool   `json:"active" bson:"active"`
}

type ProductRequest struct {
	Name        string `json:"name,omitempty" bson:"name,omitempty" validate:"required"`
	Description string `json:"description,omitempty" bson:"description,omitempty"`
	CategUUID   string `json:"categ_uuid,omitempty" bson:"categ_uuid,omitempty" validate:"required"`
}
