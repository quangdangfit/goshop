package schema

type Product struct {
	UUID        string `json:"uuid"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CategUUID   string `json:"categ_uuid"`
	Price       uint   `json:"price"`
	Active      bool   `json:"active"`
}

type ProductQueryParam struct {
	Code   string `json:"code,omitempty" form:"code"`
	Active string `json:"active" form:"active"`
}

type ProductBodyParam struct {
	Name        string `json:"name,omitempty" validate:"required"`
	Description string `json:"description,omitempty"`
	CategUUID   string `json:"categ_uuid,omitempty" validate:"required"`
	Price       uint   `json:"price,omitempty"`
}
