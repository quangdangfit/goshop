package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Preference stores one toggle: user × event_type × channel. Absence of a row is treated
// as "enabled" by the lookup helper, so users start fully opted in by default.
type Preference struct {
	ID        string    `json:"id" gorm:"primary_key"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	UserID    string `json:"user_id" gorm:"uniqueIndex:idx_pref_user_event_channel;not null"`
	EventType string `json:"event_type" gorm:"uniqueIndex:idx_pref_user_event_channel;size:64;not null"`
	Channel   string `json:"channel" gorm:"uniqueIndex:idx_pref_user_event_channel;size:32;not null"`
	Enabled   bool   `json:"enabled" gorm:"not null;default:true"`
}

func (p *Preference) BeforeCreate(tx *gorm.DB) error {
	if p.ID == "" {
		p.ID = uuid.New().String()
	}
	return nil
}
