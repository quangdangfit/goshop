package service

import (
	"context"
	"errors"
	"testing"

	"github.com/quangdangfit/gocommon/logger"
	"github.com/quangdangfit/gocommon/validation"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"goshop/internal/order/dto"
	"goshop/internal/order/model"
	"goshop/internal/order/repository/mocks"
	"goshop/pkg/config"
	"goshop/pkg/paging"
)

type OrderServiceTestSuite struct {
	suite.Suite
	mockRepo        *mocks.IOrderRepository
	mockProductRepo *mocks.IProductRepository
	service         IOrderService
}

func (suite *OrderServiceTestSuite) SetupTest() {
	logger.Initialize(config.ProductionEnv)

	validator := validation.New()
	suite.mockRepo = mocks.NewIOrderRepository(suite.T())
	suite.mockProductRepo = mocks.NewIProductRepository(suite.T())
	suite.service = NewOrderService(validator, suite.mockRepo, suite.mockProductRepo)
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
			{
				ProductID: "productID",
				Quantity:  2,
			},
		},
	}

	suite.mockProductRepo.On("GetProductByID", mock.Anything, "productID").
		Return(&model.Product{
			Name:        "product",
			Description: "product description",
			Price:       1.1,
		}, nil).Times(1)

	suite.mockRepo.On("CreateOrder", mock.Anything, "userID", mock.Anything).
		Return(&model.Order{
			UserID: "userID",
			Lines: []*model.OrderLine{
				{
					ProductID: "productID",
					Quantity:  2,
				},
			},
		}, nil).Times(1)

	order, err := suite.service.PlaceOrder(context.Background(), req)
	suite.NotNil(order)
	suite.Equal(req.UserID, order.UserID)
	suite.Equal(1, len(order.Lines))
	suite.Equal(req.Lines[0].ProductID, order.Lines[0].ProductID)
	suite.Equal(req.Lines[0].Quantity, order.Lines[0].Quantity)
	suite.Nil(err)
}

func (suite *OrderServiceTestSuite) TestPlaceOrderGetProductByIDFail() {
	req := &dto.PlaceOrderReq{
		UserID: "userID",
		Lines: []dto.PlaceOrderLineReq{
			{
				ProductID: "productID",
				Quantity:  2,
			},
		},
	}

	suite.mockProductRepo.On("GetProductByID", mock.Anything, "productID", mock.Anything).
		Return(nil, errors.New("error")).Times(1)

	order, err := suite.service.PlaceOrder(context.Background(), req)
	suite.Nil(order)
	suite.NotNil(err)
}

func (suite *OrderServiceTestSuite) TestPlaceOrderMissUserId() {
	req := &dto.PlaceOrderReq{
		Lines: []dto.PlaceOrderLineReq{
			{
				ProductID: "productID",
				Quantity:  2,
			},
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
			{
				ProductID: "productID",
				Quantity:  2,
			},
		},
	}

	suite.mockProductRepo.On("GetProductByID", mock.Anything, "productID").
		Return(&model.Product{
			Name:        "product",
			Description: "product description",
			Price:       1.1,
		}, nil).Times(1)

	suite.mockRepo.On("CreateOrder", mock.Anything, "userID", mock.Anything).
		Return(nil, errors.New("error")).Times(1)

	order, err := suite.service.PlaceOrder(context.Background(), req)
	suite.Nil(order)
	suite.NotNil(err)
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
	suite.Equal(111.1, order.TotalPrice)
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
