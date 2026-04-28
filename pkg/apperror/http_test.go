package apperror

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"goshop/pkg/config"
)

func init() {
	gin.SetMode(gin.TestMode)
}

type httpErrorResponse struct {
	Result interface{}            `json:"result"`
	Error  map[string]interface{} `json:"error"`
}

func newTestContext() (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	return c, w
}

func decodeHTTPError(t *testing.T, body []byte) httpErrorResponse {
	t.Helper()
	var resp httpErrorResponse
	assert.NoError(t, json.Unmarshal(body, &resp))
	return resp
}

func withEnvironment(env string) func() {
	cfg := config.GetConfig()
	orig := cfg.Environment
	cfg.Environment = env
	return func() { cfg.Environment = orig }
}

func TestAppError_HTTPError_NonProductionIncludesDebug(t *testing.T) {
	restore := withEnvironment("development")
	defer restore()

	c, w := newTestContext()
	appErr := Wrap(ErrNotFound, fmt.Errorf("record missing"))
	appErr.HTTPError(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
	resp := decodeHTTPError(t, w.Body.Bytes())
	assert.Nil(t, resp.Result)
	assert.Equal(t, "NOT_FOUND", resp.Error["code"])
	assert.Equal(t, "Resource not found", resp.Error["message"])
	assert.Equal(t, "Resource not found: record missing", resp.Error["debug"])
}

func TestAppError_HTTPError_ProductionOmitsDebug(t *testing.T) {
	restore := withEnvironment(config.ProductionEnv)
	defer restore()

	c, w := newTestContext()
	appErr := Wrap(ErrInternal, fmt.Errorf("secret db error"))
	appErr.HTTPError(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	resp := decodeHTTPError(t, w.Body.Bytes())
	assert.Equal(t, "INTERNAL", resp.Error["code"])
	assert.Equal(t, "Something went wrong", resp.Error["message"])
	_, hasDebug := resp.Error["debug"]
	assert.False(t, hasDebug, "debug field must be hidden in production")
}

func TestToHTTPError_UnwrapsAppError(t *testing.T) {
	restore := withEnvironment("development")
	defer restore()

	c, w := newTestContext()
	wrapped := fmt.Errorf("handler failed: %w", Wrap(ErrForbidden, errors.New("not owner")))

	ToHTTPError(c, wrapped, http.StatusTeapot, "fallback should not be used")

	assert.Equal(t, http.StatusForbidden, w.Code)
	resp := decodeHTTPError(t, w.Body.Bytes())
	assert.Equal(t, "FORBIDDEN", resp.Error["code"])
	assert.Equal(t, "Permission denied", resp.Error["message"])
}

func TestToHTTPError_NonAppErrorFallback_NonProduction(t *testing.T) {
	restore := withEnvironment("development")
	defer restore()

	c, w := newTestContext()
	raw := errors.New("boom")

	ToHTTPError(c, raw, http.StatusBadGateway, "upstream unavailable")

	assert.Equal(t, http.StatusBadGateway, w.Code)
	resp := decodeHTTPError(t, w.Body.Bytes())
	assert.Equal(t, "upstream unavailable", resp.Error["message"])
	assert.Equal(t, "boom", resp.Error["debug"])
	_, hasCode := resp.Error["code"]
	assert.False(t, hasCode, "non-AppError fallback must not expose a code field")
}

func TestToHTTPError_NonAppErrorFallback_Production(t *testing.T) {
	restore := withEnvironment(config.ProductionEnv)
	defer restore()

	c, w := newTestContext()
	ToHTTPError(c, errors.New("internal plumbing"), http.StatusServiceUnavailable, "try again later")

	assert.Equal(t, http.StatusServiceUnavailable, w.Code)
	resp := decodeHTTPError(t, w.Body.Bytes())
	assert.Equal(t, "try again later", resp.Error["message"])
	_, hasDebug := resp.Error["debug"]
	assert.False(t, hasDebug, "debug field must be hidden in production")
}
