package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Category struct {
	ID          string     `json:"id" gorm:"unique;not null;index;primary_key"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at" gorm:"index"`
	Name        string     `json:"name" gorm:"uniqueIndex;not null"`
	Slug        string     `json:"slug" gorm:"uniqueIndex;not null"`
	Description string     `json:"description"`
}

func (c *Category) BeforeCreate(tx *gorm.DB) error {
	c.ID = uuid.New().String()
	return nil
}
