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

type QuantityResponse struct {
	UUID          string `json:"uuid"`
	ProductUUID   string `json:"product_uuid"`
	WarehouseUUID string `json:"warehouse_uuid"`
	Quantity      uint   `json:"quantity"`
}

type QuantityBodyRequest struct {
	ProductUUID   string `json:"product_uuid,omitempty" validate:"required"`
	WarehouseUUID string `json:"warehouse_uuid,omitempty" validate:"required"`
	Quantity      uint   `json:"quantity" validate:"required"`
}

type QuantityQueryRequest struct {
	ProductUUID   string `json:"product_uuid,omitempty" form:"code"`
	WarehouseUUID string `json:"warehouse_uuid,omitempty" form:"active"`
}
