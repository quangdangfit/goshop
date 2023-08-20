package dto

import (
	"time"

	"goshop/pkg/paging"
)

type Product struct {
	ID          string    `json:"id"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Active      bool      `json:"active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ListProductReq struct {
	Name      string `json:"name,omitempty" form:"name"`
	Code      string `json:"code,omitempty" form:"code"`
	Page      int64  `json:"-" form:"page"`
	Limit     int64  `json:"-" form:"limit"`
	OrderBy   string `json:"-" form:"order_by"`
	OrderDesc bool   `json:"-" form:"order_desc"`
}

type ListProductRes struct {
	Products   []*Product         `json:"products"`
	Pagination *paging.Pagination `json:"pagination"`
}

type CreateProductReq struct {
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description" validate:"required"`
	Price       float64 `json:"price" validate:"gt=0"`
}

type UpdateProductReq struct {
	Name        string  `json:"name,omitempty"`
	Description string  `json:"description,omitempty"`
	Price       float64 `json:"price,omitempty" validate:"gte=0"`
}
