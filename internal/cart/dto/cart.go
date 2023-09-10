package dto

type Cart struct {
	ID     string      `json:"id"`
	UserID string      `json:"user_id"`
	Lines  []*CartLine `json:"lines"`
}

type CartLine struct {
	ProductID string `json:"product_id"`
	Quantity  uint   `json:"quantity"`
}
