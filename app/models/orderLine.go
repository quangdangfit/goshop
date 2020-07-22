package models

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"

	"goshop/dbs"
)

type OrderLine struct {
	UUID        string `json:"uuid" gorm:"unique;not null;index;primary_key"`
	ProductUUID string `json:"product_uuid"`
	OrderUUID   string `json:"order_uuid"`
	Quantity    uint   `json:"quantity"`
	Price       uint   `json:"price"`

	gorm.Model
}

func (line *OrderLine) BeforeCreate(scope *gorm.Scope) error {
	line.UUID = uuid.New().String()
	var product Product
	dbs.Database.Where("uuid = ?", line.ProductUUID).First(&product)
	line.Price = product.Price * line.Quantity

	return nil
}

type OrderLineResponse struct {
	UUID        string `json:"uuid"`
	ProductUUID string `json:"product_uuid"`
	Quantity    uint   `json:"quantity"`
	Price       uint   `json:"price"`
}

type OrderLineBodyRequest struct {
	ProductUUID string `json:"product_uuid,omitempty" validate:"required"`
	Quantity    uint   `json:"quantity,omitempty" validate:"required"`
}
