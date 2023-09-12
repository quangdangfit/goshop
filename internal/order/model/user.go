package model

import (
	"time"
)

type User struct {
	ID        string     `json:"id" gorm:"unique;not null;index;primary_key"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" gorm:"index"`
	Email     string     `json:"email" gorm:"unique;not null;index:idx_user_email"`
}
