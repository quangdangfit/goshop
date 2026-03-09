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
	mockDB *mocks.Database
	repo   OrderRepository
}

func (suite *OrderRepositoryTestSuite) SetupTest() {
	logger.Initialize(config.ProductionEnv)

	suite.mockDB = mocks.NewDatabase(suite.T())
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

func (suite *OrderRepositoryTestSuite) TestCreateOrderExecutesTransaction() {
	userID := "userID"
	orderLines := []*model.OrderLine{
		{ProductID: "p1", Quantity: 1, Price: 10.0},
		{ProductID: "p2", Quantity: 2, Price: 20.0},
	}

	// Use Run to actually invoke the transaction function
	suite.mockDB.On("WithTransaction", mock.Anything).
		Run(func(args mock.Arguments) {
			fn := args.Get(0).(func() error)
			suite.mockDB.On("Create", mock.Anything, mock.Anything).Return(nil).Once()
			suite.mockDB.On("CreateInBatches", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
			_ = fn()
		}).
		Return(nil).Times(1)

	order, err := suite.repo.CreateOrder(context.Background(), userID, orderLines)
	suite.NotNil(order)
	suite.Nil(err)
	suite.Equal(30.0, order.TotalPrice)
}

func (suite *OrderRepositoryTestSuite) TestCreateOrderTransactionCreateFails() {
	userID := "userID"
	orderLines := []*model.OrderLine{{ProductID: "p1", Quantity: 1}}

	suite.mockDB.On("WithTransaction", mock.Anything).
		Run(func(args mock.Arguments) {
			fn := args.Get(0).(func() error)
			suite.mockDB.On("Create", mock.Anything, mock.Anything).Return(errors.New("db error")).Once()
			_ = fn()
		}).
		Return(errors.New("db error")).Times(1)

	order, err := suite.repo.CreateOrder(context.Background(), userID, orderLines)
	suite.Nil(order)
	suite.NotNil(err)
}

func (suite *OrderRepositoryTestSuite) TestCreateOrderTransactionBatchFails() {
	userID := "userID"
	orderLines := []*model.OrderLine{{ProductID: "p1", Quantity: 1}}

	suite.mockDB.On("WithTransaction", mock.Anything).
		Run(func(args mock.Arguments) {
			fn := args.Get(0).(func() error)
			suite.mockDB.On("Create", mock.Anything, mock.Anything).Return(nil).Once()
			suite.mockDB.On("CreateInBatches", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("batch error")).Once()
			_ = fn()
		}).
		Return(errors.New("batch error")).Times(1)

	order, err := suite.repo.CreateOrder(context.Background(), userID, orderLines)
	suite.Nil(order)
	suite.NotNil(err)
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

func (suite *OrderRepositoryTestSuite) TestGetOrderByIDNoPreload() {
	suite.mockDB.On("FindOne", mock.Anything, &model.Order{}, mock.Anything).
		Return(nil).Times(1)

	order, err := suite.repo.GetOrderByID(context.Background(), "orderId1", false)
	suite.Nil(err)
	suite.NotNil(order)
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
