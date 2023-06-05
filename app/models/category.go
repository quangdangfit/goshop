package models

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"

	"goshop/pkg/utils"
)

type Category struct {
	Base
	Code        string `json:"code" gorm:"unique;not null;index"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Active      bool   `json:"active" gorm:"default:true"`
}

func (categ *Category) BeforeCreate(scope *gorm.Scope) error {
	categ.ID = uuid.New().String()
	categ.Code = utils.GenerateCode("C")
	categ.Active = true
	return nil
}
