package models

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type User struct {
	Base
	Email    string `json:"email" gorm:"unique;not null;index"`
	Password string `json:"password"`
	RoleID   string `json:"role_id"`
}

func (user *User) BeforeCreate(scope *gorm.Scope) error {
	user.ID = uuid.New().String()
	return nil
}
