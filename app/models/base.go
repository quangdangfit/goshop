package models

import (
	"time"
)

type Base struct {
	ID        string `json:"id" gorm:"unique;not null;index;primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}
