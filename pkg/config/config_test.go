package config

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

// repoRoot returns the repository root so tests can locate the top-level
// config.yaml regardless of the directory `go test` is invoked from.
func repoRoot(t *testing.T) string {
	t.Helper()
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(filename), "..", "..")
}

func TestLoadConfig(t *testing.T) {
	t.Setenv("CONFIG_FILE", filepath.Join(repoRoot(t), "config.yaml"))
	cfg := LoadConfig()
	assert.NotNil(t, cfg)
}

func TestGetConfig(t *testing.T) {
	t.Setenv("CONFIG_FILE", filepath.Join(repoRoot(t), "config.yaml"))
	LoadConfig()
	cfg := GetConfig()
	assert.NotNil(t, cfg)
}

func TestLoadConfig_PopulatesFromYAML(t *testing.T) {
	t.Setenv("CONFIG_FILE", filepath.Join(repoRoot(t), "config.yaml"))
	cfg := LoadConfig()

	assert.Equal(t, ProductionEnv, cfg.Environment)
	assert.Equal(t, 8888, cfg.HttpPort)
	assert.Equal(t, 8889, cfg.GrpcPort)
	assert.Equal(t, "localhost:6379", cfg.RedisURI)
	assert.Equal(t, 0, cfg.RedisDB)
}

func TestLoadConfig_AppliesDefaultsForUnsetFields(t *testing.T) {
	t.Setenv("CONFIG_FILE", filepath.Join(repoRoot(t), "config.yaml"))
	t.Setenv("cors_allowed_origins", "")
	t.Setenv("rate_limit_requests", "")
	t.Setenv("rate_limit_window_seconds", "")

	cfg := LoadConfig()

	assert.Equal(t, "*", cfg.CORSAllowedOrigins)
	assert.Equal(t, 100, cfg.RateLimitRequests)
	assert.Equal(t, 60, cfg.RateLimitWindowSeconds)
}

func TestLoadConfig_OverridesFromEnv(t *testing.T) {
	t.Setenv("CONFIG_FILE", filepath.Join(repoRoot(t), "config.yaml"))
	t.Setenv("cors_allowed_origins", "https://example.com")
	t.Setenv("rate_limit_requests", "250")
	t.Setenv("rate_limit_window_seconds", "30")

	cfg := LoadConfig()

	assert.Equal(t, "https://example.com", cfg.CORSAllowedOrigins)
	assert.Equal(t, 250, cfg.RateLimitRequests)
	assert.Equal(t, 30, cfg.RateLimitWindowSeconds)
}

func TestGetConfig_ReturnsSingleton(t *testing.T) {
	t.Setenv("CONFIG_FILE", filepath.Join(repoRoot(t), "config.yaml"))
	LoadConfig()
	a := GetConfig()
	b := GetConfig()
	assert.Same(t, a, b, "GetConfig should return the same package-level instance")
}

func TestAuthIgnoreMethods(t *testing.T) {
	assert.Contains(t, AuthIgnoreMethods, "/user.UserService/Login")
	assert.Contains(t, AuthIgnoreMethods, "/user.UserService/Register")
}

// Exercises the godotenv.Load error branch by pointing at a missing file.
// LoadConfig must still succeed because env.Parse can populate the schema
// from environment variables alone.
func TestLoadConfig_GodotenvLoadError(t *testing.T) {
	t.Setenv("CONFIG_FILE", "/nonexistent/path/config.yaml")
	cfg := LoadConfig()
	assert.NotNil(t, cfg)
}
