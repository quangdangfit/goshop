package grpc

import (
	"context"
	"errors"
	"testing"

	"github.com/quangdangfit/gocommon/logger"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"goshop/internal/order/model"
	"goshop/internal/order/service/mocks"
	"goshop/pkg/config"
	"goshop/pkg/paging"
	pb "goshop/proto/gen/go/order"
)

type OrderHandlerTestSuite struct {
	suite.Suite
	mockService *mocks.OrderService
	handler     *OrderHandler
}

func (suite *OrderHandlerTestSuite) SetupTest() {
	logger.Initialize(config.ProductionEnv)

	suite.mockService = mocks.NewOrderService(suite.T())
	suite.handler = NewOrderHandler(suite.mockService)
}

func TestOrderHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(OrderHandlerTestSuite))
}

// PlaceOrder
// =================================================================================================

func (suite *OrderHandlerTestSuite) TestPlaceOrderSuccess() {
	req := &pb.PlaceOrderReq{
		Lines: []*pb.PlaceOrderLineReq{
			{ProductId: "productId1", Quantity: 2},
		},
	}

	suite.mockService.On("PlaceOrder", mock.Anything, mock.Anything).
		Return(&model.Order{
			ID:         "orderId1",
			Code:       "SO-001",
			UserID:     "userID",
			TotalPrice: 20.0,
			Status:     model.OrderStatusNew,
			Lines: []*model.OrderLine{
				{ProductID: "productId1", Quantity: 2, Price: 20.0},
			},
		}, nil).Times(1)

	ctx := context.WithValue(context.Background(), "userId", "userID")
	res, err := suite.handler.PlaceOrder(ctx, req)

	suite.Nil(err)
	suite.NotNil(res)
	suite.Equal("orderId1", res.Order.Id)
	suite.Equal("userID", res.Order.UserId)
	suite.Equal(1, len(res.Order.Lines))
}

func (suite *OrderHandlerTestSuite) TestPlaceOrderUnauthorized() {
	req := &pb.PlaceOrderReq{
		Lines: []*pb.PlaceOrderLineReq{
			{ProductId: "productId1", Quantity: 2},
		},
	}

	res, err := suite.handler.PlaceOrder(context.Background(), req)

	suite.Nil(res)
	suite.NotNil(err)
}

func (suite *OrderHandlerTestSuite) TestPlaceOrderFail() {
	req := &pb.PlaceOrderReq{
		Lines: []*pb.PlaceOrderLineReq{
			{ProductId: "productId1", Quantity: 2},
		},
	}

	suite.mockService.On("PlaceOrder", mock.Anything, mock.Anything).
		Return(nil, errors.New("error")).Times(1)

	ctx := context.WithValue(context.Background(), "userId", "userID")
	res, err := suite.handler.PlaceOrder(ctx, req)

	suite.Nil(res)
	suite.NotNil(err)
}

// GetOrderByID
// =================================================================================================

func (suite *OrderHandlerTestSuite) TestGetOrderByIDSuccess() {
	req := &pb.GetOrderByIDReq{Id: "orderId1"}

	suite.mockService.On("GetOrderByID", mock.Anything, "orderId1").
		Return(&model.Order{
			ID:         "orderId1",
			Code:       "SO-001",
			UserID:     "userID",
			TotalPrice: 20.0,
			Status:     model.OrderStatusNew,
		}, nil).Times(1)

	ctx := context.WithValue(context.Background(), "userId", "userID")
	res, err := suite.handler.GetOrderByID(ctx, req)

	suite.Nil(err)
	suite.NotNil(res)
	suite.Equal("orderId1", res.Order.Id)
	suite.Equal("SO-001", res.Order.Code)
}

func (suite *OrderHandlerTestSuite) TestGetOrderByIDUnauthorized() {
	req := &pb.GetOrderByIDReq{Id: "orderId1"}

	res, err := suite.handler.GetOrderByID(context.Background(), req)

	suite.Nil(res)
	suite.NotNil(err)
}

func (suite *OrderHandlerTestSuite) TestGetOrderByIDMissID() {
	req := &pb.GetOrderByIDReq{}

	ctx := context.WithValue(context.Background(), "userId", "userID")
	res, err := suite.handler.GetOrderByID(ctx, req)

	suite.Nil(res)
	suite.NotNil(err)
}

