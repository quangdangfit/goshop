package jtoken

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"goshop/pkg/config"
)

func init() {
	config.LoadConfig()
}

func TestGenerateAccessToken(t *testing.T) {
	payload := map[string]interface{}{
		"id":    "user-1",
		"email": "test@example.com",
		"role":  "customer",
	}
	token := GenerateAccessToken(payload)
	assert.NotEmpty(t, token)
	assert.Equal(t, AccessTokenType, payload["type"])
}

func TestGenerateRefreshToken(t *testing.T) {
	payload := map[string]interface{}{
		"id":    "user-1",
		"email": "test@example.com",
		"role":  "customer",
	}
	token := GenerateRefreshToken(payload)
	assert.NotEmpty(t, token)
	assert.Equal(t, RefreshTokenType, payload["type"])
}

func TestValidateToken_Success(t *testing.T) {
	payload := map[string]interface{}{
		"id":    "user-1",
		"email": "test@example.com",
		"role":  "customer",
	}
	token := GenerateAccessToken(payload)

	data, err := ValidateToken(token)
	assert.NoError(t, err)
	assert.NotNil(t, data)
	assert.Equal(t, "user-1", data["id"])
}

func TestValidateToken_WithBearerPrefix(t *testing.T) {
	payload := map[string]interface{}{
		"id":   "user-2",
		"role": "admin",
	}
	token := "Bearer " + GenerateAccessToken(payload)

	data, err := ValidateToken(token)
	assert.NoError(t, err)
	assert.NotNil(t, data)
}

func TestValidateToken_InvalidToken(t *testing.T) {
	data, err := ValidateToken("invalid.token.string")
	assert.Error(t, err)
	assert.Nil(t, data)
}

func TestValidateToken_EmptyToken(t *testing.T) {
	data, err := ValidateToken("")
	assert.Error(t, err)
	assert.Nil(t, data)
}
