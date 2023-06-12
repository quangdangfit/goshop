package services

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/quangdangfit/gocommon/logger"
	"github.com/stretchr/testify/suite"

	"goshop/app/models"
	"goshop/app/serializers"
	"goshop/config"
	"goshop/mocks"
	"goshop/pkg/paging"
)

type OrderServiceTestSuite struct {
	suite.Suite
	mockRepo        *mocks.MockIOrderRepository
	mockProductRepo *mocks.MockIProductRepository
	service         IOrderService
}

func (suite *OrderServiceTestSuite) SetupTest() {
	logger.Initialize(config.TestEnv)

	mockCtrl := gomock.NewController(suite.T())
	defer mockCtrl.Finish()
	suite.mockRepo = mocks.NewMockIOrderRepository(mockCtrl)
	suite.mockProductRepo = mocks.NewMockIProductRepository(mockCtrl)
	suite.service = NewOrderService(suite.mockRepo, suite.mockProductRepo)
}

func TestOrderServiceTestSuite(t *testing.T) {
	suite.Run(t, new(OrderServiceTestSuite))
}

// GetOrderByID
// =================================================================

func (suite *OrderServiceTestSuite) TestGetOrderByIDSuccess() {
	orderID := "orderID"
	suite.mockRepo.EXPECT().GetOrderByID(gomock.Any(), orderID, true).Return(&models.Order{
		UserID:     "userID",
		TotalPrice: 111.1,
		Status:     models.OrderStatusNew,
	}, nil).Times(1)

	order, err := suite.service.GetOrderByID(context.Background(), orderID)
	suite.NotNil(order)
	suite.Equal("userID", order.UserID)
	suite.Equal(111.1, order.TotalPrice)
	suite.Equal(models.OrderStatusNew, order.Status)
	suite.Nil(err)
}

func (suite *OrderServiceTestSuite) TestGetOrderByIDFail() {
	orderID := "orderID"
	suite.mockRepo.EXPECT().GetOrderByID(gomock.Any(), orderID, true).Return(nil, errors.New("error")).Times(1)

	order, err := suite.service.GetOrderByID(context.Background(), orderID)
	suite.Nil(order)
	suite.NotNil(err)
}

// GetMyOrders
// =================================================================

func (suite *OrderServiceTestSuite) TestListOrdersSuccess() {
	req := &serializers.ListOrderReq{
		Status: "new",
	}

	suite.mockRepo.EXPECT().GetMyOrders(gomock.Any(), req).Return(
		[]*models.Order{
			{
				UserID:     "userID",
				TotalPrice: 111.2,
				Status:     models.OrderStatusNew,
			},
		},
		&paging.Pagination{
			Total:       1,
			CurrentPage: 1,
			Limit:       10,
		},
		nil).Times(1)

	orders, pagination, err := suite.service.GetMyOrders(context.Background(), req)
	suite.NotNil(orders)
	suite.Equal(1, len(orders))
	suite.Equal("userID", orders[0].UserID)
	suite.Equal(111.2, orders[0].TotalPrice)
	suite.Equal(models.OrderStatusNew, orders[0].Status)
	suite.NotNil(pagination)
	suite.Equal(int64(1), pagination.Total)
	suite.Equal(int64(1), pagination.CurrentPage)
	suite.Equal(int64(10), pagination.Limit)
	suite.Nil(err)
}

func (suite *OrderServiceTestSuite) TestListOrdersFail() {
	req := &serializers.ListOrderReq{
		Status: "new",
	}

	suite.mockRepo.EXPECT().GetMyOrders(gomock.Any(), req).Return(nil, nil, errors.New("error")).Times(1)

	orders, pagination, err := suite.service.GetMyOrders(context.Background(), req)
	suite.Nil(orders)
	suite.Nil(pagination)
	suite.NotNil(err)
}

// Place Order
// =================================================================

