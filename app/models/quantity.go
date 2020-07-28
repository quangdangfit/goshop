package models

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type Quantity struct {
	UUID          string `json:"uuid" gorm:"unique;not null;index;primary_key"`
	ProductUUID   string `json:"product_uuid" gorm:"not null;index"`
	WarehouseUUID string `json:"warehouse_uuid" gorm:"not null;index"`
	Quantity      uint   `json:"quantity"`

	gorm.Model
}

func (s *Quantity) BeforeCreate() error {
	s.UUID = uuid.New().String()
	return nil
}
