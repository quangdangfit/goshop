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

func (suite *OrderServiceTestSuite) TestPlaceOrderCreateOrderFail() {
	req := &serializers.PlaceOrderReq{
		Lines: []serializers.PlaceOrderLineReq{
			{
				ProductID: "productID",
				Quantity:  3,
			},
		},
	}
	suite.mockProductRepo.EXPECT().GetProductByID(gomock.Any(), gomock.Any()).Return(&models.Product{}, nil).Times(1)
	suite.mockRepo.EXPECT().CreateOrder(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("create order fail")).Times(1)

	order, err := suite.service.PlaceOrder(context.TODO(), req)
	suite.NotNil(err)
	suite.Nil(order)
}

func (suite *OrderServiceTestSuite) TestGetMyOrdersFail() {
	suite.mockRepo.EXPECT().GetMyOrders(gomock.Any(), gomock.Any()).Return(nil, nil, errors.New("list my orders fail")).Times(1)

	orders, pagination, err := suite.service.GetMyOrders(context.TODO(), nil)
	suite.NotNil(err)
	suite.Nil(orders)
	suite.Nil(pagination)
}

func (suite *OrderServiceTestSuite) TestCancelOrderUpdateFail() {
	suite.mockRepo.EXPECT().GetOrderByID(gomock.Any(), gomock.Any(), gomock.Any()).Return(&models.Order{
		Status: models.OrderStatusNew,
		UserID: "userID",
	}, nil).Times(1)
	suite.mockRepo.EXPECT().UpdateOrder(gomock.Any(), gomock.Any()).Return(errors.New("update order fail")).Times(1)

	order, err := suite.service.CancelOrder(context.TODO(), "orderID", "userID")
	suite.NotNil(err)
	suite.Nil(order)
}
