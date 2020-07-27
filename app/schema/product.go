package schema

type ProductQueryParams struct {
	Code   string `json:"code,omitempty" form:"code"`
	Active string `json:"active,omitempty" form:"active"`
}
