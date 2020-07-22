package models

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type Warehouse struct {
	UUID   string `json:"uuid" gorm:"unique;not null;index;primary_key"`
	Code   string `json:"code" gorm:"unique;not null;index"`
	Name   string `json:"name" gorm:"not null"`
	Active bool   `json:"active"`

	gorm.Model
}

func (w *Warehouse) BeforeCreate() error {
	w.UUID = uuid.New().String()
	w.Active = true
	return nil
}

type WarehouseResponse struct {
	UUID   string `json:"uuid"`
	Code   string `json:"code"`
	Name   string `json:"name"`
	Active bool   `json:"active"`
}

type WarehouseBodyRequest struct {
	Code string `json:"code,omitempty" validate:"required"`
	Name string `json:"name,omitempty" validate:"required"`
}

type WarehouseQueryRequest struct {
	Code   string `json:"code,omitempty" form:"code"`
	Active string `json:"active,omitempty" form:"active"`
}
