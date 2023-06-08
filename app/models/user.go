package models

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"

	"goshop/pkg/utils"
)

type User struct {
	Base
	Email    string `json:"email" gorm:"column:email;unique;not null;index"`
	Password string `json:"password" gorm:"column:password"`
	RoleID   string `json:"role_id" gorm:"column:role_id"`
}

func (user *User) BeforeCreate(scope *gorm.Scope) error {
	user.ID = uuid.New().String()
	user.Password = utils.HashAndSalt([]byte(user.Password))
	return nil
}