func (suite *OrderHandlerTestSuite) TestGetOrderByIDFail() {
	req := &pb.GetOrderByIDReq{Id: "orderId1"}

	suite.mockService.On("GetOrderByID", mock.Anything, "orderId1").
		Return(nil, errors.New("error")).Times(1)

	ctx := context.WithValue(context.Background(), "userId", "userID")
	res, err := suite.handler.GetOrderByID(ctx, req)

	suite.Nil(res)
	suite.NotNil(err)
}

// GetMyOrders
// =================================================================================================

func (suite *OrderHandlerTestSuite) TestGetMyOrdersSuccess() {
	req := &pb.GetMyOrdersReq{
		Status: "new",
		Page:   1,
		Limit:  10,
	}

	suite.mockService.On("GetMyOrders", mock.Anything, mock.Anything).
		Return(
			[]*model.Order{
				{
					ID:         "orderId1",
					Code:       "SO-001",
					UserID:     "userID",
					TotalPrice: 20.0,
					Status:     model.OrderStatusNew,
				},
			},
			&paging.Pagination{
				Total:       1,
				CurrentPage: 1,
				Limit:       10,
			},
			nil,
		).Times(1)

	ctx := context.WithValue(context.Background(), "userId", "userID")
	res, err := suite.handler.GetMyOrders(ctx, req)

	suite.Nil(err)
	suite.NotNil(res)
	suite.Equal(1, len(res.Orders))
	suite.Equal("orderId1", res.Orders[0].Id)
	suite.Equal(int64(1), res.Total)
	suite.Equal(int64(1), res.CurrentPage)
	suite.Equal(int64(10), res.Limit)
}

func (suite *OrderHandlerTestSuite) TestGetMyOrdersUnauthorized() {
	req := &pb.GetMyOrdersReq{}

	res, err := suite.handler.GetMyOrders(context.Background(), req)

	suite.Nil(res)
	suite.NotNil(err)
}

func (suite *OrderHandlerTestSuite) TestGetMyOrdersFail() {
	req := &pb.GetMyOrdersReq{}

	suite.mockService.On("GetMyOrders", mock.Anything, mock.Anything).
		Return(nil, nil, errors.New("error")).Times(1)

	ctx := context.WithValue(context.Background(), "userId", "userID")
	res, err := suite.handler.GetMyOrders(ctx, req)

	suite.Nil(res)
	suite.NotNil(err)
}

// CancelOrder
// =================================================================================================

func (suite *OrderHandlerTestSuite) TestCancelOrderSuccess() {
	req := &pb.CancelOrderReq{Id: "orderId1"}

	suite.mockService.On("CancelOrder", mock.Anything, "orderId1", "userID").
		Return(&model.Order{
			ID:         "orderId1",
			Code:       "SO-001",
			UserID:     "userID",
			TotalPrice: 20.0,
			Status:     model.OrderStatusCancelled,
		}, nil).Times(1)

	ctx := context.WithValue(context.Background(), "userId", "userID")
	res, err := suite.handler.CancelOrder(ctx, req)

	suite.Nil(err)
	suite.NotNil(res)
	suite.Equal("orderId1", res.Order.Id)
	suite.Equal(string(model.OrderStatusCancelled), res.Order.Status)
}

func (suite *OrderHandlerTestSuite) TestCancelOrderUnauthorized() {
	req := &pb.CancelOrderReq{Id: "orderId1"}

	res, err := suite.handler.CancelOrder(context.Background(), req)

	suite.Nil(res)
	suite.NotNil(err)
}

func (suite *OrderHandlerTestSuite) TestCancelOrderMissID() {
	req := &pb.CancelOrderReq{}

	ctx := context.WithValue(context.Background(), "userId", "userID")
	res, err := suite.handler.CancelOrder(ctx, req)

	suite.Nil(res)
	suite.NotNil(err)
}

func (suite *OrderHandlerTestSuite) TestCancelOrderFail() {
	req := &pb.CancelOrderReq{Id: "orderId1"}

	suite.mockService.On("CancelOrder", mock.Anything, "orderId1", "userID").
		Return(nil, errors.New("error")).Times(1)

	ctx := context.WithValue(context.Background(), "userId", "userID")
	res, err := suite.handler.CancelOrder(ctx, req)

	suite.Nil(res)
	suite.NotNil(err)
}
