package category

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"goshop/utils"
)

type Category struct {
	UUID        string `json:"uuid,omitempty" bson:"uuid,omitempty" gorm:"unique;not null;index"`
	Code        string `json:"code,omitempty" bson:"code,omitempty" gorm:"unique;not null;index"`
	Name        string `json:"name,omitempty" bson:"name,omitempty"`
	Description string `json:"description,omitempty" bson:"description,omitempty"`
	Active      bool   `json:"active,omitempty" bson:"active,omitempty" gorm:"default:true"`
	gorm.Model
}

func (categ *Category) BeforeCreate(scope *gorm.Scope) error {
	categ.UUID = uuid.New().String()
	categ.Code = utils.GenerateCode("C")
	return nil
}

type CategoryResponse struct {
	UUID        string `json:"uuid" bson:"uuid"`
	Code        string `json:"code" bson:"code"`
	Name        string `json:"name" bson:"name"`
	Description string `json:"description" bson:"description"`
	Active      bool   `json:"active" bson:"active"`
}

type CategoryRequest struct {
	Name        string `json:"name,omitempty" bson:"name,omitempty" validate:"required"`
	Description string `json:"description,omitempty" bson:"description,omitempty"`
}
