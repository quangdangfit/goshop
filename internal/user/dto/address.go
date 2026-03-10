package dto

import "time"

type Address struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Name      string    `json:"name"`
	Phone     string    `json:"phone"`
	Street    string    `json:"street"`
	City      string    `json:"city"`
	Country   string    `json:"country"`
	IsDefault bool      `json:"is_default"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateAddressReq struct {
	Name    string `json:"name" validate:"required"`
	Phone   string `json:"phone" validate:"required"`
	Street  string `json:"street" validate:"required"`
	City    string `json:"city" validate:"required"`
	Country string `json:"country" validate:"required"`
}

type UpdateAddressReq struct {
	Name    string `json:"name,omitempty"`
	Phone   string `json:"phone,omitempty"`
	Street  string `json:"street,omitempty"`
	City    string `json:"city,omitempty"`
	Country string `json:"country,omitempty"`
}
