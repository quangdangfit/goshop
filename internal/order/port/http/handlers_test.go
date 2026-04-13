package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/quangdangfit/gocommon/logger"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"goshop/internal/order/domain"
	"goshop/internal/order/model"
	"goshop/internal/order/service/mocks"
	productMocks "goshop/internal/product/service/mocks"
	"goshop/pkg/config"
	"goshop/pkg/paging"
	"goshop/pkg/response"
	"goshop/pkg/utils"
)

type OrderHandlerTestSuite struct {
	suite.Suite
	mockService        *mocks.OrderService
	mockProductService *productMocks.ProductService
	handler            *OrderHandler
}

func (suite *OrderHandlerTestSuite) SetupTest() {
	logger.Initialize(config.ProductionEnv)

	suite.mockService = mocks.NewOrderService(suite.T())
	suite.mockProductService = productMocks.NewProductService(suite.T())
	suite.handler = NewOrderHandler(suite.mockService)
}

func TestOrderHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(OrderHandlerTestSuite))
}

func (suite *OrderHandlerTestSuite) prepareContext(body any) (*gin.Context, *httptest.ResponseRecorder) {
	requestBody, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("", "/", bytes.NewBuffer(requestBody))
	c, _ := gin.CreateTestContext(w)
	c.Request = r

	return c, w
}

func (suite *OrderHandlerTestSuite) TestPlaceOrder() {
	tests := []struct {
		name      string
		body      any
		setup     func(ctx *gin.Context)
		expected  int
		checkBody func(writer *httptest.ResponseRecorder)
	}{
		{
			name: "Success",
			body: &domain.PlaceOrderReq{
				Lines: []domain.PlaceOrderLineReq{
					{ProductID: "productId1", Quantity: 2},
					{ProductID: "productId2", Quantity: 3},
				},
			},
			setup: func(ctx *gin.Context) {
				ctx.Set("userId", "123456")
				suite.mockService.On("PlaceOrder", mock.Anything, &domain.PlaceOrderReq{
					UserID: "123456",
					Lines: []domain.PlaceOrderLineReq{
						{ProductID: "productId1", Quantity: 2},
						{ProductID: "productId2", Quantity: 3},
					},
				}).Return(
					&model.Order{
						ID:         "orderId1",
						Code:       "orderCode1",
						TotalPrice: 8,
						Status:     model.OrderStatusNew,
						Lines: []*model.OrderLine{
							{ProductID: "productId1", Quantity: 2},
							{ProductID: "productId2", Quantity: 3},
						},
					},
					nil,
				).Times(1)
			},
			expected: http.StatusOK,
			checkBody: func(writer *httptest.ResponseRecorder) {
				var res response.Response
				var orderRes domain.Order
				_ = json.Unmarshal(writer.Body.Bytes(), &res)
				_ = utils.Copy(&orderRes, &res.Result)
				suite.Equal(float64(8), orderRes.TotalPrice)
				suite.Equal(string(model.OrderStatusNew), orderRes.Status)
				suite.Equal(2, len(orderRes.Lines))
			},
		},
		{
			name: "InvalidProductIdType",
			body: map[string]interface{}{
				"lines": []map[string]interface{}{
					{"product_id": 1, "quantity": 2},
				},
			},
			setup:    func(ctx *gin.Context) {},
			expected: http.StatusBadRequest,
			checkBody: func(writer *httptest.ResponseRecorder) {
				var res map[string]map[string]string
				_ = json.Unmarshal(writer.Body.Bytes(), &res)
				suite.Equal("Invalid request parameters", res["error"]["message"])
			},
		},
		{
			name: "InvalidQuantityType",
			body: map[string]interface{}{
				"lines": []map[string]interface{}{
					{"product_id": "productId1", "quantity": "1"},
				},
			},
			setup:    func(ctx *gin.Context) {},
			expected: http.StatusBadRequest,
			checkBody: func(writer *httptest.ResponseRecorder) {
				var res map[string]map[string]string
				_ = json.Unmarshal(writer.Body.Bytes(), &res)
				suite.Equal("Invalid request parameters", res["error"]["message"])
			},
		},
		{
			name: "Unauthorized",
			body: &domain.PlaceOrderReq{
				Lines: []domain.PlaceOrderLineReq{
					{ProductID: "productId1", Quantity: 2},
					{ProductID: "productId2", Quantity: 3},
				},
			},
			setup:    func(ctx *gin.Context) {},
			expected: http.StatusUnauthorized,
			checkBody: func(writer *httptest.ResponseRecorder) {
				var res map[string]map[string]string
				_ = json.Unmarshal(writer.Body.Bytes(), &res)
				suite.Equal("Unauthorized", res["error"]["message"])
			},
		},
		{
			name: "Fail",
			body: &domain.PlaceOrderReq{
				Lines: []domain.PlaceOrderLineReq{
					{ProductID: "productId1", Quantity: 2},
					{ProductID: "productId2", Quantity: 3},
				},
			},
			setup: func(ctx *gin.Context) {
				ctx.Set("userId", "123456")
				suite.mockService.On("PlaceOrder", mock.Anything, &domain.PlaceOrderReq{
					UserID: "123456",
					Lines: []domain.PlaceOrderLineReq{
						{ProductID: "productId1", Quantity: 2},
						{ProductID: "productId2", Quantity: 3},
					},
				}).Return(nil, errors.New("error")).Times(1)
			},
			expected: http.StatusInternalServerError,
			checkBody: func(writer *httptest.ResponseRecorder) {
				var res map[string]map[string]string
				_ = json.Unmarshal(writer.Body.Bytes(), &res)
				suite.Equal("Something went wrong", res["error"]["message"])
			},
		},
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			ctx, writer := suite.prepareContext(tc.body)
			tc.setup(ctx)
			suite.handler.PlaceOrder(ctx)
			suite.Equal(tc.expected, writer.Code)
			if tc.checkBody != nil {
				tc.checkBody(writer)
			}
		})
	}
}

