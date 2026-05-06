package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusSucceeded PaymentStatus = "succeeded"
	PaymentStatusFailed    PaymentStatus = "failed"
)

// Payment is the local record of a charge attempt against an external provider. Each Order
// has at most one active Payment; a previously-failed payment can be superseded by a new one.
type Payment struct {
	ID        string     `json:"id" gorm:"primary_key"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" gorm:"index"`

	OrderID          string        `json:"order_id" gorm:"uniqueIndex;not null"`
	Provider         string        `json:"provider" gorm:"not null"`
	ProviderIntentID string        `json:"provider_intent_id" gorm:"index;not null"`
	Amount           int64         `json:"amount" gorm:"not null"`
	Currency         string        `json:"currency" gorm:"not null"`
	Status           PaymentStatus `json:"status" gorm:"index;not null"`
}

func (p *Payment) BeforeCreate(tx *gorm.DB) error {
	if p.ID == "" {
		p.ID = uuid.New().String()
	}
	if p.Status == "" {
		p.Status = PaymentStatusPending
	}
	return nil
}

// ProviderEvent is the dedup table for provider webhook deliveries. Inserting a row with the
// same (provider, event_id) is rejected by the unique index, giving exactly-once processing.
type ProviderEvent struct {
	CreatedAt time.Time `json:"created_at"`
	Provider  string    `json:"provider" gorm:"primaryKey;size:32"`
	EventID   string    `json:"event_id" gorm:"primaryKey;size:128"`
}
