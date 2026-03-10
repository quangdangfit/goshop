package dto

import (
	"time"

	"goshop/pkg/paging"
)

type Product struct {
	ID            string    `json:"id"`
	Code          string    `json:"code"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	Price         float64   `json:"price"`
	Active        bool      `json:"active"`
	StockQuantity int       `json:"stock_quantity"`
	AvgRating     float64   `json:"avg_rating"`
	ReviewCount   int       `json:"review_count"`
	CategoryID    string    `json:"category_id,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type ListProductReq struct {
	Name       string `json:"name,omitempty" form:"name"`
	Code       string `json:"code,omitempty" form:"code"`
	CategoryID string `json:"category_id,omitempty" form:"category_id"`
	Page       int64  `json:"-" form:"page"`
	Limit      int64  `json:"-" form:"limit"`
	OrderBy    string `json:"-" form:"order_by"`
	OrderDesc  bool   `json:"-" form:"order_desc"`
}

type ListProductRes struct {
	Products   []*Product         `json:"products"`
	Pagination *paging.Pagination `json:"pagination"`
}

type CreateProductReq struct {
	Name          string  `json:"name" validate:"required"`
	Description   string  `json:"description" validate:"required"`
	Price         float64 `json:"price" validate:"gt=0"`
	StockQuantity int     `json:"stock_quantity" validate:"gte=0"`
	CategoryID    string  `json:"category_id,omitempty"`
}

type UpdateProductReq struct {
	Name          string  `json:"name,omitempty"`
	Description   string  `json:"description,omitempty"`
	Price         float64 `json:"price,omitempty" validate:"gte=0"`
	StockQuantity *int    `json:"stock_quantity,omitempty" validate:"omitempty,gte=0"`
	CategoryID    string  `json:"category_id,omitempty"`
}
