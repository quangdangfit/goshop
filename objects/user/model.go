package user

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type User struct {
	UUID     string `json:"uuid" gorm:"unique;not null;index;primary_key"`
	Username string `json:"username" gorm:"unique;not null;index"`
	Email    string `json:"email" gorm:"unique;not null;index"`
	Password string `json:"password"`
	RoleUUID string `json:"role_uuid"`

	gorm.Model
}

func (user *User) BeforeCreate(scope *gorm.Scope) error {
	user.UUID = uuid.New().String()
	return nil
}

type UserResponse struct {
	UUID     string      `json:"uuid"`
	Username string      `json:"username"`
	Email    string      `json:"email"`
	Extra    interface{} `json:"extra,omitempty"`
}

type RegisterRequest struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
	RoleUUID string `json:"role_uuid"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}
