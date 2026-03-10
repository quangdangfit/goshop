package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DiscountType string

const (
	DiscountTypeFixed      DiscountType = "fixed"
	DiscountTypePercentage DiscountType = "percentage"
)

type Coupon struct {
	ID             string       `json:"id" gorm:"unique;not null;index;primary_key"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
	DeletedAt      *time.Time   `json:"deleted_at" gorm:"index"`
	Code           string       `json:"code" gorm:"uniqueIndex;not null"`
	DiscountType   DiscountType `json:"discount_type"`
	DiscountValue  float64      `json:"discount_value"`
	MinOrderAmount float64      `json:"min_order_amount" gorm:"default:0"`
	MaxUsage       int          `json:"max_usage" gorm:"default:0"`
	UsedCount      int          `json:"used_count" gorm:"default:0"`
	ExpiresAt      *time.Time   `json:"expires_at"`
}

func (c *Coupon) BeforeCreate(tx *gorm.DB) error {
	c.ID = uuid.New().String()
	return nil
}
