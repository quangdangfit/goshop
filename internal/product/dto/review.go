package dto

import (
	"time"

	"goshop/pkg/paging"
)

type Review struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	ProductID string    `json:"product_id"`
	Rating    int       `json:"rating"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateReviewReq struct {
	Rating  int    `json:"rating" validate:"required,gte=1,lte=5"`
	Comment string `json:"comment"`
}

type UpdateReviewReq struct {
	Rating  int    `json:"rating,omitempty" validate:"omitempty,gte=1,lte=5"`
	Comment string `json:"comment,omitempty"`
}

type ListReviewRes struct {
	Reviews    []*Review          `json:"reviews"`
	Pagination *paging.Pagination `json:"pagination"`
}
