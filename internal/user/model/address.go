package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Address struct {
	ID        string     `json:"id" gorm:"unique;not null;index;primary_key"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" gorm:"index"`
	UserID    string     `json:"user_id" gorm:"not null;index"`
	Name      string     `json:"name"`
	Phone     string     `json:"phone"`
	Street    string     `json:"street"`
	City      string     `json:"city"`
	Country   string     `json:"country"`
	IsDefault bool       `json:"is_default" gorm:"default:false"`
}

func (a *Address) BeforeCreate(tx *gorm.DB) error {
	a.ID = uuid.New().String()
	return nil
}
