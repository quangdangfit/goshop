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
	line1 := &model.CartLine{CartID: "cartId1", ProductID: "productID1", Quantity: 4}
	line2 := &model.CartLine{CartID: "cartId1", ProductID: "productID2", Quantity: 3}
	cart := &model.Cart{
		ID:     "cartId1",
		UserID: "userID",
		Lines:  []*model.CartLine{line1, line2},
	}
	suite.mockDB.On("Update", mock.Anything, line1).Return(nil).Times(1)
	suite.mockDB.On("Update", mock.Anything, line2).Return(nil).Times(1)

	err := suite.repo.Update(context.Background(), cart)
	suite.Nil(err)
}

func (suite *CartRepositoryTestSuite) TestUpdateCartFail() {
	line1 := &model.CartLine{CartID: "cartId1", ProductID: "productID1", Quantity: 4}
	cart := &model.Cart{
		ID:     "cartId1",
		UserID: "userID",
		Lines:  []*model.CartLine{line1},
	}
	suite.mockDB.On("Update", mock.Anything, line1).Return(errors.New("error")).Times(1)

	err := suite.repo.Update(context.Background(), cart)
	suite.NotNil(err)
}

func (suite *CartRepositoryTestSuite) TestUpdateCartEmpty() {
	cart := &model.Cart{ID: "cartId1", UserID: "userID", Lines: []*model.CartLine{}}

	err := suite.repo.Update(context.Background(), cart)
	suite.Nil(err)
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

// DeleteLine
// =================================================================

func (suite *CartRepositoryTestSuite) TestDeleteLineSuccessfully() {
	suite.mockDB.On("Delete", mock.Anything, &model.CartLine{}, mock.Anything).
		Return(nil).Times(1)

	err := suite.repo.DeleteLine(context.Background(), "cartID", "productID")
	suite.Nil(err)
}

func (suite *CartRepositoryTestSuite) TestDeleteLineFail() {
	suite.mockDB.On("Delete", mock.Anything, &model.CartLine{}, mock.Anything).
		Return(errors.New("error")).Times(1)

	err := suite.repo.DeleteLine(context.Background(), "cartID", "productID")
	suite.NotNil(err)
}

// ClearCart
// =================================================================

func (suite *CartRepositoryTestSuite) TestClearCartSuccessfully() {
	suite.mockDB.On("FindOne", mock.Anything, &model.Cart{}, mock.Anything, mock.Anything).
		Return(nil).Times(1)
	suite.mockDB.On("Delete", mock.Anything, &model.CartLine{}, mock.Anything).
		Return(nil).Times(1)

	err := suite.repo.ClearCart(context.Background(), "userID")
	suite.Nil(err)
}

func (suite *CartRepositoryTestSuite) TestClearCartGetCartFail() {
	// When GetCartByUserID fails, ClearCart returns nil (cart doesn't exist)
	suite.mockDB.On("FindOne", mock.Anything, &model.Cart{}, mock.Anything, mock.Anything).
		Return(errors.New("not found")).Times(1)

	err := suite.repo.ClearCart(context.Background(), "userID")
	suite.Nil(err)
}

func (suite *CartRepositoryTestSuite) TestClearCartDeleteFail() {
	suite.mockDB.On("FindOne", mock.Anything, &model.Cart{}, mock.Anything, mock.Anything).
		Return(nil).Times(1)
	suite.mockDB.On("Delete", mock.Anything, &model.CartLine{}, mock.Anything).
		Return(errors.New("delete error")).Times(1)

	err := suite.repo.ClearCart(context.Background(), "userID")
	suite.NotNil(err)
}
