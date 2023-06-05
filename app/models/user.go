package models

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type User struct {
	Base
	Username string `json:"username" gorm:"unique;not null;index"`
	Email    string `json:"email" gorm:"unique;not null;index"`
	Password string `json:"password"`
	RoleUUID string `json:"role_uuid"`
}

func (user *User) BeforeCreate(scope *gorm.Scope) error {
	user.ID = uuid.New().String()
	return nil
}
