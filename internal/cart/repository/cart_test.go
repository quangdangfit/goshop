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
	mockDB *mocks.IDatabase
	repo   ICartRepository
}

func (suite *CartRepositoryTestSuite) SetupTest() {
	logger.Initialize(config.ProductionEnv)

	suite.mockDB = mocks.NewIDatabase(suite.T())
	suite.repo = NewCartRepository(suite.mockDB)
}

func TestCartRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(CartRepositoryTestSuite))
}

// Create
// =================================================================

func (suite *CartRepositoryTestSuite) TestCreateCartSuccessfully() {
	cart := &model.Cart{
		UserID: "userID",
		Lines: []*model.CartLine{
			{
				ProductID: "productID1",
				Quantity:  4,
			},
			{
				ProductID: "productID2",
				Quantity:  3,
			},
		},
	}
	suite.mockDB.On("Create", mock.Anything, cart).
		Return(nil).Times(1)

	err := suite.repo.Create(context.Background(), cart)
	suite.Nil(err)
}

func (suite *CartRepositoryTestSuite) TestCreateCartFail() {
	cart := &model.Cart{
		UserID: "userID",
		Lines: []*model.CartLine{
			{
				ProductID: "productID1",
				Quantity:  4,
			},
			{
				ProductID: "productID2",
				Quantity:  3,
			},
		},
	}
	suite.mockDB.On("Create", mock.Anything, cart).
		Return(errors.New("error")).Times(1)

	err := suite.repo.Create(context.Background(), cart)
	suite.NotNil(err)
}

// Update
// =================================================================

func (suite *CartRepositoryTestSuite) TestUpdateCartSuccessfully() {
	cart := &model.Cart{
		ID:     "cartId1",
		UserID: "userID",
		Lines: []*model.CartLine{
			{
				ProductID: "productID1",
				Quantity:  4,
			},
			{
				ProductID: "productID2",
				Quantity:  3,
			},
		},
	}
	suite.mockDB.On("Update", mock.Anything, cart).
		Return(nil).Times(1)

	err := suite.repo.Update(context.Background(), cart)
	suite.Nil(err)
}

func (suite *CartRepositoryTestSuite) TestUpdateCartFail() {
	cart := &model.Cart{
		ID:     "cartId1",
		UserID: "userID",
		Lines: []*model.CartLine{
			{
				ProductID: "productID1",
				Quantity:  4,
			},
			{
				ProductID: "productID2",
				Quantity:  3,
			},
		},
	}
	suite.mockDB.On("Update", mock.Anything, cart).
		Return(errors.New("error")).Times(1)

	err := suite.repo.Update(context.Background(), cart)
	suite.NotNil(err)
}

// GetCartByUserID
// =================================================================

func (suite *CartRepositoryTestSuite) TestGetCartByUserIDSuccessfully() {
	suite.mockDB.On("FindOne", mock.Anything, &model.Cart{}, mock.Anything, mock.Anything).
		Return(nil).Times(1)

	cart, err := suite.repo.GetCartByUserID(context.Background(), "userId")
	suite.Nil(err)
	suite.NotNil(cart)
}

func (suite *CartRepositoryTestSuite) TestGetCartByUserIDFail() {
	suite.mockDB.On("FindOne", mock.Anything, &model.Cart{}, mock.Anything, mock.Anything).
		Return(errors.New("error")).Times(1)

	cart, err := suite.repo.GetCartByUserID(context.Background(), "userId")
	suite.NotNil(err)
	suite.Nil(cart)
}
