package models

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type OrderLine struct {
	UUID        string `json:"uuid" gorm:"unique;not null;index;primary_key"`
	ProductUUID string `json:"product_uuid"`
	OrderUUID   string `json:"order_uuid"`
	Quantity    uint   `json:"quantity"`
	TotalPrice  uint   `json:"total_price"`

	gorm.Model
}

func (line *OrderLine) BeforeCreate() error {
	line.UUID = uuid.New().String()
	return nil
}

type OrderLineResponse struct {
	UUID        string `json:"uuid"`
	ProductUUID string `json:"product_uuid"`
	Quantity    string `json:"quantity"`
	TotalPrice  bool   `json:"total_price"`
}

type OrderLineBodyRequest struct {
	ProductUUID string `json:"product_uuid,omitempty" validate:"required"`
	OrderUUID   string `json:"order_uuid,omitempty" validate:"required"`
	Quantity    string `json:"quantity,omitempty" validate:"required"`
}
