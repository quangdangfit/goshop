package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUser_BeforeCreate_DefaultRole(t *testing.T) {
	user := &User{Password: "secret"}
	err := user.BeforeCreate(nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, user.ID)
	assert.Equal(t, UserRoleCustomer, user.Role)
	assert.NotEqual(t, "secret", user.Password) // password should be hashed
}

func TestUser_BeforeCreate_AdminRole(t *testing.T) {
	user := &User{Password: "pass", Role: UserRoleAdmin}
	err := user.BeforeCreate(nil)
	assert.NoError(t, err)
	assert.Equal(t, UserRoleAdmin, user.Role)
}
