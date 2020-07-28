package models

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"

	"goshop/pkg/utils"
)

const (
	OrderStatusNew      = "new"
	OrderStatusAssigned = "assigned"
	OrderStatusDone     = "done"
	OrderStatusCanceled = "canceled"
)

type Order struct {
	UUID       string      `json:"uuid" gorm:"unique;not null;index;primary_key"`
	Code       string      `json:"code"`
	Lines      []OrderLine `json:"lines" gorm:"foreignkey:order_uuid;association_foreignkey:uuid;save_associations:false"`
	TotalPrice uint        `json:"total_price"`
	Status     string      `json:"status"`

	gorm.Model
}

func (order *Order) BeforeCreate(scope *gorm.Scope) error {
	order.UUID = uuid.New().String()
	order.Code = utils.GenerateCode("SO")
	order.Status = OrderStatusNew

	return nil
}
