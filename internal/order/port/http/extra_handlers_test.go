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

	"goshop/internal/order/domain"
	"goshop/internal/order/model"
	"goshop/internal/order/service"
	srvMocks "goshop/internal/order/service/mocks"
	"goshop/pkg/config"
)

func TestPlaceOrder_UnauthorizedNoUserID(t *testing.T) {
	logger.Initialize(config.ProductionEnv)
	gin.SetMode(gin.TestMode)
	mockSvc := srvMocks.NewOrderService(t)
	h := NewOrderHandler(mockSvc)

	body, _ := json.Marshal(domain.PlaceOrderReq{Lines: []domain.PlaceOrderLineReq{{ProductID: "p1", Quantity: 1}}})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))
	h.PlaceOrder(c)
	require.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestPlaceOrder_InsufficientStock409(t *testing.T) {
	logger.Initialize(config.ProductionEnv)
	gin.SetMode(gin.TestMode)
	mockSvc := srvMocks.NewOrderService(t)
	mockSvc.On("PlaceOrder", mock.Anything, mock.Anything).
		Return(nil, &service.InsufficientStockError{ProductID: "p1", Requested: 5}).Once()
	h := NewOrderHandler(mockSvc)

	body, _ := json.Marshal(domain.PlaceOrderReq{Lines: []domain.PlaceOrderLineReq{{ProductID: "p1", Quantity: 5}}})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("userId", "u1")
	c.Request = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))
	h.PlaceOrder(c)
	require.Equal(t, http.StatusConflict, w.Code)
}

func TestGetOrders_Unauthorized(t *testing.T) {
	logger.Initialize(config.ProductionEnv)
	gin.SetMode(gin.TestMode)
	h := NewOrderHandler(srvMocks.NewOrderService(t))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	h.GetOrders(c)
	require.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestGetOrderByID_Unauthorized(t *testing.T) {
	logger.Initialize(config.ProductionEnv)
	gin.SetMode(gin.TestMode)
	h := NewOrderHandler(srvMocks.NewOrderService(t))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	h.GetOrderByID(c)
	require.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestUpdateOrderStatus_MissingID(t *testing.T) {
	logger.Initialize(config.ProductionEnv)
	gin.SetMode(gin.TestMode)
	h := NewOrderHandler(srvMocks.NewOrderService(t))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPut, "/?status=paid", nil)
	h.UpdateOrderStatus(c)
	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateOrderStatus_InvalidStatus(t *testing.T) {
	logger.Initialize(config.ProductionEnv)
	gin.SetMode(gin.TestMode)
	h := NewOrderHandler(srvMocks.NewOrderService(t))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "o1"}}
	c.Request = httptest.NewRequest(http.MethodPut, "/?status=bogus", nil)
	h.UpdateOrderStatus(c)
	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateOrderStatus_OK(t *testing.T) {
	logger.Initialize(config.ProductionEnv)
	gin.SetMode(gin.TestMode)
	mockSvc := srvMocks.NewOrderService(t)
	mockSvc.On("UpdateOrderStatus", mock.Anything, "o1", model.OrderStatusInProgress).
		Return(&model.Order{ID: "o1", Status: model.OrderStatusInProgress}, nil).Once()
	h := NewOrderHandler(mockSvc)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "o1"}}
	c.Request = httptest.NewRequest(http.MethodPut, "/?status=in-progress", nil)
	h.UpdateOrderStatus(c)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestCancelOrder_Unauthorized(t *testing.T) {
	logger.Initialize(config.ProductionEnv)
	gin.SetMode(gin.TestMode)
	h := NewOrderHandler(srvMocks.NewOrderService(t))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodDelete, "/", nil)
	h.CancelOrder(c)
	require.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestCreateCoupon_BadBody(t *testing.T) {
	logger.Initialize(config.ProductionEnv)
	gin.SetMode(gin.TestMode)
	h := NewCouponHandler(srvMocks.NewCouponService(t))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte("not-json")))
	h.CreateCoupon(c)
	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetCouponByCode_NotFound(t *testing.T) {
	logger.Initialize(config.ProductionEnv)
	gin.SetMode(gin.TestMode)
	mockSvc := srvMocks.NewCouponService(t)
	mockSvc.On("GetByCode", mock.Anything, "BADCODE").Return(nil, errors.New("not found")).Once()
	h := NewCouponHandler(mockSvc)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "code", Value: "BADCODE"}}
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	h.GetCouponByCode(c)
	require.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetOrderByID_MissingID(t *testing.T) {
	logger.Initialize(config.ProductionEnv)
	gin.SetMode(gin.TestMode)
	h := NewOrderHandler(srvMocks.NewOrderService(t))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("userId", "u1")
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	h.GetOrderByID(c)
	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCancelOrder_MissingID(t *testing.T) {
	logger.Initialize(config.ProductionEnv)
	gin.SetMode(gin.TestMode)
	h := NewOrderHandler(srvMocks.NewOrderService(t))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("userId", "u1")
	c.Request = httptest.NewRequest(http.MethodDelete, "/", nil)
	h.CancelOrder(c)
	require.Equal(t, http.StatusBadRequest, w.Code)
}
