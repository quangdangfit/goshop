package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/quangdangfit/gocommon/logger"
	"github.com/quangdangfit/gocommon/validation"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	cartMocks "goshop/internal/cart/repository/mocks"
	"goshop/internal/order/domain"
	"goshop/internal/order/model"
	orderMocks "goshop/internal/order/repository/mocks"
	serviceMocks "goshop/internal/order/service/mocks"
	"goshop/pkg/config"
	notifMocks "goshop/pkg/notification/mocks"
	"goshop/pkg/paging"
)

type OrderServiceTestSuite struct {
	suite.Suite
	mockRepo        *orderMocks.OrderRepository
	mockProductRepo *orderMocks.ProductRepository
	mockUserRepo    *orderMocks.UserRepository
	mockCartRepo    *cartMocks.CartRepository
	mockCouponSvc   *serviceMocks.CouponService
	mockNotifier    *notifMocks.Notifier
	service         OrderService
}

func (suite *OrderServiceTestSuite) SetupTest() {
	logger.Initialize(config.ProductionEnv)

	validator := validation.New()
	suite.mockRepo = orderMocks.NewOrderRepository(suite.T())
	suite.mockProductRepo = orderMocks.NewProductRepository(suite.T())
	suite.mockUserRepo = orderMocks.NewUserRepository(suite.T())
	suite.mockCartRepo = cartMocks.NewCartRepository(suite.T())
	suite.mockCouponSvc = serviceMocks.NewCouponService(suite.T())
	suite.mockNotifier = notifMocks.NewNotifier(suite.T())
	suite.service = NewOrderService(
		validator,
		suite.mockRepo,
		suite.mockProductRepo,
		suite.mockUserRepo,
		suite.mockCartRepo,
		suite.mockCouponSvc,
		suite.mockNotifier,
	)
}

func TestOrderServiceTestSuite(t *testing.T) {
	suite.Run(t, new(OrderServiceTestSuite))
}

func (suite *OrderServiceTestSuite) TestGetOrderByID() {
	tests := []struct {
		name    string
		setup   func()
		wantErr bool
	}{
		{
			name: "Success",
			setup: func() {
				suite.mockRepo.On("GetOrderByID", mock.Anything, "orderID", true).
					Return(&model.Order{UserID: "userID", TotalPrice: 111.1, Status: model.OrderStatusNew}, nil).Times(1)
			},
		},
		{
			name: "Not found",
			setup: func() {
				suite.mockRepo.On("GetOrderByID", mock.Anything, "orderID", true).
					Return(nil, errors.New("error")).Times(1)
			},
			wantErr: true,
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.setup()
			order, err := suite.service.GetOrderByID(context.Background(), "orderID")
			if tc.wantErr {
				suite.Nil(order)
				suite.NotNil(err)
			} else {
				suite.NotNil(order)
				suite.Equal("userID", order.UserID)
				suite.Equal(111.1, order.TotalPrice)
				suite.Nil(err)
			}
		})
	}
}

func (suite *OrderServiceTestSuite) TestGetMyOrders() {
	tests := []struct {
		name    string
		setup   func()
		wantErr bool
	}{
		{
			name: "Success",
			setup: func() {
				suite.mockRepo.On("GetMyOrders", mock.Anything, mock.Anything).
					Return([]*model.Order{{UserID: "userID", TotalPrice: 111.2, Status: model.OrderStatusNew}},
						&paging.Pagination{Total: 1, CurrentPage: 1, Limit: 10}, nil).Times(1)
			},
		},
		{
			name: "Fail",
			setup: func() {
				suite.mockRepo.On("GetMyOrders", mock.Anything, mock.Anything).
					Return(nil, nil, errors.New("error")).Times(1)
			},
			wantErr: true,
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.setup()
			req := &domain.ListOrderReq{Status: "new"}
			orders, pagination, err := suite.service.GetMyOrders(context.Background(), req)
			if tc.wantErr {
				suite.Nil(orders)
				suite.Nil(pagination)
				suite.NotNil(err)
			} else {
				suite.NotNil(orders)
				suite.Equal(1, len(orders))
				suite.NotNil(pagination)
				suite.Nil(err)
			}
		})
	}
}

