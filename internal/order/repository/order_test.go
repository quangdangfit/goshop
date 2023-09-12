package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/quangdangfit/gocommon/logger"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"goshop/internal/order/dto"
	"goshop/internal/order/model"
	"goshop/pkg/config"
	"goshop/pkg/dbs/mocks"
)

type OrderRepositoryTestSuite struct {
	suite.Suite
	mockDB *mocks.IDatabase
	repo   IOrderRepository
}

func (suite *OrderRepositoryTestSuite) SetupTest() {
	logger.Initialize(config.ProductionEnv)

	suite.mockDB = mocks.NewIDatabase(suite.T())
	suite.repo = NewOrderRepository(suite.mockDB)
}

func TestOrderRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(OrderRepositoryTestSuite))
}

// CreateOrder
// =================================================================

func (suite *OrderRepositoryTestSuite) TestCreateOrderSuccessfully() {
	userID := "userID"
	orderLines := []*model.OrderLine{
		{
			ProductID: "productID",
			Quantity:  2,
		},
	}

	suite.mockDB.On("WithTransaction", mock.Anything).Return(nil).Times(1)

	order, err := suite.repo.CreateOrder(context.Background(), userID, orderLines)
	suite.NotNil(order)
	suite.Nil(err)
}

func (suite *OrderRepositoryTestSuite) TestCreateOrderFail() {
	userID := "userID"
	orderLines := []*model.OrderLine{
		{
			ProductID: "productID",
			Quantity:  2,
		},
	}

	suite.mockDB.On("WithTransaction", mock.Anything).Return(errors.New("error")).Times(1)

	order, err := suite.repo.CreateOrder(context.Background(), userID, orderLines)
	suite.Nil(order)
	suite.NotNil(err)
}

// UpdateOrder
// =================================================================

func (suite *OrderRepositoryTestSuite) TestUpdateOrderSuccessfully() {
	order := &model.Order{
		ID:   "orderId1",
		Code: "order",
	}
	suite.mockDB.On("Update", mock.Anything, order).
		Return(nil).Times(1)

	err := suite.repo.UpdateOrder(context.Background(), order)
	suite.Nil(err)
}

func (suite *OrderRepositoryTestSuite) TestUpdateOrderFail() {
	order := &model.Order{
		ID:   "orderId1",
		Code: "order",
	}
	suite.mockDB.On("Update", mock.Anything, order).
		Return(errors.New("error")).Times(1)

	err := suite.repo.UpdateOrder(context.Background(), order)
	suite.NotNil(err)
}

// GetOrderByID
// =================================================================

func (suite *OrderRepositoryTestSuite) TestGetOrderByIDSuccessfully() {
	suite.mockDB.On("FindOne", mock.Anything, &model.Order{}, mock.Anything, mock.Anything).
		Return(nil).Times(1)

	order, err := suite.repo.GetOrderByID(context.Background(), "orderId1", true)
	suite.Nil(err)
	suite.NotNil(order)
}

func (suite *OrderRepositoryTestSuite) TestGetOrderByIDFail() {
	suite.mockDB.On("FindOne", mock.Anything, &model.Order{}, mock.Anything, mock.Anything).
		Return(errors.New("error")).Times(1)

	order, err := suite.repo.GetOrderByID(context.Background(), "orderId1", true)
	suite.NotNil(err)
	suite.Nil(order)
}

//// GetMyOrders
//// =================================================================

func (suite *OrderRepositoryTestSuite) TestListOrdersSuccessfully() {
	req := &dto.ListOrderReq{
		UserID:    "userId",
		Code:      "code",
		Status:    "new",
		Page:      2,
		Limit:     10,
		OrderBy:   "name",
		OrderDesc: true,
	}

	suite.mockDB.On("Count", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(nil).Times(1)

	suite.mockDB.On("Find", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(nil).Times(1)

	orders, pagination, err := suite.repo.GetMyOrders(context.Background(), req)
	suite.Nil(err)
	suite.Equal(0, len(orders))
	suite.NotNil(pagination)
}

func (suite *OrderRepositoryTestSuite) TestListOrdersCountFail() {
	suite.mockDB.On("Count", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(errors.New("error")).Times(1)

	orders, pagination, err := suite.repo.GetMyOrders(context.Background(), &dto.ListOrderReq{})
	suite.NotNil(err)
	suite.Nil(orders)
	suite.Nil(pagination)
}

func (suite *OrderRepositoryTestSuite) TestListOrdersFindFail() {
	suite.mockDB.On("Count", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(nil).Times(1)

	suite.mockDB.On("Find", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(errors.New("error")).Times(1)

	orders, pagination, err := suite.repo.GetMyOrders(context.Background(), &dto.ListOrderReq{})
	suite.NotNil(err)
	suite.Nil(orders)
	suite.Nil(pagination)
}
