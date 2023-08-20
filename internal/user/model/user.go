package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"goshop/pkg/utils"
)

type UserRole string

const (
	UserRoleAdmin    UserRole = "admin"
	UserRoleCustomer UserRole = "customer"
)

type User struct {
	ID        string     `json:"id" gorm:"unique;not null;index;primary_key"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" gorm:"index"`
	Email     string     `json:"email" gorm:"unique;not null;index:idx_user_email"`
	Password  string     `json:"password"`
	Role      UserRole   `json:"role"`
}

func (user *User) BeforeCreate(tx *gorm.DB) error {
	user.ID = uuid.New().String()
	user.Password = utils.HashAndSalt([]byte(user.Password))
	if user.Role == "" {
		user.Role = UserRoleCustomer
	}
	return nil
}
