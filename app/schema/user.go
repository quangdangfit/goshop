package schema

type User struct {
	UUID     string      `json:"uuid"`
	Username string      `json:"username"`
	Email    string      `json:"email"`
	Extra    interface{} `json:"extra,omitempty"`
}

type Register struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
	RoleUUID string `json:"role_uuid"`
}

type Login struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}
