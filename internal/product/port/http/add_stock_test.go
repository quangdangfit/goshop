package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/quangdangfit/gocommon/logger"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"goshop/internal/product/model"
	srvMocks "goshop/internal/product/service/mocks"
	"goshop/pkg/config"
	redisMocks "goshop/pkg/redis/mocks"
)

func newAddStockHandler(t *testing.T) (*ProductHandler, *srvMocks.ProductService, *redisMocks.Redis) {
	logger.Initialize(config.ProductionEnv)
	gin.SetMode(gin.TestMode)
	svc := srvMocks.NewProductService(t)
	cache := redisMocks.NewRedis(t)
	return NewProductHandler(cache, svc), svc, cache
}

func TestAddStock_BadBody(t *testing.T) {
	h, _, _ := newAddStockHandler(t)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "p1"}}
	c.Request = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte("not-json")))
	h.AddStock(c)
	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAddStock_InvalidQuantity(t *testing.T) {
	h, _, _ := newAddStockHandler(t)
	body, _ := json.Marshal(map[string]any{"quantity": 0})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "p1"}}
	c.Request = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))
	h.AddStock(c)
	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAddStock_ServiceError(t *testing.T) {
	h, svc, _ := newAddStockHandler(t)
	svc.On("AddStock", mock.Anything, "p1", 10, "admin1").
		Return(nil, errors.New("db down")).Once()
	body, _ := json.Marshal(map[string]any{"quantity": 10})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("userId", "admin1")
	c.Params = gin.Params{{Key: "id", Value: "p1"}}
	c.Request = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))
	h.AddStock(c)
	require.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestAddStock_Success(t *testing.T) {
	h, svc, cache := newAddStockHandler(t)
	svc.On("AddStock", mock.Anything, "p1", 10, "admin1").
		Return(&model.Product{ID: "p1", StockQuantity: 110}, nil).Once()
	cache.On("RemovePattern", "*product*").Return(nil).Once()
	body, _ := json.Marshal(map[string]any{"quantity": 10})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("userId", "admin1")
	c.Params = gin.Params{{Key: "id", Value: "p1"}}
	c.Request = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))
	h.AddStock(c)
	require.Equal(t, http.StatusOK, w.Code)
}
