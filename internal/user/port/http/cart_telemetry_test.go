package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/quangdangfit/gocommon/logger"
	"github.com/stretchr/testify/require"

	"goshop/pkg/config"
)

func setupCartTelemetryCtx(userID string, body any) (*gin.Context, *httptest.ResponseRecorder) {
	logger.Initialize(config.ProductionEnv)
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	if userID != "" {
		c.Set("userId", userID)
	}
	var buf []byte
	switch v := body.(type) {
	case string:
		buf = []byte(v)
	default:
		buf, _ = json.Marshal(body)
	}
	c.Request = httptest.NewRequest(http.MethodPut, "/", bytes.NewBuffer(buf))
	return c, w
}

func TestPutCartSnapshot_Unauthorized(t *testing.T) {
	c, w := setupCartTelemetryCtx("", map[string]any{"items": []any{}})
	PutCartSnapshot(c)
	require.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestPutCartSnapshot_BadBody(t *testing.T) {
	c, w := setupCartTelemetryCtx("u1", "not-json")
	PutCartSnapshot(c)
	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPutCartSnapshot_TooManyItems(t *testing.T) {
	items := make([]map[string]any, 101)
	for i := range items {
		items[i] = map[string]any{"product_id": "p", "quantity": 1}
	}
	c, w := setupCartTelemetryCtx("u1", map[string]any{"items": items})
	PutCartSnapshot(c)
	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPutCartSnapshot_Success(t *testing.T) {
	body := strings.NewReader(`{"items":[{"product_id":"p1","quantity":2}]}`)
	logger.Initialize(config.ProductionEnv)
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("userId", "u1")
	c.Request = httptest.NewRequest(http.MethodPut, "/", body)
	PutCartSnapshot(c)
	require.Equal(t, http.StatusNoContent, w.Code)
}
