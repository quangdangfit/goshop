package service

import (
	"context"
	"errors"
	"testing"

	"github.com/quangdangfit/gocommon/logger"
	"github.com/quangdangfit/gocommon/validation"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	cartMocks "goshop/internal/cart/repository/mocks"
	"goshop/internal/order/dto"
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

// GetOrderByID
// =================================================================

func (suite *OrderServiceTestSuite) TestGetOrderByIDSuccess() {
	orderID := "orderID"
	suite.mockRepo.On("GetOrderByID", mock.Anything, orderID, true).
		Return(&model.Order{
			UserID:     "userID",
			TotalPrice: 111.1,
			Status:     model.OrderStatusNew,
		}, nil).Times(1)

	order, err := suite.service.GetOrderByID(context.Background(), orderID)
	suite.NotNil(order)
	suite.Equal("userID", order.UserID)
	suite.Equal(111.1, order.TotalPrice)
	suite.Equal(model.OrderStatusNew, order.Status)
	suite.Nil(err)
}

func (suite *OrderServiceTestSuite) TestGetOrderByIDFail() {
	orderID := "orderID"
	suite.mockRepo.On("GetOrderByID", mock.Anything, orderID, true).
		Return(nil, errors.New("error")).Times(1)

	order, err := suite.service.GetOrderByID(context.Background(), orderID)
	suite.Nil(order)
	suite.NotNil(err)
}

// GetMyOrders
// =================================================================

func (suite *OrderServiceTestSuite) TestListOrdersSuccess() {
	req := &dto.ListOrderReq{
		Status: "new",
	}

	suite.mockRepo.On("GetMyOrders", mock.Anything, req).
		Return(
			[]*model.Order{
				{
					UserID:     "userID",
					TotalPrice: 111.2,
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

	orders, pagination, err := suite.service.GetMyOrders(context.Background(), req)
	suite.NotNil(orders)
	suite.Equal(1, len(orders))
	suite.Equal("userID", orders[0].UserID)
	suite.Equal(111.2, orders[0].TotalPrice)
	suite.Equal(model.OrderStatusNew, orders[0].Status)
	suite.NotNil(pagination)
	suite.Equal(int64(1), pagination.Total)
	suite.Equal(int64(1), pagination.CurrentPage)
	suite.Equal(int64(10), pagination.Limit)
	suite.Nil(err)
}

func (suite *OrderServiceTestSuite) TestListOrdersFail() {
	req := &dto.ListOrderReq{
		Status: "new",
	}

	suite.mockRepo.On("GetMyOrders", mock.Anything, req).
		Return(nil, nil, errors.New("error")).Times(1)

	orders, pagination, err := suite.service.GetMyOrders(context.Background(), req)
	suite.Nil(orders)
	suite.Nil(pagination)
	suite.NotNil(err)
}

// Place Order
// =================================================================

func (suite *OrderServiceTestSuite) TestPlaceOrderSuccess() {
	req := &dto.PlaceOrderReq{
		UserID: "userID",
		Lines: []dto.PlaceOrderLineReq{
			{ProductID: "productID", Quantity: 2},
		},
	}

	suite.mockProductRepo.On("GetProductByID", mock.Anything, "productID").
		Return(&model.Product{Name: "product", Price: 1.1}, nil).Times(1)

	suite.mockRepo.On("CreateOrder", mock.Anything, "userID", mock.Anything, "", float64(0)).
		Return(&model.Order{
			UserID: "userID",
			Lines:  []*model.OrderLine{{ProductID: "productID", Quantity: 2}},
		}, nil).Times(1)

	// DecrementStock is called best-effort
	suite.mockProductRepo.On("DecrementStock", mock.Anything, "productID", 2).
		Return(nil).Maybe()

	// ClearCart is called best-effort after order creation
	suite.mockCartRepo.On("ClearCart", mock.Anything, "userID").Return(nil).Maybe()

	// Notification goroutine (best-effort, may or may not complete before test ends)
	suite.mockUserRepo.On("GetUserByID", mock.Anything, "userID").
		Return(&model.User{ID: "userID", Email: "user@test.com"}, nil).Maybe()
	suite.mockNotifier.On("SendOrderPlaced", mock.Anything, mock.Anything, "user@test.com").
		Return(nil).Maybe()

	order, err := suite.service.PlaceOrder(context.Background(), req)
	suite.NotNil(order)
	suite.Equal(req.UserID, order.UserID)
	suite.Equal(1, len(order.Lines))
	suite.Nil(err)
}

func (suite *OrderServiceTestSuite) TestPlaceOrderGetProductByIDFail() {
	req := &dto.PlaceOrderReq{
		UserID: "userID",
		Lines: []dto.PlaceOrderLineReq{
			{ProductID: "productID", Quantity: 2},
		},
	}

	suite.mockProductRepo.On("GetProductByID", mock.Anything, "productID").
		Return(nil, errors.New("error")).Times(1)

	order, err := suite.service.PlaceOrder(context.Background(), req)
	suite.Nil(order)
	suite.NotNil(err)
}

func (suite *OrderServiceTestSuite) TestPlaceOrderMissUserId() {
	req := &dto.PlaceOrderReq{
		Lines: []dto.PlaceOrderLineReq{
			{ProductID: "productID", Quantity: 2},
		},
	}

	order, err := suite.service.PlaceOrder(context.Background(), req)
	suite.Nil(order)
	suite.NotNil(err)
}

func (suite *OrderServiceTestSuite) TestPlaceOrderCreateFail() {
	req := &dto.PlaceOrderReq{
		UserID: "userID",
		Lines: []dto.PlaceOrderLineReq{
			{ProductID: "productID", Quantity: 2},
		},
	}

	suite.mockProductRepo.On("GetProductByID", mock.Anything, "productID").
		Return(&model.Product{Name: "product", Price: 1.1}, nil).Times(1)

	suite.mockRepo.On("CreateOrder", mock.Anything, "userID", mock.Anything, "", float64(0)).
		Return(nil, errors.New("error")).Times(1)

	order, err := suite.service.PlaceOrder(context.Background(), req)
	suite.Nil(order)
	suite.NotNil(err)
}

func (suite *OrderServiceTestSuite) TestPlaceOrderWithCoupon() {
	req := &dto.PlaceOrderReq{
		UserID:     "userID",
		CouponCode: "SAVE10",
		Lines: []dto.PlaceOrderLineReq{
			{ProductID: "productID", Quantity: 2},
		},
	}

	suite.mockProductRepo.On("GetProductByID", mock.Anything, "productID").
		Return(&model.Product{Name: "product", Price: 10.0}, nil).Times(1)

	suite.mockCouponSvc.On("Apply", mock.Anything, "SAVE10", float64(20)).
		Return(float64(2), &model.Coupon{ID: "c1", Code: "SAVE10"}, nil).Times(1)

	suite.mockCouponSvc.On("IncrUsedCount", mock.Anything, "c1").
		Return(nil).Times(1)

	suite.mockRepo.On("CreateOrder", mock.Anything, "userID", mock.Anything, "SAVE10", float64(2)).
		Return(&model.Order{
			UserID:         "userID",
			TotalPrice:     20,
			DiscountAmount: 2,
			FinalPrice:     18,
			CouponCode:     "SAVE10",
		}, nil).Times(1)

	suite.mockProductRepo.On("DecrementStock", mock.Anything, "productID", 2).
		Return(nil).Maybe()
	suite.mockCartRepo.On("ClearCart", mock.Anything, "userID").Return(nil).Maybe()
	suite.mockUserRepo.On("GetUserByID", mock.Anything, "userID").
		Return(&model.User{ID: "userID", Email: "user@test.com"}, nil).Maybe()
	suite.mockNotifier.On("SendOrderPlaced", mock.Anything, mock.Anything, mock.Anything).
		Return(nil).Maybe()

	order, err := suite.service.PlaceOrder(context.Background(), req)
	suite.NotNil(order)
	suite.Equal(float64(2), order.DiscountAmount)
	suite.Equal(float64(18), order.FinalPrice)
	suite.Equal("SAVE10", order.CouponCode)
	suite.Nil(err)
}

func (suite *OrderServiceTestSuite) TestPlaceOrderInvalidCoupon() {
	req := &dto.PlaceOrderReq{
		UserID:     "userID",
		CouponCode: "INVALID",
		Lines: []dto.PlaceOrderLineReq{
			{ProductID: "productID", Quantity: 2},
		},
	}

	suite.mockProductRepo.On("GetProductByID", mock.Anything, "productID").
		Return(&model.Product{Name: "product", Price: 10.0}, nil).Times(1)

	suite.mockCouponSvc.On("Apply", mock.Anything, "INVALID", float64(20)).
		Return(float64(0), (*model.Coupon)(nil), errors.New("coupon not found")).Times(1)

	order, err := suite.service.PlaceOrder(context.Background(), req)
	suite.Nil(order)
	suite.NotNil(err)
}

func (suite *OrderServiceTestSuite) TestPlaceOrderCouponIncrFail() {
	req := &dto.PlaceOrderReq{
		UserID:     "userID",
		CouponCode: "SAVE10",
		Lines: []dto.PlaceOrderLineReq{
			{ProductID: "productID", Quantity: 2},
		},
	}

	suite.mockProductRepo.On("GetProductByID", mock.Anything, "productID").
		Return(&model.Product{Name: "product", Price: 10.0}, nil).Times(1)

	suite.mockCouponSvc.On("Apply", mock.Anything, "SAVE10", float64(20)).
		Return(float64(2), &model.Coupon{ID: "c1", Code: "SAVE10"}, nil).Times(1)

	// IncrUsedCount fails but execution continues
	suite.mockCouponSvc.On("IncrUsedCount", mock.Anything, "c1").
		Return(errors.New("incr error")).Times(1)

	suite.mockRepo.On("CreateOrder", mock.Anything, "userID", mock.Anything, "SAVE10", float64(2)).
		Return(&model.Order{
			UserID:         "userID",
			TotalPrice:     20,
			DiscountAmount: 2,
			FinalPrice:     18,
			CouponCode:     "SAVE10",
			Lines:          []*model.OrderLine{{ProductID: "productID", Quantity: 2}},
		}, nil).Times(1)

	suite.mockProductRepo.On("DecrementStock", mock.Anything, "productID", 2).
		Return(nil).Maybe()
	suite.mockCartRepo.On("ClearCart", mock.Anything, "userID").Return(nil).Maybe()
	suite.mockUserRepo.On("GetUserByID", mock.Anything, "userID").
		Return(&model.User{ID: "userID", Email: "user@test.com"}, nil).Maybe()
	suite.mockNotifier.On("SendOrderPlaced", mock.Anything, mock.Anything, mock.Anything).
		Return(nil).Maybe()

	order, err := suite.service.PlaceOrder(context.Background(), req)
	suite.NotNil(order)
	suite.Nil(err)
}

// Cancel Order
// =================================================================

func (suite *OrderServiceTestSuite) TestCancelOrderSuccess() {
	userID := "userID"
	orderID := "orderID"

	suite.mockRepo.On("GetOrderByID", mock.Anything, orderID, false).
		Return(&model.Order{
			UserID:     userID,
			TotalPrice: 111.1,
			Status:     model.OrderStatusNew,
		}, nil).Times(1)

	suite.mockRepo.On("UpdateOrder", mock.Anything, &model.Order{
		UserID:     userID,
		TotalPrice: 111.1,
		Status:     model.OrderStatusCancelled,
	}).Return(nil).Times(1)

	order, err := suite.service.CancelOrder(context.Background(), orderID, userID)
	suite.NotNil(order)
	suite.Equal(userID, order.UserID)
	suite.Equal(model.OrderStatusCancelled, order.Status)
	suite.Nil(err)
}

func (suite *OrderServiceTestSuite) TestCancelOrderFail() {
	userID := "userID"
	orderID := "orderID"

	suite.mockRepo.On("GetOrderByID", mock.Anything, orderID, false).
		Return(&model.Order{
			UserID:     userID,
			TotalPrice: 111.1,
			Status:     model.OrderStatusNew,
		}, nil).Times(1)

	suite.mockRepo.On("UpdateOrder", mock.Anything, &model.Order{
		UserID:     userID,
		TotalPrice: 111.1,
		Status:     model.OrderStatusCancelled,
	}).Return(errors.New("error")).Times(1)

	order, err := suite.service.CancelOrder(context.Background(), orderID, userID)
	suite.Nil(order)
	suite.NotNil(err)
}

func (suite *OrderServiceTestSuite) TestCancelOrderDifferenceUserId() {
	userID := "userID"
	orderID := "orderID"

	suite.mockRepo.On("GetOrderByID", mock.Anything, orderID, false).
		Return(&model.Order{
			UserID:     "userID1",
			TotalPrice: 111.1,
			Status:     model.OrderStatusNew,
		}, nil).Times(1)

	order, err := suite.service.CancelOrder(context.Background(), orderID, userID)
	suite.Nil(order)
	suite.NotNil(err)
}

func (suite *OrderServiceTestSuite) TestCancelOrderInvalidStatus() {
	userID := "userID"
	orderID := "orderID"

	suite.mockRepo.On("GetOrderByID", mock.Anything, orderID, false).
		Return(&model.Order{
			UserID:     userID,
			TotalPrice: 111.1,
			Status:     model.OrderStatusCancelled,
		}, nil).Times(1)

	order, err := suite.service.CancelOrder(context.Background(), orderID, userID)
	suite.Nil(order)
	suite.NotNil(err)
}

func (suite *OrderServiceTestSuite) TestCancelOrderGetOrderByIDFail() {
	userID := "userID"
	orderID := "orderID"

	suite.mockRepo.On("GetOrderByID", mock.Anything, orderID, false).
		Return(nil, errors.New("error")).Times(1)

	order, err := suite.service.CancelOrder(context.Background(), orderID, userID)
	suite.Nil(order)
	suite.NotNil(err)
}

// UpdateOrderStatus
// =================================================================

func (suite *OrderServiceTestSuite) TestUpdateOrderStatusSuccess() {
	orderID := "orderID"

	suite.mockRepo.On("GetOrderByID", mock.Anything, orderID, false).
		Return(&model.Order{
			ID:     orderID,
			UserID: "userID",
			Status: model.OrderStatusNew,
		}, nil).Times(1)

	suite.mockRepo.On("UpdateOrder", mock.Anything, &model.Order{
		ID:     orderID,
		UserID: "userID",
		Status: model.OrderStatusDone,
	}).Return(nil).Times(1)

	suite.mockUserRepo.On("GetUserByID", mock.Anything, "userID").
		Return(&model.User{ID: "userID", Email: "user@test.com"}, nil).Maybe()
	suite.mockNotifier.On("SendOrderStatusChanged", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(nil).Maybe()

	order, err := suite.service.UpdateOrderStatus(context.Background(), orderID, model.OrderStatusDone)
	suite.NotNil(order)
	suite.Equal(model.OrderStatusDone, order.Status)
	suite.Nil(err)
}

func (suite *OrderServiceTestSuite) TestUpdateOrderStatusGetOrderFail() {
	orderID := "orderID"

	suite.mockRepo.On("GetOrderByID", mock.Anything, orderID, false).
		Return(nil, errors.New("not found")).Times(1)

	order, err := suite.service.UpdateOrderStatus(context.Background(), orderID, model.OrderStatusDone)
	suite.Nil(order)
	suite.NotNil(err)
}

func (suite *OrderServiceTestSuite) TestUpdateOrderStatusUpdateFail() {
	orderID := "orderID"

	suite.mockRepo.On("GetOrderByID", mock.Anything, orderID, false).
		Return(&model.Order{
			ID:     orderID,
			UserID: "userID",
			Status: model.OrderStatusNew,
		}, nil).Times(1)

	suite.mockRepo.On("UpdateOrder", mock.Anything, &model.Order{
		ID:     orderID,
		UserID: "userID",
		Status: model.OrderStatusDone,
	}).Return(errors.New("db error")).Times(1)

	order, err := suite.service.UpdateOrderStatus(context.Background(), orderID, model.OrderStatusDone)
	suite.Nil(order)
	suite.NotNil(err)
}
