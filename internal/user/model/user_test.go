package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUser_BeforeCreate(t *testing.T) {
	tests := []struct {
		name         string
		user         *User
		expectedRole UserRole
		checkHash    bool
	}{
		{
			name:         "DefaultRole",
			user:         &User{Password: "secret"},
			expectedRole: UserRoleCustomer,
			checkHash:    true,
		},
		{
			name:         "AdminRole",
			user:         &User{Password: "pass", Role: UserRoleAdmin},
			expectedRole: UserRoleAdmin,
			checkHash:    false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			originalPassword := tc.user.Password
			err := tc.user.BeforeCreate(nil)
			assert.NoError(t, err)
			assert.NotEmpty(t, tc.user.ID)
			assert.Equal(t, tc.expectedRole, tc.user.Role)
			if tc.checkHash {
				assert.NotEqual(t, originalPassword, tc.user.Password)
			}
		})
	}
}
