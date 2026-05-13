package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderLine struct {
	ID        string     `json:"id" gorm:"unique;not null;index;primary_key"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" gorm:"index"`
	OrderID   string     `json:"order_id"`
	ProductID string     `json:"product_id"`
	Product   *Product
	Quantity  uint    `json:"quantity"`
	Price     float64 `json:"price"`
}

// BeforeCreate generates a UUID only when one isn't already set. Unconditional
// assignment was a subtle bug: GORM's Save(parentOrder) upserts the has-many
// Lines slice as INSERT…ON CONFLICT, which fires BeforeCreate; overwriting the
// ID gave each insert a fresh PK so the conflict never hit and we ended up with
// duplicate rows on every Save of the parent order.
func (line *OrderLine) BeforeCreate(tx *gorm.DB) error {
	if line.ID == "" {
		line.ID = uuid.New().String()
	}
	return nil
}
