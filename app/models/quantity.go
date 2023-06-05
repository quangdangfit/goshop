package models

import (
	"github.com/google/uuid"
)

type Quantity struct {
	Base
	ProductID   string `json:"product_id" gorm:"not null;index"`
	WarehouseID string `json:"warehouse_id" gorm:"not null;index"`
	Quantity    uint   `json:"quantity"`
}

func (s *Quantity) BeforeCreate() error {
	s.ID = uuid.New().String()
	return nil
}
