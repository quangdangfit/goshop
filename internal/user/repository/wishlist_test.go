package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/quangdangfit/gocommon/logger"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"goshop/internal/user/model"
	"goshop/pkg/config"
	dbsMocks "goshop/pkg/dbs/mocks"
)

type WishlistRepositoryTestSuite struct {
	suite.Suite
	mockDB *dbsMocks.Database
	repo   WishlistRepository
}

func (suite *WishlistRepositoryTestSuite) SetupTest() {
	logger.Initialize(config.ProductionEnv)
	suite.mockDB = dbsMocks.NewDatabase(suite.T())
	suite.repo = NewWishlistRepository(suite.mockDB)
}

func TestWishlistRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(WishlistRepositoryTestSuite))
}

// GetWishlist
// =================================================================

func (suite *WishlistRepositoryTestSuite) TestGetWishlistSuccess() {
	suite.mockDB.On("Find", mock.Anything, mock.Anything, mock.Anything).Return(nil).Times(1)

	items, err := suite.repo.GetWishlist(context.Background(), "u1")
	suite.Nil(err)
	suite.Equal(0, len(items))
}

func (suite *WishlistRepositoryTestSuite) TestGetWishlistFail() {
	suite.mockDB.On("Find", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("db error")).Times(1)

	items, err := suite.repo.GetWishlist(context.Background(), "u1")
	suite.NotNil(err)
	suite.Nil(items)
}

// Add
// =================================================================

func (suite *WishlistRepositoryTestSuite) TestAddSuccess() {
	suite.mockDB.On("Create", mock.Anything, &model.Wishlist{UserID: "u1", ProductID: "p1"}).Return(nil).Times(1)

	err := suite.repo.Add(context.Background(), "u1", "p1")
	suite.Nil(err)
}

func (suite *WishlistRepositoryTestSuite) TestAddFail() {
	suite.mockDB.On("Create", mock.Anything, &model.Wishlist{UserID: "u1", ProductID: "p1"}).Return(errors.New("duplicate")).Times(1)

	err := suite.repo.Add(context.Background(), "u1", "p1")
	suite.NotNil(err)
}

// Remove
// =================================================================

func (suite *WishlistRepositoryTestSuite) TestRemoveSuccess() {
	suite.mockDB.On("Delete", mock.Anything, mock.Anything, mock.Anything).Return(nil).Times(1)

	err := suite.repo.Remove(context.Background(), "u1", "p1")
	suite.Nil(err)
}

func (suite *WishlistRepositoryTestSuite) TestRemoveFail() {
	suite.mockDB.On("Delete", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("not found")).Times(1)

	err := suite.repo.Remove(context.Background(), "u1", "p1")
	suite.NotNil(err)
}
