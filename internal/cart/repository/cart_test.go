package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/quangdangfit/gocommon/logger"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"goshop/internal/cart/model"
	"goshop/pkg/config"
	"goshop/pkg/dbs/mocks"
)

type CartRepositoryTestSuite struct {
	suite.Suite
	mockDB *mocks.Database
	repo   CartRepository
}

func (suite *CartRepositoryTestSuite) SetupTest() {
	logger.Initialize(config.ProductionEnv)

	suite.mockDB = mocks.NewDatabase(suite.T())
	suite.repo = NewCartRepository(suite.mockDB)
}

func TestCartRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(CartRepositoryTestSuite))
}

func (suite *CartRepositoryTestSuite) TestCreate() {
	tests := []struct {
		name    string
		setup   func()
		wantErr bool
	}{
		{
			name: "Success",
			setup: func() {
				suite.mockDB.On("Create", mock.Anything, mock.Anything).Return(nil).Times(1)
			},
		},
		{
			name: "DB error",
			setup: func() {
				suite.mockDB.On("Create", mock.Anything, mock.Anything).Return(errors.New("error")).Times(1)
			},
			wantErr: true,
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.setup()
			cart := &model.Cart{
				UserID: "userID",
				Lines: []*model.CartLine{
					{ProductID: "productID1", Quantity: 4},
					{ProductID: "productID2", Quantity: 3},
				},
			}
			err := suite.repo.Create(context.Background(), cart)
			if tc.wantErr {
				suite.NotNil(err)
			} else {
				suite.Nil(err)
			}
		})
	}
}

func (suite *CartRepositoryTestSuite) TestUpdate() {
	tests := []struct {
		name    string
		setup   func()
		wantErr bool
		empty   bool
	}{
		{
			name: "Success",
			setup: func() {
				suite.mockDB.On("Update", mock.Anything, mock.Anything).Return(nil).Times(2)
			},
		},
		{
			name: "Fail",
			setup: func() {
				suite.mockDB.On("Update", mock.Anything, mock.Anything).Return(errors.New("error")).Times(1)
			},
			wantErr: true,
		},
		{
			name:  "Empty lines",
			setup: func() {},
			empty: true,
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.setup()
			var lines []*model.CartLine
			if !tc.empty {
				if tc.wantErr {
					lines = []*model.CartLine{{CartID: "cartId1", ProductID: "productID1", Quantity: 4}}
				} else {
					lines = []*model.CartLine{
						{CartID: "cartId1", ProductID: "productID1", Quantity: 4},
						{CartID: "cartId1", ProductID: "productID2", Quantity: 3},
					}
				}
			} else {
				lines = []*model.CartLine{}
			}
			cart := &model.Cart{ID: "cartId1", UserID: "userID", Lines: lines}
			err := suite.repo.Update(context.Background(), cart)
			if tc.wantErr {
				suite.NotNil(err)
			} else {
				suite.Nil(err)
			}
		})
	}
}

func (suite *CartRepositoryTestSuite) TestGetCartByUserID() {
	tests := []struct {
		name    string
		setup   func()
		wantErr bool
	}{
		{
			name: "Success",
			setup: func() {
				suite.mockDB.On("FindOne", mock.Anything, &model.Cart{}, mock.Anything, mock.Anything).Return(nil).Times(1)
			},
		},
		{
			name: "Not found",
			setup: func() {
				suite.mockDB.On("FindOne", mock.Anything, &model.Cart{}, mock.Anything, mock.Anything).Return(errors.New("error")).Times(1)
			},
			wantErr: true,
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.setup()
			cart, err := suite.repo.GetCartByUserID(context.Background(), "userId")
			if tc.wantErr {
				suite.NotNil(err)
				suite.Nil(cart)
			} else {
				suite.Nil(err)
				suite.NotNil(cart)
			}
		})
	}
}

func (suite *CartRepositoryTestSuite) TestDeleteLine() {
	tests := []struct {
		name    string
		setup   func()
		wantErr bool
	}{
		{
			name: "Success",
			setup: func() {
				suite.mockDB.On("Delete", mock.Anything, &model.CartLine{}, mock.Anything).Return(nil).Times(1)
			},
		},
		{
			name: "DB error",
			setup: func() {
				suite.mockDB.On("Delete", mock.Anything, &model.CartLine{}, mock.Anything).Return(errors.New("error")).Times(1)
			},
			wantErr: true,
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.setup()
			err := suite.repo.DeleteLine(context.Background(), "cartID", "productID")
			if tc.wantErr {
				suite.NotNil(err)
			} else {
				suite.Nil(err)
			}
		})
	}
}

func (suite *CartRepositoryTestSuite) TestClearCart() {
	tests := []struct {
		name    string
		setup   func()
		wantErr bool
	}{
		{
			name: "Success",
			setup: func() {
				suite.mockDB.On("FindOne", mock.Anything, &model.Cart{}, mock.Anything, mock.Anything).Return(nil).Times(1)
				suite.mockDB.On("Delete", mock.Anything, &model.CartLine{}, mock.Anything).Return(nil).Times(1)
			},
		},
		{
			name: "Cart not found",
			setup: func() {
				suite.mockDB.On("FindOne", mock.Anything, &model.Cart{}, mock.Anything, mock.Anything).Return(errors.New("not found")).Times(1)
			},
		},
		{
			name: "Delete fail",
			setup: func() {
				suite.mockDB.On("FindOne", mock.Anything, &model.Cart{}, mock.Anything, mock.Anything).Return(nil).Times(1)
				suite.mockDB.On("Delete", mock.Anything, &model.CartLine{}, mock.Anything).Return(errors.New("delete error")).Times(1)
			},
			wantErr: true,
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.setup()
			err := suite.repo.ClearCart(context.Background(), "userID")
			if tc.wantErr {
				suite.NotNil(err)
			} else {
				suite.Nil(err)
			}
		})
	}
}
