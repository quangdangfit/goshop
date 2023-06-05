package models

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type Role struct {
	Base
	Name        string `json:"name" gorm:"unique;not null;index"`
	Description string `json:"description"`
}

func (role *Role) BeforeCreate(scope *gorm.Scope) error {
	role.ID = uuid.New().String()
	return nil
}
