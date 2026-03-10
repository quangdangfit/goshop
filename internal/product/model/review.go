package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Review struct {
	ID        string     `json:"id" gorm:"unique;not null;index;primary_key"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" gorm:"index"`
	UserID    string     `json:"user_id" gorm:"uniqueIndex:idx_review_user_product;not null"`
	ProductID string     `json:"product_id" gorm:"uniqueIndex:idx_review_user_product;not null"`
	Rating    int        `json:"rating"`
	Comment   string     `json:"comment"`
}

func (r *Review) BeforeCreate(tx *gorm.DB) error {
	r.ID = uuid.New().String()
	return nil
}
