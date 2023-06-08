package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Quantity struct {
	Base
	ProductID   string `json:"product_id" gorm:"not null;index"`
	WarehouseID string `json:"warehouse_id" gorm:"not null;index"`
	Quantity    uint   `json:"quantity"`
}

func (s *Quantity) BeforeCreate(tx *gorm.DB) error {
	s.ID = uuid.New().String()
	return nil
}
