package serializers

type Order struct {
	UUID       string      `json:"uuid"`
	Code       string      `json:"code"`
	Lines      []OrderLine `json:"lines"`
	TotalPrice uint        `json:"total_price"`
	Status     string      `json:"status"`
}

type OrderBodyParam struct {
	Lines []OrderLineBodyParam `json:"lines,omitempty" validate:"required"`
}

type OrderQueryParam struct {
	Code   string `json:"code,omitempty" form:"code"`
	Status string `json:"status,omitempty" form:"active"`
}
