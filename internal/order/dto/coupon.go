package dto

import "time"

type Coupon struct {
	ID             string     `json:"id"`
	Code           string     `json:"code"`
	DiscountType   string     `json:"discount_type"`
	DiscountValue  float64    `json:"discount_value"`
	MinOrderAmount float64    `json:"min_order_amount"`
	MaxUsage       int        `json:"max_usage"`
	UsedCount      int        `json:"used_count"`
	ExpiresAt      *time.Time `json:"expires_at"`
}

type CreateCouponReq struct {
	Code           string     `json:"code" validate:"required"`
	DiscountType   string     `json:"discount_type" validate:"required,oneof=fixed percentage"`
	DiscountValue  float64    `json:"discount_value" validate:"required,gt=0"`
	MinOrderAmount float64    `json:"min_order_amount" validate:"gte=0"`
	MaxUsage       int        `json:"max_usage" validate:"gte=0"`
	ExpiresAt      *time.Time `json:"expires_at"`
}
