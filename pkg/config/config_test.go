package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	cfg := LoadConfig()
	assert.NotNil(t, cfg)
}

func TestGetConfig(t *testing.T) {
	LoadConfig()
	cfg := GetConfig()
	assert.NotNil(t, cfg)
}
