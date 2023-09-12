package dto

type Cart struct {
	ID    string         `json:"id"`
	User  *User          `json:"user"`
	Lines []*CartLineReq `json:"lines"`
}

type CartLine struct {
	Product  *Product `json:"product"`
	Quantity uint     `json:"quantity" validate:"required"`
}

type CartLineReq struct {
	ProductID string `json:"product_id" validate:"required"`
	Quantity  uint   `json:"quantity" validate:"required"`
}

type AddProductReq struct {
	UserID string       `json:"user_id" validate:"required"`
	Line   *CartLineReq `json:"line"  validate:"required,dive"`
}

type RemoveProductReq struct {
	UserID    string `json:"user_id" validate:"required"`
	ProductID string `json:"product_id"  validate:"required"`
}
