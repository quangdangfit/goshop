package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Wishlist struct {
	ID        string    `json:"id" gorm:"unique;not null;index;primary_key"`
	CreatedAt time.Time `json:"created_at"`
	UserID    string    `json:"user_id" gorm:"uniqueIndex:idx_wishlist_user_product;not null"`
	ProductID string    `json:"product_id" gorm:"uniqueIndex:idx_wishlist_user_product;not null"`
}

func (w *Wishlist) BeforeCreate(tx *gorm.DB) error {
	w.ID = uuid.New().String()
	return nil
}