func (suite *OrderServiceTestSuite) TestPlaceOrderSuccess() {
	req := &serializers.PlaceOrderReq{
		UserID: "userID",
		Lines: []serializers.PlaceOrderLineReq{
			{
				ProductID: "productID",
				Quantity:  2,
			},
		},
	}

	suite.mockProductRepo.EXPECT().GetProductByID(gomock.Any(), "productID").Return(&models.Product{
		Name:        "product",
		Description: "product description",
		Price:       1.1,
	}, nil).Times(1)
	suite.mockRepo.EXPECT().CreateOrder(gomock.Any(), "userID", gomock.Any()).Return(
		&models.Order{
			UserID: "userID",
			Lines: []*models.OrderLine{
				{
					ProductID: "productID",
					Quantity:  2,
				},
			},
		},
		nil).Times(1)

	order, err := suite.service.PlaceOrder(context.Background(), req)
	suite.NotNil(order)
	suite.Equal(req.UserID, order.UserID)
	suite.Equal(1, len(order.Lines))
	suite.Equal(req.Lines[0].ProductID, order.Lines[0].ProductID)
	suite.Equal(req.Lines[0].Quantity, order.Lines[0].Quantity)
	suite.Nil(err)
}

func (suite *OrderServiceTestSuite) TestCreateGetProductByIDFail() {
	req := &serializers.PlaceOrderReq{
		UserID: "userID",
		Lines: []serializers.PlaceOrderLineReq{
			{
				ProductID: "productID",
				Quantity:  2,
			},
		},
	}

	suite.mockProductRepo.EXPECT().GetProductByID(gomock.Any(), "productID").Return(nil, errors.New("error")).Times(1)

	order, err := suite.service.PlaceOrder(context.Background(), req)
	suite.Nil(order)
	suite.NotNil(err)
}

func (suite *OrderServiceTestSuite) TestCreateFail() {
	req := &serializers.PlaceOrderReq{
		UserID: "userID",
		Lines: []serializers.PlaceOrderLineReq{
			{
				ProductID: "productID",
				Quantity:  2,
			},
		},
	}

	suite.mockProductRepo.EXPECT().GetProductByID(gomock.Any(), "productID").Return(&models.Product{
		Name:        "product",
		Description: "product description",
		Price:       1.1,
	}, nil).Times(1)
	suite.mockRepo.EXPECT().CreateOrder(gomock.Any(), "userID", gomock.Any()).Return(nil, errors.New("error")).Times(1)

	order, err := suite.service.PlaceOrder(context.Background(), req)
	suite.Nil(order)
	suite.NotNil(err)
}

// Cancel Order
// =================================================================

func (suite *OrderServiceTestSuite) TestCancelOrderSuccess() {
	userID := "orderID"
	orderID := "orderID"
	suite.mockRepo.EXPECT().GetOrderByID(gomock.Any(), orderID, false).Return(
		&models.Order{
			UserID:     userID,
			TotalPrice: 111.1,
			Status:     models.OrderStatusNew,
		},
		nil).Times(1)
	suite.mockRepo.EXPECT().UpdateOrder(gomock.Any(), &models.Order{
		UserID:     userID,
		TotalPrice: 111.1,
		Status:     models.OrderStatusCancelled,
	}).Return(nil).Times(1)

	order, err := suite.service.CancelOrder(context.Background(), orderID, userID)
	suite.NotNil(order)
	suite.Equal(userID, order.UserID)
	suite.Equal(111.1, order.TotalPrice)
	suite.Equal(models.OrderStatusCancelled, order.Status)
	suite.Nil(err)
}

func (suite *OrderServiceTestSuite) TestCancelOrderFail() {
	userID := "orderID"
	orderID := "orderID"
	suite.mockRepo.EXPECT().GetOrderByID(gomock.Any(), orderID, false).Return(
		&models.Order{
			UserID:     userID,
			TotalPrice: 111.1,
			Status:     models.OrderStatusNew,
		},
		nil).Times(1)
	suite.mockRepo.EXPECT().UpdateOrder(gomock.Any(), &models.Order{
		UserID:     userID,
		TotalPrice: 111.1,
		Status:     models.OrderStatusCancelled,
	}).Return(errors.New("error")).Times(1)

	order, err := suite.service.CancelOrder(context.Background(), orderID, userID)
	suite.Nil(order)
	suite.NotNil(err)
}

func (suite *OrderServiceTestSuite) TestCancelOrderGetOrderByIDFail() {
	userID := "orderID"
	orderID := "orderID"
	suite.mockRepo.EXPECT().GetOrderByID(gomock.Any(), orderID, false).Return(nil, errors.New("error")).Times(1)

	order, err := suite.service.CancelOrder(context.Background(), orderID, userID)
	suite.Nil(order)
	suite.NotNil(err)
}
