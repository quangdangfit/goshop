package serializers

type Warehouse struct {
	UUID   string `json:"uuid"`
	Code   string `json:"code"`
	Name   string `json:"name"`
	Active bool   `json:"active"`
}

type WarehouseBodyParam struct {
	Code string `json:"code,omitempty" validate:"required"`
	Name string `json:"name,omitempty" validate:"required"`
}

type WarehouseQueryParam struct {
	Code   string `json:"code,omitempty" form:"code"`
	Active string `json:"active" form:"active"`
}
