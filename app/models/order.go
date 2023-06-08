package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"goshop/pkg/utils"
)

const (
	OrderStatusNew      = "new"
	OrderStatusAssigned = "assigned"
	OrderStatusDone     = "done"
	OrderStatusCanceled = "canceled"
)

type Order struct {
	Base
	Code       string      `json:"code"`
	Lines      []OrderLine `json:"lines" gorm:"foreignkey:order_uuid;association_foreignkey:uuid;save_associations:false"`
	TotalPrice float64     `json:"total_price"`
	Status     string      `json:"status"`
}

func (order *Order) BeforeCreate(tx *gorm.DB) error {
	order.ID = uuid.New().String()
	order.Code = utils.GenerateCode("SO")
	order.Status = OrderStatusNew

	return nil
}
