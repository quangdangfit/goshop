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
	"goshop/internal/order/service/mocks"
	productMocks "goshop/internal/product/service/mocks"
	"goshop/pkg/config"
	"goshop/pkg/paging"
	"goshop/pkg/response"
	"goshop/pkg/utils"
)

type OrderHandlerTestSuite struct {
	suite.Suite
	mockService        *mocks.IOrderService
	mockProductService *productMocks.IProductService
	handler            *OrderHandler
}

func (suite *OrderHandlerTestSuite) SetupTest() {
	logger.Initialize(config.ProductionEnv)

	suite.mockService = mocks.NewIOrderService(suite.T())
	suite.mockProductService = productMocks.NewIProductService(suite.T())
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

// PlaceOrder
// =================================================================================================

func (suite *OrderHandlerTestSuite) TestOrderAPI_PlaceOrderSuccess() {
	req := &dto.PlaceOrderReq{
		Lines: []dto.PlaceOrderLineReq{
			{
				ProductID: "productId1",
				Quantity:  2,
			},
			{
				ProductID: "productId2",
				Quantity:  3,
			},
		},
	}

	ctx, writer := suite.prepareContext(req)
	ctx.Set("userId", "123456")
	req.UserID = "123456"

	suite.mockService.On("PlaceOrder", mock.Anything, req).
		Return(
			&model.Order{
				ID:         "orderId1",
				Code:       "orderCode1",
				TotalPrice: 8,
				Status:     model.OrderStatusNew,
				Lines: []*model.OrderLine{
					{
						ProductID: "productId1",
						Quantity:  2,
					},
					{
						ProductID: "productId2",
						Quantity:  3,
					},
				},
			},
			nil,
		).Times(1)

	suite.handler.PlaceOrder(ctx)

	var res response.Response
	var orderRes dto.Order

	_ = json.Unmarshal(writer.Body.Bytes(), &res)
	utils.Copy(&orderRes, &res.Result)
	suite.Equal(http.StatusOK, writer.Code)
	suite.Equal(float64(8), orderRes.TotalPrice)
	suite.Equal(string(model.OrderStatusNew), orderRes.Status)
	suite.Equal(2, len(orderRes.Lines))
}

func (suite *OrderHandlerTestSuite) TestOrderAPI_PlaceOrderInvalidProductIdType() {
	req := map[string]interface{}{
		"lines": []map[string]interface{}{
			{
				"product_id": 1,
				"quantity":   2,
			},
		},
	}

	ctx, writer := suite.prepareContext(req)

	suite.handler.PlaceOrder(ctx)

	var res map[string]map[string]string
	_ = json.Unmarshal(writer.Body.Bytes(), &res)
	suite.Equal(http.StatusBadRequest, writer.Code)
	suite.Equal("Invalid parameters", res["error"]["message"])
}

func (suite *OrderHandlerTestSuite) TestOrderAPI_PlaceOrderInvalidQuantityType() {
	req := map[string]interface{}{
		"lines": []map[string]interface{}{
			{
				"product_id": "productId1",
				"quantity":   "1",
			},
		},
	}

	ctx, writer := suite.prepareContext(req)

	suite.handler.PlaceOrder(ctx)

	var res map[string]map[string]string
	_ = json.Unmarshal(writer.Body.Bytes(), &res)
	suite.Equal(http.StatusBadRequest, writer.Code)
	suite.Equal("Invalid parameters", res["error"]["message"])
}

func (suite *OrderHandlerTestSuite) TestOrderAPI_PlaceOrderUnauthorized() {
	req := &dto.PlaceOrderReq{
		Lines: []dto.PlaceOrderLineReq{
			{
				ProductID: "productId1",
				Quantity:  2,
			},
			{
				ProductID: "productId2",
				Quantity:  3,
			},
		},
	}

	ctx, writer := suite.prepareContext(req)

	suite.handler.PlaceOrder(ctx)

	var res map[string]map[string]string
	_ = json.Unmarshal(writer.Body.Bytes(), &res)
	suite.Equal(http.StatusUnauthorized, writer.Code)
	suite.Equal("Unauthorized", res["error"]["message"])
}

func (suite *OrderHandlerTestSuite) TestOrderAPI_PlaceOrderFail() {
	req := &dto.PlaceOrderReq{
		Lines: []dto.PlaceOrderLineReq{
			{
				ProductID: "productId1",
				Quantity:  2,
			},
			{
				ProductID: "productId2",
				Quantity:  3,
			},
		},
	}

	ctx, writer := suite.prepareContext(req)
	ctx.Set("userId", "123456")
	req.UserID = "123456"

	suite.mockService.On("PlaceOrder", mock.Anything, req).
		Return(nil, errors.New("error")).Times(1)

	suite.handler.PlaceOrder(ctx)

	var res map[string]map[string]string
	_ = json.Unmarshal(writer.Body.Bytes(), &res)
	suite.Equal(http.StatusInternalServerError, writer.Code)
	suite.Equal("Something went wrong", res["error"]["message"])
}

// Get Order Detail
// =================================================================================================

func (suite *OrderHandlerTestSuite) TestOrderAPI_GetOrderByIDSuccess() {
	ctx, writer := suite.prepareContext(nil)
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

	suite.handler.GetOrderByID(ctx)

	var res response.Response
	var orderRes dto.Order

	_ = json.Unmarshal(writer.Body.Bytes(), &res)
	utils.Copy(&orderRes, &res.Result)
	suite.Equal(http.StatusOK, writer.Code)
	suite.Equal(float64(5), orderRes.TotalPrice)
	suite.Equal(string(model.OrderStatusNew), orderRes.Status)
	suite.Equal(0, len(orderRes.Lines))
}

func (suite *OrderHandlerTestSuite) TestOrderAPI_GetOrderByIDMissID() {
	ctx, writer := suite.prepareContext(nil)
	ctx.Set("userId", "123456")

	suite.handler.GetOrderByID(ctx)

	var res response.Response
	_ = json.Unmarshal(writer.Body.Bytes(), &res)
	suite.Equal(http.StatusBadRequest, writer.Code)
	suite.NotNil(res.Error)
}

func (suite *OrderHandlerTestSuite) TestOrderAPI_GetOrderByIDUnauthorized() {
	ctx, writer := suite.prepareContext(nil)
	ctx.AddParam("id", "orderId1")

	suite.handler.GetOrderByID(ctx)

	var res response.Response
	_ = json.Unmarshal(writer.Body.Bytes(), &res)
	suite.Equal(http.StatusUnauthorized, writer.Code)
	suite.NotNil(res.Error)
}

func (suite *OrderHandlerTestSuite) TestOrderAPI_GetOrderByIDFail() {
	ctx, writer := suite.prepareContext(nil)
	ctx.Set("userId", "123456")
	ctx.AddParam("id", "orderId1")

	suite.mockService.On("GetOrderByID", mock.Anything, "orderId1").
		Return(nil, errors.New("error")).Times(1)

	suite.handler.GetOrderByID(ctx)

	var res response.Response
	_ = json.Unmarshal(writer.Body.Bytes(), &res)
	suite.Equal(http.StatusNotFound, writer.Code)
	suite.NotNil(res.Error)
}

// GetOrders
// =================================================================================================

func (suite *OrderHandlerTestSuite) TestOrderAPI_GetMyOrdersSuccess() {
	ctx, writer := suite.prepareContext(nil)
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

	suite.handler.GetOrders(ctx)

	var res response.Response
	var orderRes dto.ListOrderRes

	_ = json.Unmarshal(writer.Body.Bytes(), &res)
	utils.Copy(&orderRes, &res.Result)
	suite.Equal(http.StatusOK, writer.Code)
	suite.Equal(1, len(orderRes.Orders))
	suite.Equal("orderId1", orderRes.Orders[0].ID)
	suite.Equal(float64(5), orderRes.Orders[0].TotalPrice)
	suite.Equal(string(model.OrderStatusNew), orderRes.Orders[0].Status)
}

func (suite *OrderHandlerTestSuite) TestOrderAPI_GetMyOrdersUnauthorized() {
	ctx, writer := suite.prepareContext(nil)

	suite.handler.GetOrders(ctx)

	var res response.Response
	_ = json.Unmarshal(writer.Body.Bytes(), &res)
	suite.Equal(http.StatusUnauthorized, writer.Code)
	suite.NotNil(res.Error)
}

func (suite *OrderHandlerTestSuite) TestOrderAPI_GetMyOrdersInvalidFieldType() {
	ctx, writer := suite.prepareContext(nil)
	ctx.Request.URL, _ = url.Parse("?page=q")

	suite.handler.GetOrders(ctx)

	var res response.Response
	_ = json.Unmarshal(writer.Body.Bytes(), &res)
	suite.Equal(http.StatusBadRequest, writer.Code)
	suite.NotNil(res.Error)
}

func (suite *OrderHandlerTestSuite) TestOrderAPI_GetMyOrdersFail() {
	ctx, writer := suite.prepareContext(nil)
	ctx.Set("userId", "123456")

	suite.mockService.On("GetMyOrders", mock.Anything, mock.Anything).
		Return(nil, nil, errors.New("error")).Times(1)

	suite.handler.GetOrders(ctx)

	var res response.Response
	_ = json.Unmarshal(writer.Body.Bytes(), &res)
	suite.Equal(http.StatusInternalServerError, writer.Code)
	suite.NotNil(res.Error)
}

// CancelOrder
// =================================================================================================

func (suite *OrderHandlerTestSuite) TestOrderAPI_CancelOrderSuccess() {
	ctx, writer := suite.prepareContext(nil)
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

	suite.handler.CancelOrder(ctx)

	var res response.Response
	var orderRes dto.Order

	_ = json.Unmarshal(writer.Body.Bytes(), &res)
	utils.Copy(&orderRes, &res.Result)
	suite.Equal(http.StatusOK, writer.Code)
	suite.Equal("orderId1", orderRes.ID)
	suite.Equal(float64(5), orderRes.TotalPrice)
	suite.Equal(string(model.OrderStatusNew), orderRes.Status)
}

func (suite *OrderHandlerTestSuite) TestOrderAPI_CancelOrderUnauthorized() {
	ctx, writer := suite.prepareContext(nil)
	ctx.AddParam("id", "orderId1")

	suite.handler.CancelOrder(ctx)

	var res response.Response
	_ = json.Unmarshal(writer.Body.Bytes(), &res)
	suite.Equal(http.StatusUnauthorized, writer.Code)
	suite.NotNil(res.Error)
}

func (suite *OrderHandlerTestSuite) TestOrderAPI_CancelOrderMissID() {
	ctx, writer := suite.prepareContext(nil)
	ctx.Set("userId", "123456")

	suite.handler.CancelOrder(ctx)

	var res response.Response
	_ = json.Unmarshal(writer.Body.Bytes(), &res)
	suite.Equal(http.StatusBadRequest, writer.Code)
	suite.NotNil(res.Error)
}

func (suite *OrderHandlerTestSuite) TestOrderAPI_CancelOrderFail() {
	ctx, writer := suite.prepareContext(nil)
	ctx.Set("userId", "123456")
	ctx.AddParam("id", "orderId1")

	suite.mockService.On("CancelOrder", mock.Anything, "orderId1", "123456").
		Return(nil, errors.New("error")).Times(1)

	suite.handler.CancelOrder(ctx)

	var res response.Response
	_ = json.Unmarshal(writer.Body.Bytes(), &res)
	suite.Equal(http.StatusInternalServerError, writer.Code)
	suite.NotNil(res.Error)
}
