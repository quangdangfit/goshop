package serializers

type Role struct {
	UUID        string `json:"uuid"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type RoleBodyParam struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}
