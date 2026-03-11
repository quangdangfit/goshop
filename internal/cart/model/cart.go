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
	Lines     []*CartLine `json:"lines" gorm:"foreignKey:CartID"`
}

type CartLine struct {
	ID        string   `json:"id" gorm:"unique;not null;index;primary_key"`
	CartID    string   `json:"cart_id" gorm:"not null;index"`
	ProductID string   `json:"product_id" gorm:"not null;index"`
	Product   *Product `gorm:"foreignKey:ProductID;references:ID"`
	Quantity  uint     `json:"quantity"`
}

func (cart *Cart) BeforeCreate(tx *gorm.DB) error {
	cart.ID = uuid.New().String()
	return nil
}

func (line *CartLine) BeforeCreate(tx *gorm.DB) error {
	line.ID = uuid.New().String()
	return nil
}