func (suite *OrderHandlerTestSuite) TestGetOrderByID() {
	tests := []struct {
		name      string
		setup     func(ctx *gin.Context)
		expected  int
		checkBody func(writer *httptest.ResponseRecorder)
	}{
		{
			name: "Success",
			setup: func(ctx *gin.Context) {
				ctx.Set("userId", "123456")
				ctx.AddParam("id", "orderId1")
				suite.mockService.On("GetOrderByID", mock.Anything, "orderId1").
					Return(
						&model.Order{
							ID:         "orderId1",
							UserID:     "123456",
							TotalPrice: 5,
							Status:     model.OrderStatusNew,
						},
						nil,
					).Times(1)
			},
			expected: http.StatusOK,
			checkBody: func(writer *httptest.ResponseRecorder) {
				var res response.Response
				var orderRes domain.Order
				_ = json.Unmarshal(writer.Body.Bytes(), &res)
				_ = utils.Copy(&orderRes, &res.Result)
				suite.Equal(float64(5), orderRes.TotalPrice)
				suite.Equal(string(model.OrderStatusNew), orderRes.Status)
				suite.Equal(0, len(orderRes.Lines))
			},
		},
		{
			name: "MissID",
			setup: func(ctx *gin.Context) {
				ctx.Set("userId", "123456")
			},
			expected: http.StatusBadRequest,
			checkBody: func(writer *httptest.ResponseRecorder) {
				var res response.Response
				_ = json.Unmarshal(writer.Body.Bytes(), &res)
				suite.NotNil(res.Error)
			},
		},
		{
			name: "Unauthorized",
			setup: func(ctx *gin.Context) {
				ctx.AddParam("id", "orderId1")
			},
			expected: http.StatusUnauthorized,
			checkBody: func(writer *httptest.ResponseRecorder) {
				var res response.Response
				_ = json.Unmarshal(writer.Body.Bytes(), &res)
				suite.NotNil(res.Error)
			},
		},
		{
			name: "Fail",
			setup: func(ctx *gin.Context) {
				ctx.Set("userId", "123456")
				ctx.AddParam("id", "orderId1")
				suite.mockService.On("GetOrderByID", mock.Anything, "orderId1").
					Return(nil, errors.New("error")).Times(1)
			},
			expected: http.StatusNotFound,
			checkBody: func(writer *httptest.ResponseRecorder) {
				var res response.Response
				_ = json.Unmarshal(writer.Body.Bytes(), &res)
				suite.NotNil(res.Error)
			},
		},
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			ctx, writer := suite.prepareContext(nil)
			tc.setup(ctx)
			suite.handler.GetOrderByID(ctx)
			suite.Equal(tc.expected, writer.Code)
			if tc.checkBody != nil {
				tc.checkBody(writer)
			}
		})
	}
}

