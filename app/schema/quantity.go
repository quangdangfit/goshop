package schema

type Quantity struct {
	UUID          string `json:"uuid"`
	ProductUUID   string `json:"product_uuid"`
	WarehouseUUID string `json:"warehouse_uuid"`
	Quantity      uint   `json:"quantity"`
}

type QuantityQueryParam struct {
	ProductUUID   string `json:"product_uuid,omitempty" form:"code"`
	WarehouseUUID string `json:"warehouse_uuid,omitempty" form:"active"`
}

type QuantityBodyParam struct {
	ProductUUID   string `json:"product_uuid,omitempty" validate:"required"`
	WarehouseUUID string `json:"warehouse_uuid,omitempty" validate:"required"`
	Quantity      uint   `json:"quantity" validate:"required"`
}
