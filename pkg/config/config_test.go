package config

import (
	"os"
	"path/filepath"
	"runtime"
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

func TestLoadConfig_PopulatesFromYAML(t *testing.T) {
	cfg := LoadConfig()

	// Values come from pkg/config/config.yaml shipped with the repo.
	assert.Equal(t, ProductionEnv, cfg.Environment)
	assert.Equal(t, 8888, cfg.HttpPort)
	assert.Equal(t, 8889, cfg.GrpcPort)
	assert.Equal(t, "localhost:6379", cfg.RedisURI)
	assert.Equal(t, 0, cfg.RedisDB)
}

func TestLoadConfig_AppliesDefaultsForUnsetFields(t *testing.T) {
	t.Setenv("cors_allowed_origins", "")
	t.Setenv("rate_limit_requests", "")
	t.Setenv("rate_limit_window_seconds", "")

	cfg := LoadConfig()

	// envDefault tags kick in when the env var is missing/empty.
	assert.Equal(t, "*", cfg.CORSAllowedOrigins)
	assert.Equal(t, 100, cfg.RateLimitRequests)
	assert.Equal(t, 60, cfg.RateLimitWindowSeconds)
}

func TestLoadConfig_OverridesFromEnv(t *testing.T) {
	t.Setenv("cors_allowed_origins", "https://example.com")
	t.Setenv("rate_limit_requests", "250")
	t.Setenv("rate_limit_window_seconds", "30")

	cfg := LoadConfig()

	assert.Equal(t, "https://example.com", cfg.CORSAllowedOrigins)
	assert.Equal(t, 250, cfg.RateLimitRequests)
	assert.Equal(t, 30, cfg.RateLimitWindowSeconds)
}

func TestGetConfig_ReturnsSingleton(t *testing.T) {
	LoadConfig()
	a := GetConfig()
	b := GetConfig()
	assert.Same(t, a, b, "GetConfig should return the same package-level instance")
}

func TestAuthIgnoreMethods(t *testing.T) {
	assert.Contains(t, AuthIgnoreMethods, "/user.UserService/Login")
	assert.Contains(t, AuthIgnoreMethods, "/user.UserService/Register")
}

// Exercises the godotenv.Load error branch by temporarily moving config.yaml
// out of the way. LoadConfig must still succeed because env.Parse can
// populate the schema from environment variables alone.
func TestLoadConfig_GodotenvLoadError(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)
	original := filepath.Join(dir, "config.yaml")
	moved := filepath.Join(dir, "config.yaml.bak")

	if err := os.Rename(original, moved); err != nil {
		t.Skipf("cannot move config.yaml: %v", err)
	}
	t.Cleanup(func() { _ = os.Rename(moved, original) })

	cfg := LoadConfig()
	assert.NotNil(t, cfg)
}