func (suite *OrderHandlerTestSuite) TestGetOrders() {
	tests := []struct {
		name      string
		setup     func(ctx *gin.Context)
		expected  int
		checkBody func(writer *httptest.ResponseRecorder)
	}{
		{
			name: "Success",
			setup: func(ctx *gin.Context) {
				ctx.Set("userId", "123456")
				suite.mockService.On("GetMyOrders", mock.Anything, mock.Anything).
					Return(
						[]*model.Order{
							{
								ID:         "orderId1",
								UserID:     "123456",
								TotalPrice: 5,
								Status:     model.OrderStatusNew,
							},
						},
						&paging.Pagination{
							Total:       1,
							CurrentPage: 1,
						},
						nil,
					).Times(1)
			},
			expected: http.StatusOK,
			checkBody: func(writer *httptest.ResponseRecorder) {
				var res response.Response
				var orderRes domain.ListOrderRes
				_ = json.Unmarshal(writer.Body.Bytes(), &res)
				_ = utils.Copy(&orderRes, &res.Result)
				suite.Equal(1, len(orderRes.Orders))
				suite.Equal("orderId1", orderRes.Orders[0].ID)
				suite.Equal(float64(5), orderRes.Orders[0].TotalPrice)
				suite.Equal(string(model.OrderStatusNew), orderRes.Orders[0].Status)
			},
		},
		{
			name:     "Unauthorized",
			setup:    func(ctx *gin.Context) {},
			expected: http.StatusUnauthorized,
			checkBody: func(writer *httptest.ResponseRecorder) {
				var res response.Response
				_ = json.Unmarshal(writer.Body.Bytes(), &res)
				suite.NotNil(res.Error)
			},
		},
		{
			name: "InvalidFieldType",
			setup: func(ctx *gin.Context) {
				ctx.Request.URL, _ = url.Parse("?page=q")
			},
			expected: http.StatusBadRequest,
			checkBody: func(writer *httptest.ResponseRecorder) {
				var res response.Response
				_ = json.Unmarshal(writer.Body.Bytes(), &res)
				suite.NotNil(res.Error)
			},
		},
		{
			name: "Fail",
			setup: func(ctx *gin.Context) {
				ctx.Set("userId", "123456")
				suite.mockService.On("GetMyOrders", mock.Anything, mock.Anything).
					Return(nil, nil, errors.New("error")).Times(1)
			},
			expected: http.StatusInternalServerError,
			checkBody: func(writer *httptest.ResponseRecorder) {
				var res response.Response
				_ = json.Unmarshal(writer.Body.Bytes(), &res)
				suite.NotNil(res.Error)
			},
		},
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			ctx, writer := suite.prepareContext(nil)
			tc.setup(ctx)
			suite.handler.GetOrders(ctx)
			suite.Equal(tc.expected, writer.Code)
			if tc.checkBody != nil {
				tc.checkBody(writer)
			}
		})
	}
}

func (suite *OrderHandlerTestSuite) TestCancelOrder() {
	tests := []struct {
		name      string
		setup     func(ctx *gin.Context)
		expected  int
		checkBody func(writer *httptest.ResponseRecorder)
	}{
		{
			name: "Success",
			setup: func(ctx *gin.Context) {
				ctx.Set("userId", "123456")
				ctx.AddParam("id", "orderId1")
				suite.mockService.On("CancelOrder", mock.Anything, "orderId1", "123456").
					Return(
						&model.Order{
							ID:         "orderId1",
							UserID:     "123456",
							TotalPrice: 5,
							Status:     model.OrderStatusNew,
						},
						nil,
					).Times(1)
			},
			expected: http.StatusOK,
			checkBody: func(writer *httptest.ResponseRecorder) {
				var res response.Response
				var orderRes domain.Order
				_ = json.Unmarshal(writer.Body.Bytes(), &res)
				_ = utils.Copy(&orderRes, &res.Result)
				suite.Equal("orderId1", orderRes.ID)
				suite.Equal(float64(5), orderRes.TotalPrice)
				suite.Equal(string(model.OrderStatusNew), orderRes.Status)
			},
		},
		{
			name: "Unauthorized",
			setup: func(ctx *gin.Context) {
				ctx.AddParam("id", "orderId1")
			},
			expected: http.StatusUnauthorized,
			checkBody: func(writer *httptest.ResponseRecorder) {
				var res response.Response
				_ = json.Unmarshal(writer.Body.Bytes(), &res)
				suite.NotNil(res.Error)
			},
		},
		{
			name: "MissID",
			setup: func(ctx *gin.Context) {
				ctx.Set("userId", "123456")
			},
			expected: http.StatusBadRequest,
			checkBody: func(writer *httptest.ResponseRecorder) {
				var res response.Response
				_ = json.Unmarshal(writer.Body.Bytes(), &res)
				suite.NotNil(res.Error)
			},
		},
		{
			name: "Fail",
			setup: func(ctx *gin.Context) {
				ctx.Set("userId", "123456")
				ctx.AddParam("id", "orderId1")
				suite.mockService.On("CancelOrder", mock.Anything, "orderId1", "123456").
					Return(nil, errors.New("error")).Times(1)
			},
			expected: http.StatusInternalServerError,
			checkBody: func(writer *httptest.ResponseRecorder) {
				var res response.Response
				_ = json.Unmarshal(writer.Body.Bytes(), &res)
				suite.NotNil(res.Error)
			},
		},
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			ctx, writer := suite.prepareContext(nil)
			tc.setup(ctx)
			suite.handler.CancelOrder(ctx)
			suite.Equal(tc.expected, writer.Code)
			if tc.checkBody != nil {
				tc.checkBody(writer)
			}
		})
	}
}
