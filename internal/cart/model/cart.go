package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Cart struct {
	ID        string     `json:"id" gorm:"unique;not null;index;primary_key"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" gorm:"index"`
	UserID    string     `json:"user_id" gorm:"unique;not null;index"`
	User      *User
	Lines     []*CartLine `json:"lines"`
}

type CartLine struct {
	ProductID string `json:"product_id"`
	Product   *Product
	Quantity  uint `json:"quantity"`
}

func (cart *Cart) BeforeCreate(tx *gorm.DB) error {
	cart.ID = uuid.New().String()
	return nil
}
