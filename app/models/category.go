package models

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"

	"goshop/pkg/utils"
)

type Category struct {
	UUID        string `json:"uuid" gorm:"unique;not null;index;primary_key"`
	Code        string `json:"code" gorm:"unique;not null;index"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Active      bool   `json:"active" gorm:"default:true"`

	gorm.Model
}

func (categ *Category) BeforeCreate(scope *gorm.Scope) error {
	categ.UUID = uuid.New().String()
	categ.Code = utils.GenerateCode("C")
	categ.Active = true
	return nil
}
