package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"goshop/pkg/config"
	"goshop/pkg/jtoken"
	redisMocks "goshop/pkg/redis/mocks"
)

func setupGinTest() {
	gin.SetMode(gin.TestMode)
}

// JWT / JWTAuth / JWTRefresh
// =================================================================

func TestJWTAuth_NoToken(t *testing.T) {
	setupGinTest()
	w := httptest.NewRecorder()
	c, engine := gin.CreateTestContext(w)
	engine.Use(JWTAuth())
	engine.GET("/test", func(c *gin.Context) { c.Status(http.StatusOK) })

	c.Request, _ = http.NewRequest(http.MethodGet, "/test", nil)
	engine.ServeHTTP(w, c.Request)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestJWTAuth_ValidToken(t *testing.T) {
	config.LoadConfig()
	setupGinTest()

	payload := map[string]interface{}{
		"id":    "user-1",
		"email": "test@example.com",
		"role":  "customer",
	}
	token := jtoken.GenerateAccessToken(payload)

	w := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(w)
	engine.Use(JWTAuth())
	engine.GET("/test", func(c *gin.Context) { c.Status(http.StatusOK) })

	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", token)
	engine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestJWTAuth_InvalidToken(t *testing.T) {
	setupGinTest()
	w := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(w)
	engine.Use(JWTAuth())
	engine.GET("/test", func(c *gin.Context) { c.Status(http.StatusOK) })

	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "invalid.token.value")
	engine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestJWTAuth_WrongTokenType(t *testing.T) {
	config.LoadConfig()
	setupGinTest()

	// Generate a refresh token but use it where access token is expected
	payload := map[string]interface{}{
		"id":    "user-1",
		"email": "test@example.com",
		"role":  "customer",
	}
	token := jtoken.GenerateRefreshToken(payload)

	w := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(w)
	engine.Use(JWTAuth()) // expects access token
	engine.GET("/test", func(c *gin.Context) { c.Status(http.StatusOK) })

	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", token)
	engine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestJWTRefresh_ValidToken(t *testing.T) {
	config.LoadConfig()
	setupGinTest()

	payload := map[string]interface{}{
		"id":    "user-1",
		"email": "test@example.com",
		"role":  "customer",
	}
	token := jtoken.GenerateRefreshToken(payload)

	w := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(w)
	engine.Use(JWTRefresh())
	engine.GET("/test", func(c *gin.Context) { c.Status(http.StatusOK) })

	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", token)
	engine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// CORS
// =================================================================

func TestCORS_WildcardOrigin(t *testing.T) {
	config.LoadConfig()
	config.GetConfig().CORSAllowedOrigins = "*"
	setupGinTest()

	w := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(w)
	engine.Use(CORS())
	engine.GET("/test", func(c *gin.Context) { c.Status(http.StatusOK) })

	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Origin", "http://example.com")
	engine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "http://example.com", w.Header().Get("Access-Control-Allow-Origin"))
}

func TestCORS_SpecificOriginAllowed(t *testing.T) {
	config.LoadConfig()
	config.GetConfig().CORSAllowedOrigins = "http://example.com"
	setupGinTest()

	w := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(w)
	engine.Use(CORS())
	engine.GET("/test", func(c *gin.Context) { c.Status(http.StatusOK) })

	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Origin", "http://example.com")
	engine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCORS_OriginNotAllowed(t *testing.T) {
	config.LoadConfig()
	config.GetConfig().CORSAllowedOrigins = "http://allowed.com"
	setupGinTest()

	w := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(w)
	engine.Use(CORS())
	engine.GET("/test", func(c *gin.Context) { c.Status(http.StatusOK) })

	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Origin", "http://notallowed.com")
	engine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "http://allowed.com", w.Header().Get("Access-Control-Allow-Origin"))
}

func TestCORS_OptionsMethod(t *testing.T) {
	config.LoadConfig()
	config.GetConfig().CORSAllowedOrigins = "*"
	setupGinTest()

	w := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(w)
	engine.Use(CORS())
	engine.OPTIONS("/test", func(c *gin.Context) { c.Status(http.StatusOK) })

	req, _ := http.NewRequest(http.MethodOptions, "/test", nil)
	req.Header.Set("Origin", "http://example.com")
	engine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

// RateLimit
// =================================================================

func TestRateLimit_Disabled(t *testing.T) {
	config.LoadConfig()
	config.GetConfig().RateLimitRequests = 0
	setupGinTest()

	mockRedis := redisMocks.NewRedis(t)

	w := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(w)
	engine.Use(RateLimit(mockRedis))
	engine.GET("/test", func(c *gin.Context) { c.Status(http.StatusOK) })

	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	engine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRateLimit_Allowed(t *testing.T) {
	config.LoadConfig()
	config.GetConfig().RateLimitRequests = 100
	config.GetConfig().RateLimitWindowSeconds = 60
	setupGinTest()

	mockRedis := redisMocks.NewRedis(t)
	mockRedis.On("Incr", mock.Anything, mock.Anything).Return(int64(1), nil).Times(1)

	w := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(w)
	engine.Use(RateLimit(mockRedis))
	engine.GET("/test", func(c *gin.Context) { c.Status(http.StatusOK) })

	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	engine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRateLimit_Exceeded(t *testing.T) {
	config.LoadConfig()
	config.GetConfig().RateLimitRequests = 1
	config.GetConfig().RateLimitWindowSeconds = 60
	setupGinTest()

	mockRedis := redisMocks.NewRedis(t)
	mockRedis.On("Incr", mock.Anything, mock.Anything).Return(int64(2), nil).Times(1)

	w := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(w)
	engine.Use(RateLimit(mockRedis))
	engine.GET("/test", func(c *gin.Context) { c.Status(http.StatusOK) })

	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	engine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusTooManyRequests, w.Code)
}

func TestRateLimit_RedisError(t *testing.T) {
	config.LoadConfig()
	config.GetConfig().RateLimitRequests = 100
	config.GetConfig().RateLimitWindowSeconds = 60
	setupGinTest()

	mockRedis := redisMocks.NewRedis(t)
	mockRedis.On("Incr", mock.Anything, mock.Anything).Return(int64(0), errors.New("redis error")).Times(1)

	w := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(w)
	engine.Use(RateLimit(mockRedis))
	engine.GET("/test", func(c *gin.Context) { c.Status(http.StatusOK) })

	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	engine.ServeHTTP(w, req)

	// On Redis error, request is allowed through
	assert.Equal(t, http.StatusOK, w.Code)
}
