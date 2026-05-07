package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// DeadLetterNotification is an audit row for notifications that exhausted their retry budget.
// Operators can inspect this table to triage delivery failures and trigger manual replays.
type DeadLetterNotification struct {
	ID        string    `json:"id" gorm:"primary_key"`
	CreatedAt time.Time `json:"created_at"`

	EventType string `json:"event_type" gorm:"index;size:64;not null"`
	UserEmail string `json:"user_email" gorm:"size:255;not null"`
	Payload   string `json:"payload" gorm:"type:text"`
	LastError string `json:"last_error" gorm:"type:text"`
}

func (d *DeadLetterNotification) BeforeCreate(tx *gorm.DB) error {
	if d.ID == "" {
		d.ID = uuid.New().String()
	}
	return nil
}
