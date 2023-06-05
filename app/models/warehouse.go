package models

import (
	"github.com/google/uuid"
)

type Warehouse struct {
	Base
	Code   string `json:"code" gorm:"unique;not null;index"`
	Name   string `json:"name" gorm:"not null"`
	Active bool   `json:"active"`
}

func (w *Warehouse) BeforeCreate() error {
	w.ID = uuid.New().String()
	w.Active = true
	return nil
}
