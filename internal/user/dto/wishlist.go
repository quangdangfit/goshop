package dto

import "time"

type WishlistItem struct {
	ID        string    `json:"id"`
	ProductID string    `json:"product_id"`
	CreatedAt time.Time `json:"created_at"`
}

type AddToWishlistReq struct {
	ProductID string `json:"product_id" validate:"required"`
}
