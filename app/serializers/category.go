package serializers

type Category struct {
	UUID        string `json:"uuid"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Active      bool   `json:"active"`
}

type CategoryBodyParam struct {
	Name        string `json:"name,omitempty" validate:"required"`
	Description string `json:"description,omitempty"`
}

type CategoryQueryParam struct {
	Code   string `json:"code,omitempty" form:"code"`
	Active bool   `json:"active" form:"active"`
}