func (suite *OrderServiceTestSuite) TestPlaceOrder() {
	tests := []struct {
		name    string
		req     *domain.PlaceOrderReq
		setup   func()
		wantErr bool
		sleep   bool
	}{
		{
			name: "Success",
			req: &domain.PlaceOrderReq{
				UserID: "userID",
				Lines:  []domain.PlaceOrderLineReq{{ProductID: "productID", Quantity: 2}},
			},
			setup: func() {
				suite.mockProductRepo.On("GetProductByID", mock.Anything, "productID").
					Return(&model.Product{Name: "product", Price: 1.1}, nil).Times(1)
				suite.mockRepo.On("CreateOrder", mock.Anything, "userID", mock.Anything, "", float64(0)).
					Return(&model.Order{UserID: "userID", Lines: []*model.OrderLine{{ProductID: "productID", Quantity: 2}}}, nil).Times(1)
				suite.mockProductRepo.On("DecrementStock", mock.Anything, "productID", 2).Return(nil).Maybe()
				suite.mockCartRepo.On("ClearCart", mock.Anything, "userID").Return(nil).Maybe()
				suite.mockUserRepo.On("GetUserByID", mock.Anything, "userID").
					Return(&model.User{ID: "userID", Email: "user@test.com"}, nil).Maybe()
				suite.mockNotifier.On("SendOrderPlaced", mock.Anything, mock.Anything, "user@test.com").Return(nil).Maybe()
			},
		},
		{
			name: "GetProductByID fail",
			req: &domain.PlaceOrderReq{
				UserID: "userID",
				Lines:  []domain.PlaceOrderLineReq{{ProductID: "productID", Quantity: 2}},
			},
			setup: func() {
				suite.mockProductRepo.On("GetProductByID", mock.Anything, "productID").
					Return(nil, errors.New("error")).Times(1)
			},
			wantErr: true,
		},
		{
			name: "Missing UserID",
			req: &domain.PlaceOrderReq{
				Lines: []domain.PlaceOrderLineReq{{ProductID: "productID", Quantity: 2}},
			},
			setup:   func() {},
			wantErr: true,
		},
		{
			name: "CreateOrder fail",
			req: &domain.PlaceOrderReq{
				UserID: "userID",
				Lines:  []domain.PlaceOrderLineReq{{ProductID: "productID", Quantity: 2}},
			},
			setup: func() {
				suite.mockProductRepo.On("GetProductByID", mock.Anything, "productID").
					Return(&model.Product{Name: "product", Price: 1.1}, nil).Times(1)
				suite.mockRepo.On("CreateOrder", mock.Anything, "userID", mock.Anything, "", float64(0)).
					Return(nil, errors.New("error")).Times(1)
			},
			wantErr: true,
		},
		{
			name: "With coupon",
			req: &domain.PlaceOrderReq{
				UserID: "userID", CouponCode: "SAVE10",
				Lines: []domain.PlaceOrderLineReq{{ProductID: "productID", Quantity: 2}},
			},
			setup: func() {
				suite.mockProductRepo.On("GetProductByID", mock.Anything, "productID").
					Return(&model.Product{Name: "product", Price: 10.0}, nil).Times(1)
				suite.mockCouponSvc.On("Apply", mock.Anything, "SAVE10", float64(20)).
					Return(float64(2), &model.Coupon{ID: "c1", Code: "SAVE10"}, nil).Times(1)
				suite.mockCouponSvc.On("IncrUsedCount", mock.Anything, "c1").Return(nil).Times(1)
				suite.mockRepo.On("CreateOrder", mock.Anything, "userID", mock.Anything, "SAVE10", float64(2)).
					Return(&model.Order{UserID: "userID", TotalPrice: 20, DiscountAmount: 2, FinalPrice: 18, CouponCode: "SAVE10"}, nil).Times(1)
				suite.mockProductRepo.On("DecrementStock", mock.Anything, "productID", 2).Return(nil).Maybe()
				suite.mockCartRepo.On("ClearCart", mock.Anything, "userID").Return(nil).Maybe()
				suite.mockUserRepo.On("GetUserByID", mock.Anything, "userID").
					Return(&model.User{ID: "userID", Email: "user@test.com"}, nil).Maybe()
				suite.mockNotifier.On("SendOrderPlaced", mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()
			},
		},
		{
			name: "Invalid coupon",
			req: &domain.PlaceOrderReq{
				UserID: "userID", CouponCode: "INVALID",
				Lines: []domain.PlaceOrderLineReq{{ProductID: "productID", Quantity: 2}},
			},
			setup: func() {
				suite.mockProductRepo.On("GetProductByID", mock.Anything, "productID").
					Return(&model.Product{Name: "product", Price: 10.0}, nil).Times(1)
				suite.mockCouponSvc.On("Apply", mock.Anything, "INVALID", float64(20)).
					Return(float64(0), (*model.Coupon)(nil), errors.New("coupon not found")).Times(1)
			},
			wantErr: true,
		},
		{
			name: "ClearCart fail (best-effort)",
			req: &domain.PlaceOrderReq{
				UserID: "userID",
				Lines:  []domain.PlaceOrderLineReq{{ProductID: "productID", Quantity: 2}},
			},
			setup: func() {
				suite.mockProductRepo.On("GetProductByID", mock.Anything, "productID").
					Return(&model.Product{Name: "product", Price: 1.1}, nil).Times(1)
				suite.mockRepo.On("CreateOrder", mock.Anything, "userID", mock.Anything, "", float64(0)).
					Return(&model.Order{UserID: "userID", Lines: []*model.OrderLine{{ProductID: "productID", Quantity: 2}}}, nil).Times(1)
				suite.mockProductRepo.On("DecrementStock", mock.Anything, "productID", 2).Return(nil).Maybe()
				suite.mockCartRepo.On("ClearCart", mock.Anything, "userID").Return(errors.New("clear error")).Times(1)
				suite.mockUserRepo.On("GetUserByID", mock.Anything, "userID").
					Return(&model.User{ID: "userID", Email: "user@test.com"}, nil).Maybe()
				suite.mockNotifier.On("SendOrderPlaced", mock.Anything, mock.Anything, "user@test.com").Return(nil).Maybe()
			},
		},
		{
			name:    "DecrementStock fail",
			wantErr: true,
			req: &domain.PlaceOrderReq{
				UserID: "userID",
				Lines:  []domain.PlaceOrderLineReq{{ProductID: "productID", Quantity: 2}},
			},
			setup: func() {
				suite.mockProductRepo.On("GetProductByID", mock.Anything, "productID").
					Return(&model.Product{Name: "product", Price: 1.1}, nil).Times(1)
				suite.mockRepo.On("CreateOrder", mock.Anything, "userID", mock.Anything, "", float64(0)).
					Return(&model.Order{UserID: "userID", Lines: []*model.OrderLine{{ProductID: "productID", Quantity: 2}}}, nil).Times(1)
				suite.mockProductRepo.On("DecrementStock", mock.Anything, "productID", 2).Return(errors.New("stock error")).Times(1)
				suite.mockCartRepo.On("ClearCart", mock.Anything, "userID").Return(nil).Maybe()
			},
		},
		{
			name: "Coupon IncrUsedCount fail",
			req: &domain.PlaceOrderReq{
				UserID: "userID", CouponCode: "SAVE10",
				Lines: []domain.PlaceOrderLineReq{{ProductID: "productID", Quantity: 2}},
			},
			setup: func() {
				suite.mockProductRepo.On("GetProductByID", mock.Anything, "productID").
					Return(&model.Product{Name: "product", Price: 10.0}, nil).Times(1)
				suite.mockCouponSvc.On("Apply", mock.Anything, "SAVE10", float64(20)).
					Return(float64(2), &model.Coupon{ID: "c1", Code: "SAVE10"}, nil).Times(1)
				suite.mockCouponSvc.On("IncrUsedCount", mock.Anything, "c1").Return(errors.New("incr error")).Times(1)
				suite.mockRepo.On("CreateOrder", mock.Anything, "userID", mock.Anything, "SAVE10", float64(2)).
					Return(&model.Order{UserID: "userID", TotalPrice: 20, DiscountAmount: 2, FinalPrice: 18, CouponCode: "SAVE10",
						Lines: []*model.OrderLine{{ProductID: "productID", Quantity: 2}}}, nil).Times(1)
				suite.mockProductRepo.On("DecrementStock", mock.Anything, "productID", 2).Return(nil).Maybe()
				suite.mockCartRepo.On("ClearCart", mock.Anything, "userID").Return(nil).Maybe()
				suite.mockUserRepo.On("GetUserByID", mock.Anything, "userID").
					Return(&model.User{ID: "userID", Email: "user@test.com"}, nil).Maybe()
				suite.mockNotifier.On("SendOrderPlaced", mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()
			},
		},
		{
			name: "GetUser for notification fail",
			req: &domain.PlaceOrderReq{
				UserID: "userID",
				Lines:  []domain.PlaceOrderLineReq{{ProductID: "productID", Quantity: 2}},
			},
			setup: func() {
				suite.mockProductRepo.On("GetProductByID", mock.Anything, "productID").
					Return(&model.Product{Name: "product", Price: 1.1}, nil).Times(1)
				suite.mockRepo.On("CreateOrder", mock.Anything, "userID", mock.Anything, "", float64(0)).
					Return(&model.Order{UserID: "userID", Lines: []*model.OrderLine{{ProductID: "productID", Quantity: 2}}}, nil).Times(1)
				suite.mockProductRepo.On("DecrementStock", mock.Anything, "productID", 2).Return(nil).Maybe()
				suite.mockCartRepo.On("ClearCart", mock.Anything, "userID").Return(nil).Maybe()
				suite.mockUserRepo.On("GetUserByID", mock.Anything, "userID").
					Return(nil, errors.New("user not found")).Times(1)
			},
			sleep: true,
		},
		{
			name: "SendOrderPlaced notification fail",
			req: &domain.PlaceOrderReq{
				UserID: "userID",
				Lines:  []domain.PlaceOrderLineReq{{ProductID: "productID", Quantity: 2}},
			},
			setup: func() {
				suite.mockProductRepo.On("GetProductByID", mock.Anything, "productID").
					Return(&model.Product{Name: "product", Price: 1.1}, nil).Times(1)
				suite.mockRepo.On("CreateOrder", mock.Anything, "userID", mock.Anything, "", float64(0)).
					Return(&model.Order{UserID: "userID", Lines: []*model.OrderLine{{ProductID: "productID", Quantity: 2}}}, nil).Times(1)
				suite.mockProductRepo.On("DecrementStock", mock.Anything, "productID", 2).Return(nil).Maybe()
				suite.mockCartRepo.On("ClearCart", mock.Anything, "userID").Return(nil).Maybe()
				suite.mockUserRepo.On("GetUserByID", mock.Anything, "userID").
					Return(&model.User{ID: "userID", Email: "user@test.com"}, nil).Times(1)
				suite.mockNotifier.On("SendOrderPlaced", mock.Anything, mock.Anything, "user@test.com").
					Return(errors.New("notif error")).Times(1)
			},
			sleep: true,
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.setup()
			order, err := suite.service.PlaceOrder(context.Background(), tc.req)
			if tc.wantErr {
				suite.Nil(order)
				suite.NotNil(err)
			} else {
				suite.NotNil(order)
				suite.Nil(err)
			}
			if tc.sleep {
				time.Sleep(20 * time.Millisecond)
			}
		})
	}
}

func (suite *OrderServiceTestSuite) TestCancelOrder() {
	tests := []struct {
		name    string
		setup   func()
		wantErr bool
	}{
		{
			name: "Success",
			setup: func() {
				suite.mockRepo.On("GetOrderByID", mock.Anything, "orderID", false).
					Return(&model.Order{UserID: "userID", TotalPrice: 111.1, Status: model.OrderStatusNew}, nil).Times(1)
				suite.mockRepo.On("UpdateOrder", mock.Anything, &model.Order{
					UserID: "userID", TotalPrice: 111.1, Status: model.OrderStatusCancelled,
				}).Return(nil).Times(1)
			},
		},
		{
			name: "Update fail",
			setup: func() {
				suite.mockRepo.On("GetOrderByID", mock.Anything, "orderID", false).
					Return(&model.Order{UserID: "userID", TotalPrice: 111.1, Status: model.OrderStatusNew}, nil).Times(1)
				suite.mockRepo.On("UpdateOrder", mock.Anything, &model.Order{
					UserID: "userID", TotalPrice: 111.1, Status: model.OrderStatusCancelled,
				}).Return(errors.New("error")).Times(1)
			},
			wantErr: true,
		},
		{
			name: "Different user",
			setup: func() {
				suite.mockRepo.On("GetOrderByID", mock.Anything, "orderID", false).
					Return(&model.Order{UserID: "userID1", TotalPrice: 111.1, Status: model.OrderStatusNew}, nil).Times(1)
			},
			wantErr: true,
		},
		{
			name: "Invalid status",
			setup: func() {
				suite.mockRepo.On("GetOrderByID", mock.Anything, "orderID", false).
					Return(&model.Order{UserID: "userID", TotalPrice: 111.1, Status: model.OrderStatusCancelled}, nil).Times(1)
			},
			wantErr: true,
		},
		{
			name: "GetOrderByID fail",
			setup: func() {
				suite.mockRepo.On("GetOrderByID", mock.Anything, "orderID", false).
					Return(nil, errors.New("error")).Times(1)
			},
			wantErr: true,
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.setup()
			order, err := suite.service.CancelOrder(context.Background(), "orderID", "userID")
			if tc.wantErr {
				suite.Nil(order)
				suite.NotNil(err)
			} else {
				suite.NotNil(order)
				suite.Equal(model.OrderStatusCancelled, order.Status)
				suite.Nil(err)
			}
		})
	}
}

