package dto

type Cart struct {
	ID     string      `json:"id"`
	UserID string      `json:"user_id"`
	Lines  []*CartLine `json:"lines"`
}

type CartLine struct {
	ProductID string `json:"product_id" validate:"required"`
	Quantity  uint   `json:"quantity" validate:"required"`
}

type AddProductReq struct {
	UserID string    `json:"user_id" validate:"required"`
	Line   *CartLine `json:"line"  validate:"required,dive"`
}

type RemoveProductReq struct {
	UserID    string `json:"user_id" validate:"required"`
	ProductID string `json:"product_id"  validate:"required"`
}
