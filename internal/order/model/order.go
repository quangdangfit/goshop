package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"goshop/pkg/utils"
)

type OrderStatus string

const (
	OrderStatusNew        OrderStatus = "new"
	OrderStatusInProgress OrderStatus = "in-progress"
	OrderStatusDone       OrderStatus = "done"
	OrderStatusCancelled  OrderStatus = "cancelled"
)

func (s OrderStatus) IsValid() bool {
	switch s {
	case OrderStatusNew, OrderStatusInProgress, OrderStatusDone, OrderStatusCancelled:
		return true
	}
	return false
}

type Order struct {
	ID             string     `json:"id" gorm:"unique;not null;index;primary_key"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at" gorm:"index"`
	Code           string     `json:"code"`
	UserID         string     `json:"user_id"`
	User           *User
	Lines          []*OrderLine `json:"lines"`
	TotalPrice     float64      `json:"total_price"`
	DiscountAmount float64      `json:"discount_amount" gorm:"default:0"`
	FinalPrice     float64      `json:"final_price" gorm:"default:0"`
	CouponCode     string       `json:"coupon_code"`
	Status         OrderStatus  `json:"status"`
}

func (order *Order) BeforeCreate(tx *gorm.DB) error {
	order.ID = uuid.New().String()
	order.Code = utils.GenerateCode("SO")

	if order.Status == "" {
		order.Status = OrderStatusNew
	}

	return nil
}
