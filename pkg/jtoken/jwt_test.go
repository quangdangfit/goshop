package jtoken

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"goshop/pkg/config"
)

func init() {
	config.LoadConfig()
}

func TestGenerateToken(t *testing.T) {
	tests := []struct {
		name     string
		generate func(payload map[string]interface{}) string
		wantType string
	}{
		{
			name:     "access token",
			generate: GenerateAccessToken,
			wantType: AccessTokenType,
		},
		{
			name:     "refresh token",
			generate: GenerateRefreshToken,
			wantType: RefreshTokenType,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			payload := map[string]interface{}{
				"id":    "user-1",
				"email": "test@example.com",
				"role":  "customer",
			}
			token := tc.generate(payload)
			assert.NotEmpty(t, token)
			assert.Equal(t, tc.wantType, payload["type"])
		})
	}
}

func TestValidateToken(t *testing.T) {
	tests := []struct {
		name    string
		token   func() string
		wantErr bool
		wantID  string
	}{
		{
			name: "valid access token",
			token: func() string {
				return GenerateAccessToken(map[string]interface{}{
					"id":    "user-1",
					"email": "test@example.com",
					"role":  "customer",
				})
			},
			wantErr: false,
			wantID:  "user-1",
		},
		{
			name: "with Bearer prefix",
			token: func() string {
				return "Bearer " + GenerateAccessToken(map[string]interface{}{
					"id":   "user-2",
					"role": "admin",
				})
			},
			wantErr: false,
			wantID:  "user-2",
		},
		{
			name: "invalid token",
			token: func() string {
				return "invalid.token.string"
			},
			wantErr: true,
		},
		{
			name: "empty token",
			token: func() string {
				return ""
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			data, err := ValidateToken(tc.token())
			if tc.wantErr {
				assert.Error(t, err)
				assert.Nil(t, data)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, data)
				if tc.wantID != "" {
					assert.Equal(t, tc.wantID, data["id"])
				}
			}
		})
	}
}

func TestGenerateToken_Returns(t *testing.T) {
	tests := []struct {
		name     string
		generate func(payload map[string]interface{}) string
	}{
		{
			name:     "access token returns string",
			generate: GenerateAccessToken,
		},
		{
			name:     "refresh token returns string",
			generate: GenerateRefreshToken,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cfg := config.GetConfig()
			orig := cfg.AuthSecret
			defer func() { cfg.AuthSecret = orig }()

			token := tc.generate(map[string]interface{}{
				"id":   "user-1",
				"role": "customer",
			})
			assert.IsType(t, "", token)
		})
	}
}
