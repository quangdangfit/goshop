package serializers

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
	Name  string `json:"name,omitempty" form:"name"`
	Code  string `json:"code,omitempty" form:"code"`
	Page  int    `json:"page,omitempty" form:"page"`
	Limit int    `json:"limit,omitempty" form:"limit"`
	Sort  string `json:"sort,omitempty" form:"sort"`
}

type ListProductRes struct {
	Products   []Product         `json:"products,omitempty"`
	Pagination paging.Pagination `json:"pagination,omitempty"`
}

type CreateProductReq struct {
	Name        string  `json:"name,omitempty" validate:"required"`
	Description string  `json:"description,omitempty"`
	Price       float64 `json:"price,omitempty" validate:"gt=0"`
}
