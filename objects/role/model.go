package role

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type Role struct {
	UUID        string `json:"uuid" gorm:"unique;not null;index;primary_key"`
	Name        string `json:"name" gorm:"unique;not null;index"`
	Description string `json:"description"`

	gorm.Model
}

func (role *Role) BeforeCreate(scope *gorm.Scope) error {
	role.UUID = uuid.New().String()
	return nil
}

type RoleResponse struct {
	UUID        string `json:"uuid"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type RoleRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}
