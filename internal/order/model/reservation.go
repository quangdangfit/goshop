package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ReservationStatus string

const (
	ReservationStatusActive    ReservationStatus = "active"
	ReservationStatusCommitted ReservationStatus = "committed"
	ReservationStatusReleased  ReservationStatus = "released"
)

// StockReservation holds units of a product committed to an in-flight order. While active, the
// quantity counts toward Product.ReservedQuantity. Active reservations expire at ExpiresAt; a
// background sweeper releases expired ones and cancels their parent order if still unpaid.
type StockReservation struct {
	ID        string     `json:"id" gorm:"primary_key"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" gorm:"index"`

	OrderID   string            `json:"order_id" gorm:"index;not null"`
	ProductID string            `json:"product_id" gorm:"index;not null"`
	Quantity  int               `json:"quantity" gorm:"not null"`
	Status    ReservationStatus `json:"status" gorm:"index;not null"`
	ExpiresAt time.Time         `json:"expires_at" gorm:"index;not null"`
}

func (r *StockReservation) BeforeCreate(tx *gorm.DB) error {
	if r.ID == "" {
		r.ID = uuid.New().String()
	}
	if r.Status == "" {
		r.Status = ReservationStatusActive
	}
	return nil
}
