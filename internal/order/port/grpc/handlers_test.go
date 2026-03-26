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

func (suite *OrderHandlerTestSuite) TestPlaceOrder() {
	tests := []struct {
		name      string
		setup     func()
		req       *pb.PlaceOrderReq
		ctx       context.Context
		expectNil bool
		expectErr bool
		validate  func(res *pb.PlaceOrderRes)
	}{
		{
			name: "Success",
			setup: func() {
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
			},
			req: &pb.PlaceOrderReq{
				Lines: []*pb.PlaceOrderLineReq{
					{ProductId: "productId1", Quantity: 2},
				},
			},
			ctx:       context.WithValue(context.Background(), "userId", "userID"),
			expectNil: false,
			expectErr: false,
			validate: func(res *pb.PlaceOrderRes) {
				suite.Equal("orderId1", res.Order.Id)
				suite.Equal("userID", res.Order.UserId)
				suite.Equal(1, len(res.Order.Lines))
			},
		},
		{
			name:  "Unauthorized",
			setup: func() {},
			req: &pb.PlaceOrderReq{
				Lines: []*pb.PlaceOrderLineReq{
					{ProductId: "productId1", Quantity: 2},
				},
			},
			ctx:       context.Background(),
			expectNil: true,
			expectErr: true,
		},
		{
			name: "Fail",
			setup: func() {
				suite.mockService.On("PlaceOrder", mock.Anything, mock.Anything).
					Return(nil, errors.New("error")).Times(1)
			},
			req: &pb.PlaceOrderReq{
				Lines: []*pb.PlaceOrderLineReq{
					{ProductId: "productId1", Quantity: 2},
				},
			},
			ctx:       context.WithValue(context.Background(), "userId", "userID"),
			expectNil: true,
			expectErr: true,
		},
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.setup()

			res, err := suite.handler.PlaceOrder(tc.ctx, tc.req)

			if tc.expectErr {
				suite.NotNil(err)
			} else {
				suite.Nil(err)
			}
			if tc.expectNil {
				suite.Nil(res)
			} else {
				suite.NotNil(res)
				if tc.validate != nil {
					tc.validate(res)
				}
			}
		})
	}
}

// GetOrderByID
// =================================================================================================

func (suite *OrderHandlerTestSuite) TestGetOrderByID() {
	tests := []struct {
		name      string
		setup     func()
		req       *pb.GetOrderByIDReq
		ctx       context.Context
		expectNil bool
		expectErr bool
		validate  func(res *pb.GetOrderByIDRes)
	}{
		{
			name: "Success",
			setup: func() {
				suite.mockService.On("GetOrderByID", mock.Anything, "orderId1").
					Return(&model.Order{
						ID:         "orderId1",
						Code:       "SO-001",
						UserID:     "userID",
						TotalPrice: 20.0,
						Status:     model.OrderStatusNew,
					}, nil).Times(1)
			},
			req:       &pb.GetOrderByIDReq{Id: "orderId1"},
			ctx:       context.WithValue(context.Background(), "userId", "userID"),
			expectNil: false,
			expectErr: false,
			validate: func(res *pb.GetOrderByIDRes) {
				suite.Equal("orderId1", res.Order.Id)
				suite.Equal("SO-001", res.Order.Code)
			},
		},
		{
			name:      "Unauthorized",
			setup:     func() {},
			req:       &pb.GetOrderByIDReq{Id: "orderId1"},
			ctx:       context.Background(),
			expectNil: true,
			expectErr: true,
		},
		{
			name:      "MissID",
			setup:     func() {},
			req:       &pb.GetOrderByIDReq{},
			ctx:       context.WithValue(context.Background(), "userId", "userID"),
			expectNil: true,
			expectErr: true,
		},
		{
			name: "Fail",
			setup: func() {
				suite.mockService.On("GetOrderByID", mock.Anything, "orderId1").
					Return(nil, errors.New("error")).Times(1)
			},
			req:       &pb.GetOrderByIDReq{Id: "orderId1"},
			ctx:       context.WithValue(context.Background(), "userId", "userID"),
			expectNil: true,
			expectErr: true,
		},
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.setup()

			res, err := suite.handler.GetOrderByID(tc.ctx, tc.req)

			if tc.expectErr {
				suite.NotNil(err)
			} else {
				suite.Nil(err)
			}
			if tc.expectNil {
				suite.Nil(res)
			} else {
				suite.NotNil(res)
				if tc.validate != nil {
					tc.validate(res)
				}
			}
		})
	}
}

