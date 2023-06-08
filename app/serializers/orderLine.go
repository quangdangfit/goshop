package serializers

type OrderLine struct {
	UUID        string `json:"uuid"`
	ProductUUID string `json:"product_uuid"`
	Quantity    uint   `json:"quantity"`
	Price       uint   `json:"price"`
}

type OrderLineQueryParam struct {
	ProductUUID string `json:"product_uuid,omitempty" validate:"required"`
}

type OrderLineBodyParam struct {
	ProductUUID string `json:"product_uuid,omitempty" validate:"required"`
	Quantity    uint   `json:"quantity,omitempty" validate:"required"`
}
