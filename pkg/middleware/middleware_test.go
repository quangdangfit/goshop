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

func TestJWTAuth(t *testing.T) {
	tests := []struct {
		name     string
		setup    func()
		token    string
		expected int
	}{
		{
			name:     "NoToken",
			setup:    func() {},
			token:    "",
			expected: http.StatusUnauthorized,
		},
		{
			name: "ValidToken",
			setup: func() {
				config.LoadConfig()
			},
			token: func() string {
				config.LoadConfig()
				return jtoken.GenerateAccessToken(map[string]interface{}{
					"id": "user-1", "email": "test@example.com", "role": "customer",
				})
			}(),
			expected: http.StatusOK,
		},
		{
			name:     "InvalidToken",
			setup:    func() {},
			token:    "invalid.token.value",
			expected: http.StatusUnauthorized,
		},
		{
			name: "WrongTokenType",
			setup: func() {
				config.LoadConfig()
			},
			token: func() string {
				config.LoadConfig()
				return jtoken.GenerateRefreshToken(map[string]interface{}{
					"id": "user-1", "email": "test@example.com", "role": "customer",
				})
			}(),
			expected: http.StatusUnauthorized,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()
			setupGinTest()

			w := httptest.NewRecorder()
			_, engine := gin.CreateTestContext(w)
			engine.Use(JWTAuth())
			engine.GET("/test", func(c *gin.Context) { c.Status(http.StatusOK) })

			req, _ := http.NewRequest(http.MethodGet, "/test", nil)
			if tc.token != "" {
				req.Header.Set("Authorization", tc.token)
			}
			engine.ServeHTTP(w, req)

			assert.Equal(t, tc.expected, w.Code)
		})
	}
}

func TestJWTRefresh(t *testing.T) {
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

func TestAdminOnly(t *testing.T) {
	tests := []struct {
		name     string
		role     string
		expected int
	}{
		{
			name:     "AdminAllowed",
			role:     "admin",
			expected: http.StatusOK,
		},
		{
			name:     "CustomerForbidden",
			role:     "customer",
			expected: http.StatusForbidden,
		},
		{
			name:     "EmptyRoleForbidden",
			role:     "",
			expected: http.StatusForbidden,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			config.LoadConfig()
			setupGinTest()

			w := httptest.NewRecorder()
			_, engine := gin.CreateTestContext(w)
			engine.Use(func(c *gin.Context) {
				c.Set("role", tc.role)
				c.Next()
			})
			engine.Use(AdminOnly())
			engine.GET("/test", func(c *gin.Context) { c.Status(http.StatusOK) })

			req, _ := http.NewRequest(http.MethodGet, "/test", nil)
			engine.ServeHTTP(w, req)

			assert.Equal(t, tc.expected, w.Code)
		})
	}
}

func TestCORS(t *testing.T) {
	tests := []struct {
		name           string
		allowedOrigins string
		origin         string
		method         string
		expectedStatus int
		checkHeader    func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name:           "WildcardOrigin",
			allowedOrigins: "*",
			origin:         "http://example.com",
			method:         http.MethodGet,
			expectedStatus: http.StatusOK,
			checkHeader: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.Equal(t, "http://example.com", w.Header().Get("Access-Control-Allow-Origin"))
			},
		},
		{
			name:           "SpecificOriginAllowed",
			allowedOrigins: "http://example.com",
			origin:         "http://example.com",
			method:         http.MethodGet,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "OriginNotAllowed",
			allowedOrigins: "http://allowed.com",
			origin:         "http://notallowed.com",
			method:         http.MethodGet,
			expectedStatus: http.StatusOK,
			checkHeader: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.Equal(t, "http://allowed.com", w.Header().Get("Access-Control-Allow-Origin"))
			},
		},
		{
			name:           "OptionsMethod",
			allowedOrigins: "*",
			origin:         "http://example.com",
			method:         http.MethodOptions,
			expectedStatus: http.StatusNoContent,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			config.LoadConfig()
			config.GetConfig().CORSAllowedOrigins = tc.allowedOrigins
			setupGinTest()

			w := httptest.NewRecorder()
			_, engine := gin.CreateTestContext(w)
			engine.Use(CORS())
			if tc.method == http.MethodOptions {
				engine.OPTIONS("/test", func(c *gin.Context) { c.Status(http.StatusOK) })
			} else {
				engine.GET("/test", func(c *gin.Context) { c.Status(http.StatusOK) })
			}

			req, _ := http.NewRequest(tc.method, "/test", nil)
			req.Header.Set("Origin", tc.origin)
			engine.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatus, w.Code)
			if tc.checkHeader != nil {
				tc.checkHeader(t, w)
			}
		})
	}
}

func TestRateLimit(t *testing.T) {
	tests := []struct {
		name          string
		maxRequests   int
		windowSeconds int
		setupMock     func(m *redisMocks.Redis)
		expected      int
	}{
		{
			name:        "Disabled",
			maxRequests: 0,
			setupMock:   func(m *redisMocks.Redis) {},
			expected:    http.StatusOK,
		},
		{
			name:          "Allowed",
			maxRequests:   100,
			windowSeconds: 60,
			setupMock: func(m *redisMocks.Redis) {
				m.On("Incr", mock.Anything, mock.Anything).Return(int64(1), nil).Times(1)
			},
			expected: http.StatusOK,
		},
		{
			name:          "Exceeded",
			maxRequests:   1,
			windowSeconds: 60,
			setupMock: func(m *redisMocks.Redis) {
				m.On("Incr", mock.Anything, mock.Anything).Return(int64(2), nil).Times(1)
			},
			expected: http.StatusTooManyRequests,
		},
		{
			name:          "RedisError",
			maxRequests:   100,
			windowSeconds: 60,
			setupMock: func(m *redisMocks.Redis) {
				m.On("Incr", mock.Anything, mock.Anything).Return(int64(0), errors.New("redis error")).Times(1)
			},
			expected: http.StatusOK,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			config.LoadConfig()
			config.GetConfig().RateLimitRequests = tc.maxRequests
			if tc.windowSeconds > 0 {
				config.GetConfig().RateLimitWindowSeconds = tc.windowSeconds
			}
			setupGinTest()

			mockRedis := redisMocks.NewRedis(t)
			tc.setupMock(mockRedis)

			w := httptest.NewRecorder()
			_, engine := gin.CreateTestContext(w)
			engine.Use(RateLimit(mockRedis))
			engine.GET("/test", func(c *gin.Context) { c.Status(http.StatusOK) })

			req, _ := http.NewRequest(http.MethodGet, "/test", nil)
			engine.ServeHTTP(w, req)

			assert.Equal(t, tc.expected, w.Code)
		})
	}
}
