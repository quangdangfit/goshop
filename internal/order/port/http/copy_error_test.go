package http

import (
	"bytes"
	"encoding/json"
	"math"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/quangdangfit/gocommon/logger"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"goshop/internal/order/domain"
	"goshop/internal/order/model"
	srvMocks "goshop/internal/order/service/mocks"
	"goshop/pkg/config"
	"goshop/pkg/paging"
)

func nanOrder() *model.Order {
	return &model.Order{ID: "o1", TotalPrice: math.NaN(), Status: model.OrderStatusNew}
}

func setupOrderHandler(t *testing.T) (*OrderHandler, *srvMocks.OrderService) {
	logger.Initialize(config.ProductionEnv)
	gin.SetMode(gin.TestMode)
	svc := srvMocks.NewOrderService(t)
	return NewOrderHandler(svc), svc
}

func TestPlaceOrder_CopyError(t *testing.T) {
	h, svc := setupOrderHandler(t)
	svc.On("PlaceOrder", mock.Anything, mock.Anything).Return(nanOrder(), nil).Once()

	body, _ := json.Marshal(domain.PlaceOrderReq{
		Lines: []domain.PlaceOrderLineReq{{ProductID: "p1", Quantity: 1}},
	})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("userId", "u1")
	c.Request = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))
	h.PlaceOrder(c)
	require.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestGetOrders_CopyError(t *testing.T) {
	h, svc := setupOrderHandler(t)
	svc.On("GetMyOrders", mock.Anything, mock.Anything).
		Return([]*model.Order{nanOrder()}, &paging.Pagination{}, nil).Once()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("userId", "u1")
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	h.GetOrders(c)
	require.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestGetOrderByID_CopyError(t *testing.T) {
	h, svc := setupOrderHandler(t)
	svc.On("GetOrderByID", mock.Anything, "o1").Return(nanOrder(), nil).Once()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("userId", "u1")
	c.Params = gin.Params{{Key: "id", Value: "o1"}}
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	h.GetOrderByID(c)
	require.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestUpdateOrderStatus_CopyError(t *testing.T) {
	h, svc := setupOrderHandler(t)
	svc.On("UpdateOrderStatus", mock.Anything, "o1", model.OrderStatusInProgress).
		Return(nanOrder(), nil).Once()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "o1"}}
	c.Request = httptest.NewRequest(http.MethodPut, "/?status=in-progress", nil)
	h.UpdateOrderStatus(c)
	require.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestCancelOrder_CopyError(t *testing.T) {
	h, svc := setupOrderHandler(t)
	svc.On("CancelOrder", mock.Anything, "o1", "u1").Return(nanOrder(), nil).Once()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("userId", "u1")
	c.Params = gin.Params{{Key: "id", Value: "o1"}}
	c.Request = httptest.NewRequest(http.MethodDelete, "/", nil)
	h.CancelOrder(c)
	require.Equal(t, http.StatusInternalServerError, w.Code)
}
