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

func (suite *OrderRepositoryTestSuite) TestCreateOrder() {
	tests := []struct {
		name           string
		setup          func()
		wantErr        bool
		wantTotalPrice float64
	}{
		{
			name: "Success",
			setup: func() {
				suite.mockDB.On("WithTransaction", mock.Anything).Return(nil).Times(1)
			},
		},
		{
			name: "Executes transaction",
			setup: func() {
				suite.mockDB.On("WithTransaction", mock.Anything).
					Run(func(args mock.Arguments) {
						fn := args.Get(0).(func() error)
						suite.mockDB.On("Create", mock.Anything, mock.Anything).Return(nil).Once()
						suite.mockDB.On("CreateInBatches", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
						_ = fn()
					}).
					Return(nil).Times(1)
			},
			wantTotalPrice: 30.0,
		},
		{
			name: "Transaction create fails",
			setup: func() {
				suite.mockDB.On("WithTransaction", mock.Anything).
					Run(func(args mock.Arguments) {
						fn := args.Get(0).(func() error)
						suite.mockDB.On("Create", mock.Anything, mock.Anything).Return(errors.New("db error")).Once()
						_ = fn()
					}).
					Return(errors.New("db error")).Times(1)
			},
			wantErr: true,
		},
		{
			name: "Transaction batch fails",
			setup: func() {
				suite.mockDB.On("WithTransaction", mock.Anything).
					Run(func(args mock.Arguments) {
						fn := args.Get(0).(func() error)
						suite.mockDB.On("Create", mock.Anything, mock.Anything).Return(nil).Once()
						suite.mockDB.On("CreateInBatches", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("batch error")).Once()
						_ = fn()
					}).
					Return(errors.New("batch error")).Times(1)
			},
			wantErr: true,
		},
		{
			name: "WithTransaction fails",
			setup: func() {
				suite.mockDB.On("WithTransaction", mock.Anything).Return(errors.New("error")).Times(1)
			},
			wantErr: true,
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.setup()
			orderLines := []*model.OrderLine{{ProductID: "p1", Quantity: 1, Price: 10.0}, {ProductID: "p2", Quantity: 2, Price: 20.0}}
			if tc.name == "Success" || tc.name == "WithTransaction fails" {
				orderLines = []*model.OrderLine{{ProductID: "productID", Quantity: 2}}
			}
			order, err := suite.repo.CreateOrder(context.Background(), "userID", orderLines, "", 0)
			if tc.wantErr {
				suite.Nil(order)
				suite.NotNil(err)
			} else {
				suite.NotNil(order)
				suite.Nil(err)
				if tc.wantTotalPrice > 0 {
					suite.Equal(tc.wantTotalPrice, order.TotalPrice)
				}
			}
		})
	}
}

func (suite *OrderRepositoryTestSuite) TestUpdateOrder() {
	tests := []struct {
		name    string
		setup   func()
		wantErr bool
	}{
		{
			name: "Success",
			setup: func() {
				suite.mockDB.On("Update", mock.Anything, mock.Anything).Return(nil).Times(1)
			},
		},
		{
			name: "DB error",
			setup: func() {
				suite.mockDB.On("Update", mock.Anything, mock.Anything).Return(errors.New("error")).Times(1)
			},
			wantErr: true,
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.setup()
			order := &model.Order{ID: "orderId1", Code: "order"}
			err := suite.repo.UpdateOrder(context.Background(), order)
			if tc.wantErr {
				suite.NotNil(err)
			} else {
				suite.Nil(err)
			}
		})
	}
}

func (suite *OrderRepositoryTestSuite) TestGetOrderByID() {
	tests := []struct {
		name    string
		preload bool
		setup   func()
		wantErr bool
	}{
		{
			name:    "Success with preload",
			preload: true,
			setup: func() {
				suite.mockDB.On("FindOne", mock.Anything, &model.Order{}, mock.Anything, mock.Anything).Return(nil).Times(1)
			},
		},
		{
			name:    "Fail",
			preload: true,
			setup: func() {
				suite.mockDB.On("FindOne", mock.Anything, &model.Order{}, mock.Anything, mock.Anything).Return(errors.New("error")).Times(1)
			},
			wantErr: true,
		},
		{
			name:    "No preload",
			preload: false,
			setup: func() {
				suite.mockDB.On("FindOne", mock.Anything, &model.Order{}, mock.Anything).Return(nil).Times(1)
			},
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.setup()
			order, err := suite.repo.GetOrderByID(context.Background(), "orderId1", tc.preload)
			if tc.wantErr {
				suite.NotNil(err)
				suite.Nil(order)
			} else {
				suite.Nil(err)
				suite.NotNil(order)
			}
		})
	}
}

func (suite *OrderRepositoryTestSuite) TestGetMyOrders() {
	tests := []struct {
		name    string
		req     *dto.ListOrderReq
		setup   func()
		wantErr bool
	}{
		{
			name: "Success",
			req: &dto.ListOrderReq{
				UserID: "userId", Code: "code", Status: "new",
				Page: 2, Limit: 10, OrderBy: "name", OrderDesc: true,
			},
			setup: func() {
				suite.mockDB.On("Count", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Times(1)
				suite.mockDB.On("Find", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Times(1)
			},
		},
		{
			name: "Count fail",
			req:  &dto.ListOrderReq{},
			setup: func() {
				suite.mockDB.On("Count", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("error")).Times(1)
			},
			wantErr: true,
		},
		{
			name: "Find fail",
			req:  &dto.ListOrderReq{},
			setup: func() {
				suite.mockDB.On("Count", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Times(1)
				suite.mockDB.On("Find", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("error")).Times(1)
			},
			wantErr: true,
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.setup()
			orders, pagination, err := suite.repo.GetMyOrders(context.Background(), tc.req)
			if tc.wantErr {
				suite.NotNil(err)
				suite.Nil(orders)
				suite.Nil(pagination)
			} else {
				suite.Nil(err)
				suite.Equal(0, len(orders))
				suite.NotNil(pagination)
			}
		})
	}
}
