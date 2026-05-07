package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"goshop/pkg/utils"
)

type OrderStatus string

const (
	OrderStatusNew            OrderStatus = "new"
	OrderStatusPendingPayment OrderStatus = "pending_payment"
	OrderStatusPaid           OrderStatus = "paid"
	OrderStatusInProgress     OrderStatus = "in-progress"
	OrderStatusDone           OrderStatus = "done"
	OrderStatusCancelled      OrderStatus = "cancelled"
	OrderStatusPaymentFailed  OrderStatus = "payment_failed"
)

func (s OrderStatus) IsValid() bool {
	switch s {
	case OrderStatusNew, OrderStatusPendingPayment, OrderStatusPaid,
		OrderStatusInProgress, OrderStatusDone, OrderStatusCancelled, OrderStatusPaymentFailed:
		return true
	}
	return false
}

// allowedTransitions maps each status to the set of statuses it can advance to. Terminal
// statuses (done, cancelled) have no outbound transitions. Designed for admin-driven
// fulfillment moves; payment-driven transitions (pending_payment -> paid/payment_failed)
// also live here so the webhook path goes through the same gate.
var allowedTransitions = map[OrderStatus]map[OrderStatus]struct{}{
	OrderStatusNew:            {OrderStatusInProgress: {}, OrderStatusCancelled: {}},
	OrderStatusPendingPayment: {OrderStatusPaid: {}, OrderStatusPaymentFailed: {}, OrderStatusCancelled: {}},
	OrderStatusPaid:           {OrderStatusInProgress: {}, OrderStatusCancelled: {}},
	OrderStatusInProgress:     {OrderStatusDone: {}, OrderStatusCancelled: {}},
	OrderStatusPaymentFailed:  {OrderStatusCancelled: {}, OrderStatusPendingPayment: {}},
	OrderStatusDone:           {},
	OrderStatusCancelled:      {},
}

// CanTransitionTo reports whether `s` is allowed to advance to `next`. Idempotent moves
// (s == next) are accepted so retried webhooks don't fail the gate.
func (s OrderStatus) CanTransitionTo(next OrderStatus) bool {
	if s == next {
		return true
	}
	allowed, ok := allowedTransitions[s]
	if !ok {
		return false
	}
	_, ok = allowed[next]
	return ok
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