func (suite *OrderServiceTestSuite) TestUpdateOrderStatus() {
	tests := []struct {
		name    string
		setup   func()
		wantErr bool
		sleep   bool
	}{
		{
			name: "Success",
			setup: func() {
				suite.mockRepo.On("GetOrderByID", mock.Anything, "orderID", false).
					Return(&model.Order{ID: "orderID", UserID: "userID", Status: model.OrderStatusNew}, nil).Times(1)
				suite.mockRepo.On("UpdateOrder", mock.Anything, &model.Order{
					ID: "orderID", UserID: "userID", Status: model.OrderStatusDone,
				}).Return(nil).Times(1)
				suite.mockUserRepo.On("GetUserByID", mock.Anything, "userID").
					Return(&model.User{ID: "userID", Email: "user@test.com"}, nil).Maybe()
				suite.mockNotifier.On("SendOrderStatusChanged", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(nil).Maybe()
			},
		},
		{
			name: "GetOrder fail",
			setup: func() {
				suite.mockRepo.On("GetOrderByID", mock.Anything, "orderID", false).
					Return(nil, errors.New("not found")).Times(1)
			},
			wantErr: true,
		},
		{
			name: "Update fail",
			setup: func() {
				suite.mockRepo.On("GetOrderByID", mock.Anything, "orderID", false).
					Return(&model.Order{ID: "orderID", UserID: "userID", Status: model.OrderStatusNew}, nil).Times(1)
				suite.mockRepo.On("UpdateOrder", mock.Anything, &model.Order{
					ID: "orderID", UserID: "userID", Status: model.OrderStatusDone,
				}).Return(errors.New("db error")).Times(1)
			},
			wantErr: true,
		},
		{
			name: "GetUser fail (goroutine)",
			setup: func() {
				suite.mockRepo.On("GetOrderByID", mock.Anything, "orderID", false).
					Return(&model.Order{ID: "orderID", UserID: "userID", Status: model.OrderStatusNew}, nil).Times(1)
				suite.mockRepo.On("UpdateOrder", mock.Anything, &model.Order{
					ID: "orderID", UserID: "userID", Status: model.OrderStatusDone,
				}).Return(nil).Times(1)
				suite.mockUserRepo.On("GetUserByID", mock.Anything, "userID").
					Return(nil, errors.New("user not found")).Times(1)
			},
			sleep: true,
		},
		{
			name: "Send notification fail (goroutine)",
			setup: func() {
				suite.mockRepo.On("GetOrderByID", mock.Anything, "orderID", false).
					Return(&model.Order{ID: "orderID", UserID: "userID", Status: model.OrderStatusNew}, nil).Times(1)
				suite.mockRepo.On("UpdateOrder", mock.Anything, &model.Order{
					ID: "orderID", UserID: "userID", Status: model.OrderStatusDone,
				}).Return(nil).Times(1)
				suite.mockUserRepo.On("GetUserByID", mock.Anything, "userID").
					Return(&model.User{ID: "userID", Email: "user@test.com"}, nil).Times(1)
				suite.mockNotifier.On("SendOrderStatusChanged", mock.Anything, mock.Anything, "user@test.com", string(model.OrderStatusDone)).
					Return(errors.New("notif error")).Times(1)
			},
			sleep: true,
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.setup()
			order, err := suite.service.UpdateOrderStatus(context.Background(), "orderID", model.OrderStatusDone)
			if tc.wantErr {
				suite.Nil(order)
				suite.NotNil(err)
			} else {
				suite.NotNil(order)
				suite.Equal(model.OrderStatusDone, order.Status)
				suite.Nil(err)
			}
			if tc.sleep {
				time.Sleep(20 * time.Millisecond)
			}
		})
	}
}
