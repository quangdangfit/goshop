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

	"goshop/internal/order/dto"
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

// CreateCoupon
// =================================================================================================

func (suite *CouponHandlerTestSuite) TestCreateCouponSuccess() {
	req := &dto.CreateCouponReq{
		Code:          "SAVE10",
		DiscountType:  "fixed",
		DiscountValue: 10,
	}
	ctx, writer := suite.prepareContext("POST", "/api/v1/coupons", req)

	suite.mockCouponService.On("Create", mock.Anything, req).
		Return(&model.Coupon{
			ID:            "c1",
			Code:          "SAVE10",
			DiscountType:  model.DiscountTypeFixed,
			DiscountValue: 10,
		}, nil).Times(1)

	suite.couponHandler.CreateCoupon(ctx)

	var res response.Response
	var coupon dto.Coupon
	_ = json.Unmarshal(writer.Body.Bytes(), &res)
	utils.Copy(&coupon, &res.Result)
	suite.Equal(http.StatusOK, writer.Code)
	suite.Equal("SAVE10", coupon.Code)
}

func (suite *CouponHandlerTestSuite) TestCreateCouponInvalidBody() {
	req := map[string]any{"code": 123}
	ctx, writer := suite.prepareContext("POST", "/api/v1/coupons", req)

	suite.couponHandler.CreateCoupon(ctx)

	suite.Equal(http.StatusBadRequest, writer.Code)
}

func (suite *CouponHandlerTestSuite) TestCreateCouponFail() {
	req := &dto.CreateCouponReq{
		Code:          "SAVE10",
		DiscountType:  "fixed",
		DiscountValue: 10,
	}
	ctx, writer := suite.prepareContext("POST", "/api/v1/coupons", req)

	suite.mockCouponService.On("Create", mock.Anything, req).
		Return(nil, errors.New("duplicate code")).Times(1)

	suite.couponHandler.CreateCoupon(ctx)

	suite.Equal(http.StatusInternalServerError, writer.Code)
}

// GetCouponByCode
// =================================================================================================

func (suite *CouponHandlerTestSuite) TestGetCouponByCodeSuccess() {
	ctx, writer := suite.prepareContext("GET", "/api/v1/coupons/SAVE10", nil)
	ctx.Params = gin.Params{{Key: "code", Value: "SAVE10"}}

	suite.mockCouponService.On("GetByCode", mock.Anything, "SAVE10").
		Return(&model.Coupon{
			ID:            "c1",
			Code:          "SAVE10",
			DiscountType:  model.DiscountTypeFixed,
			DiscountValue: 10,
		}, nil).Times(1)

	suite.couponHandler.GetCouponByCode(ctx)

	var res response.Response
	var coupon dto.Coupon
	_ = json.Unmarshal(writer.Body.Bytes(), &res)
	utils.Copy(&coupon, &res.Result)
	suite.Equal(http.StatusOK, writer.Code)
	suite.Equal("SAVE10", coupon.Code)
}

func (suite *CouponHandlerTestSuite) TestGetCouponByCodeNotFound() {
	ctx, writer := suite.prepareContext("GET", "/api/v1/coupons/INVALID", nil)
	ctx.Params = gin.Params{{Key: "code", Value: "INVALID"}}

	suite.mockCouponService.On("GetByCode", mock.Anything, "INVALID").
		Return(nil, errors.New("not found")).Times(1)

	suite.couponHandler.GetCouponByCode(ctx)

	suite.Equal(http.StatusNotFound, writer.Code)
}

// UpdateOrderStatus
// =================================================================================================

func (suite *CouponHandlerTestSuite) TestUpdateOrderStatusSuccess() {
	q := url.Values{}
	q.Set("status", "done")
	ctx, writer := suite.prepareContextWithQuery("/api/v1/orders/o1/status", q)
	ctx.Params = gin.Params{{Key: "id", Value: "o1"}}

	suite.mockOrderService.On("UpdateOrderStatus", mock.Anything, "o1", model.OrderStatus("done")).
		Return(&model.Order{ID: "o1", Status: model.OrderStatusDone}, nil).Times(1)

	suite.orderHandler.UpdateOrderStatus(ctx)

	suite.Equal(http.StatusOK, writer.Code)
}

func (suite *CouponHandlerTestSuite) TestUpdateOrderStatusMissingOrderID() {
	q := url.Values{}
	q.Set("status", "done")
	ctx, writer := suite.prepareContextWithQuery("/api/v1/orders//status", q)
	// id param is empty

	suite.orderHandler.UpdateOrderStatus(ctx)

	suite.Equal(http.StatusBadRequest, writer.Code)
}

func (suite *CouponHandlerTestSuite) TestUpdateOrderStatusMissingStatus() {
	ctx, writer := suite.prepareContextWithQuery("/api/v1/orders/o1/status", url.Values{})
	ctx.Params = gin.Params{{Key: "id", Value: "o1"}}

	suite.orderHandler.UpdateOrderStatus(ctx)

	suite.Equal(http.StatusBadRequest, writer.Code)
}

func (suite *CouponHandlerTestSuite) TestUpdateOrderStatusFail() {
	q := url.Values{}
	q.Set("status", "done")
	ctx, writer := suite.prepareContextWithQuery("/api/v1/orders/o1/status", q)
	ctx.Params = gin.Params{{Key: "id", Value: "o1"}}

	suite.mockOrderService.On("UpdateOrderStatus", mock.Anything, "o1", model.OrderStatus("done")).
		Return(nil, errors.New("not found")).Times(1)

	suite.orderHandler.UpdateOrderStatus(ctx)

	suite.Equal(http.StatusInternalServerError, writer.Code)
}