// GetMyOrders
// =================================================================================================

func (suite *OrderHandlerTestSuite) TestGetMyOrders() {
	tests := []struct {
		name      string
		setup     func()
		req       *pb.GetMyOrdersReq
		ctx       context.Context
		expectNil bool
		expectErr bool
		validate  func(res *pb.GetMyOrdersRes)
	}{
		{
			name: "Success",
			setup: func() {
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
			},
			req: &pb.GetMyOrdersReq{
				Status: "new",
				Page:   1,
				Limit:  10,
			},
			ctx:       context.WithValue(context.Background(), "userId", "userID"),
			expectNil: false,
			expectErr: false,
			validate: func(res *pb.GetMyOrdersRes) {
				suite.Equal(1, len(res.Orders))
				suite.Equal("orderId1", res.Orders[0].Id)
				suite.Equal(int64(1), res.Total)
				suite.Equal(int64(1), res.CurrentPage)
				suite.Equal(int64(10), res.Limit)
			},
		},
		{
			name:      "Unauthorized",
			setup:     func() {},
			req:       &pb.GetMyOrdersReq{},
			ctx:       context.Background(),
			expectNil: true,
			expectErr: true,
		},
		{
			name: "Fail",
			setup: func() {
				suite.mockService.On("GetMyOrders", mock.Anything, mock.Anything).
					Return(nil, nil, errors.New("error")).Times(1)
			},
			req:       &pb.GetMyOrdersReq{},
			ctx:       context.WithValue(context.Background(), "userId", "userID"),
			expectNil: true,
			expectErr: true,
		},
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.setup()

			res, err := suite.handler.GetMyOrders(tc.ctx, tc.req)

			if tc.expectErr {
				suite.NotNil(err)
			} else {
				suite.Nil(err)
			}
			if tc.expectNil {
				suite.Nil(res)
			} else {
				suite.NotNil(res)
				if tc.validate != nil {
					tc.validate(res)
				}
			}
		})
	}
}

// CancelOrder
// =================================================================================================

func (suite *OrderHandlerTestSuite) TestCancelOrder() {
	tests := []struct {
		name      string
		setup     func()
		req       *pb.CancelOrderReq
		ctx       context.Context
		expectNil bool
		expectErr bool
		validate  func(res *pb.CancelOrderRes)
	}{
		{
			name: "Success",
			setup: func() {
				suite.mockService.On("CancelOrder", mock.Anything, "orderId1", "userID").
					Return(&model.Order{
						ID:         "orderId1",
						Code:       "SO-001",
						UserID:     "userID",
						TotalPrice: 20.0,
						Status:     model.OrderStatusCancelled,
					}, nil).Times(1)
			},
			req:       &pb.CancelOrderReq{Id: "orderId1"},
			ctx:       context.WithValue(context.Background(), "userId", "userID"),
			expectNil: false,
			expectErr: false,
			validate: func(res *pb.CancelOrderRes) {
				suite.Equal("orderId1", res.Order.Id)
				suite.Equal(string(model.OrderStatusCancelled), res.Order.Status)
			},
		},
		{
			name:      "Unauthorized",
			setup:     func() {},
			req:       &pb.CancelOrderReq{Id: "orderId1"},
			ctx:       context.Background(),
			expectNil: true,
			expectErr: true,
		},
		{
			name:      "MissID",
			setup:     func() {},
			req:       &pb.CancelOrderReq{},
			ctx:       context.WithValue(context.Background(), "userId", "userID"),
			expectNil: true,
			expectErr: true,
		},
		{
			name: "Fail",
			setup: func() {
				suite.mockService.On("CancelOrder", mock.Anything, "orderId1", "userID").
					Return(nil, errors.New("error")).Times(1)
			},
			req:       &pb.CancelOrderReq{Id: "orderId1"},
			ctx:       context.WithValue(context.Background(), "userId", "userID"),
			expectNil: true,
			expectErr: true,
		},
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.setup()

			res, err := suite.handler.CancelOrder(tc.ctx, tc.req)

			if tc.expectErr {
				suite.NotNil(err)
			} else {
				suite.Nil(err)
			}
			if tc.expectNil {
				suite.Nil(res)
			} else {
				suite.NotNil(res)
				if tc.validate != nil {
					tc.validate(res)
				}
			}
		})
	}
}
