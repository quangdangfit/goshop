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
	svcMocks "goshop/internal/order/service/mocks"
	"goshop/pkg/config"
	"goshop/pkg/response"
	"goshop/pkg/utils"
)

type CouponHandlerTestSuite struct {
	suite.Suite
	mockCouponService *svcMocks.CouponService
	mockOrderService  *svcMocks.OrderService
	couponHandler     *CouponHandler
	orderHandler      *OrderHandler
}

func (suite *CouponHandlerTestSuite) SetupTest() {
	logger.Initialize(config.ProductionEnv)
	suite.mockCouponService = svcMocks.NewCouponService(suite.T())
	suite.mockOrderService = svcMocks.NewOrderService(suite.T())
	suite.couponHandler = NewCouponHandler(suite.mockCouponService)
	suite.orderHandler = NewOrderHandler(suite.mockOrderService)
}

func TestCouponHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(CouponHandlerTestSuite))
}

func (suite *CouponHandlerTestSuite) prepareContext(method, path string, body any) (*gin.Context, *httptest.ResponseRecorder) {
	requestBody, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(method, path, bytes.NewBuffer(requestBody))
	c, _ := gin.CreateTestContext(w)
	c.Request = r
	return c, w
}

func (suite *CouponHandlerTestSuite) prepareContextWithQuery(path string, query url.Values) (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("PUT", path+"?"+query.Encode(), nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = r
	return c, w
}

func (suite *CouponHandlerTestSuite) TestCreateCoupon() {
	tests := []struct {
		name      string
		body      any
		setup     func()
		expected  int
		checkBody func(writer *httptest.ResponseRecorder)
	}{
		{
			name: "Success",
			body: &domain.CreateCouponReq{
				Code:          "SAVE10",
				DiscountType:  "fixed",
				DiscountValue: 10,
			},
			setup: func() {
				suite.mockCouponService.On("Create", mock.Anything, &domain.CreateCouponReq{
					Code:          "SAVE10",
					DiscountType:  "fixed",
					DiscountValue: 10,
				}).Return(&model.Coupon{
					ID:            "c1",
					Code:          "SAVE10",
					DiscountType:  model.DiscountTypeFixed,
					DiscountValue: 10,
				}, nil).Times(1)
			},
			expected: http.StatusOK,
			checkBody: func(writer *httptest.ResponseRecorder) {
				var res response.Response
				var coupon domain.Coupon
				_ = json.Unmarshal(writer.Body.Bytes(), &res)
				utils.Copy(&coupon, &res.Result)
				suite.Equal("SAVE10", coupon.Code)
			},
		},
		{
			name:     "InvalidBody",
			body:     map[string]any{"code": 123},
			setup:    func() {},
			expected: http.StatusBadRequest,
		},
		{
			name: "Fail",
			body: &domain.CreateCouponReq{
				Code:          "SAVE10",
				DiscountType:  "fixed",
				DiscountValue: 10,
			},
			setup: func() {
				suite.mockCouponService.On("Create", mock.Anything, &domain.CreateCouponReq{
					Code:          "SAVE10",
					DiscountType:  "fixed",
					DiscountValue: 10,
				}).Return(nil, errors.New("duplicate code")).Times(1)
			},
			expected: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			ctx, writer := suite.prepareContext("POST", "/api/v1/coupons", tc.body)
			tc.setup()
			suite.couponHandler.CreateCoupon(ctx)
			suite.Equal(tc.expected, writer.Code)
			if tc.checkBody != nil {
				tc.checkBody(writer)
			}
		})
	}
}

func (suite *CouponHandlerTestSuite) TestGetCouponByCode() {
	tests := []struct {
		name      string
		params    gin.Params
		setup     func()
		expected  int
		checkBody func(writer *httptest.ResponseRecorder)
	}{
		{
			name:   "Success",
			params: gin.Params{{Key: "code", Value: "SAVE10"}},
			setup: func() {
				suite.mockCouponService.On("GetByCode", mock.Anything, "SAVE10").
					Return(&model.Coupon{
						ID:            "c1",
						Code:          "SAVE10",
						DiscountType:  model.DiscountTypeFixed,
						DiscountValue: 10,
					}, nil).Times(1)
			},
			expected: http.StatusOK,
			checkBody: func(writer *httptest.ResponseRecorder) {
				var res response.Response
				var coupon domain.Coupon
				_ = json.Unmarshal(writer.Body.Bytes(), &res)
				utils.Copy(&coupon, &res.Result)
				suite.Equal("SAVE10", coupon.Code)
			},
		},
		{
			name:   "NotFound",
			params: gin.Params{{Key: "code", Value: "INVALID"}},
			setup: func() {
				suite.mockCouponService.On("GetByCode", mock.Anything, "INVALID").
					Return(nil, errors.New("not found")).Times(1)
			},
			expected: http.StatusNotFound,
		},
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			ctx, writer := suite.prepareContext("GET", "/api/v1/coupons/"+tc.params[0].Value, nil)
			ctx.Params = tc.params
			tc.setup()
			suite.couponHandler.GetCouponByCode(ctx)
			suite.Equal(tc.expected, writer.Code)
			if tc.checkBody != nil {
				tc.checkBody(writer)
			}
		})
	}
}

func (suite *CouponHandlerTestSuite) TestUpdateOrderStatus() {
	tests := []struct {
		name     string
		params   gin.Params
		query    url.Values
		setup    func()
		expected int
	}{
		{
			name:   "Success",
			params: gin.Params{{Key: "id", Value: "o1"}},
			query:  url.Values{"status": []string{"done"}},
			setup: func() {
				suite.mockOrderService.On("UpdateOrderStatus", mock.Anything, "o1", model.OrderStatus("done")).
					Return(&model.Order{ID: "o1", Status: model.OrderStatusDone}, nil).Times(1)
			},
			expected: http.StatusOK,
		},
		{
			name:     "MissingOrderID",
			params:   nil,
			query:    url.Values{"status": []string{"done"}},
			setup:    func() {},
			expected: http.StatusBadRequest,
		},
		{
			name:     "MissingStatus",
			params:   gin.Params{{Key: "id", Value: "o1"}},
			query:    url.Values{},
			setup:    func() {},
			expected: http.StatusBadRequest,
		},
		{
			name:   "Fail",
			params: gin.Params{{Key: "id", Value: "o1"}},
			query:  url.Values{"status": []string{"done"}},
			setup: func() {
				suite.mockOrderService.On("UpdateOrderStatus", mock.Anything, "o1", model.OrderStatus("done")).
					Return(nil, errors.New("not found")).Times(1)
			},
			expected: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			ctx, writer := suite.prepareContextWithQuery("/api/v1/orders/o1/status", tc.query)
			if tc.params != nil {
				ctx.Params = tc.params
			}
			tc.setup()
			suite.orderHandler.UpdateOrderStatus(ctx)
			suite.Equal(tc.expected, writer.Code)
		})
	}
}
