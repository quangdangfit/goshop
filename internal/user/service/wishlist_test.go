package service

import (
	"context"
	"errors"
	"testing"

	"github.com/quangdangfit/gocommon/logger"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"goshop/internal/user/dto"
	"goshop/internal/user/model"
	"goshop/internal/user/repository/mocks"
	"goshop/pkg/config"
)

type WishlistServiceTestSuite struct {
	suite.Suite
	mockRepo *mocks.WishlistRepository
	service  WishlistService
}

func (suite *WishlistServiceTestSuite) SetupTest() {
	logger.Initialize(config.ProductionEnv)
	suite.mockRepo = mocks.NewWishlistRepository(suite.T())
	suite.service = NewWishlistService(suite.mockRepo)
}

func TestWishlistServiceTestSuite(t *testing.T) {
	suite.Run(t, new(WishlistServiceTestSuite))
}

// GetWishlist
// =================================================================================================

func (suite *WishlistServiceTestSuite) TestGetWishlistSuccess() {
	suite.mockRepo.On("GetWishlist", mock.Anything, "u1").
		Return([]*model.Wishlist{
			{UserID: "u1", ProductID: "p1"},
			{UserID: "u1", ProductID: "p2"},
		}, nil).Times(1)

	items, err := suite.service.GetWishlist(context.Background(), "u1")
	suite.Nil(err)
	suite.Equal(2, len(items))
}

func (suite *WishlistServiceTestSuite) TestGetWishlistFail() {
	suite.mockRepo.On("GetWishlist", mock.Anything, "u1").
		Return(nil, errors.New("db error")).Times(1)

	items, err := suite.service.GetWishlist(context.Background(), "u1")
	suite.NotNil(err)
	suite.Nil(items)
}

// AddProduct
// =================================================================================================

func (suite *WishlistServiceTestSuite) TestAddProductSuccess() {
	req := &dto.AddToWishlistReq{ProductID: "p1"}
	suite.mockRepo.On("Add", mock.Anything, "u1", "p1").Return(nil).Times(1)

	err := suite.service.AddProduct(context.Background(), "u1", req)
	suite.Nil(err)
}

func (suite *WishlistServiceTestSuite) TestAddProductFail() {
	req := &dto.AddToWishlistReq{ProductID: "p1"}
	suite.mockRepo.On("Add", mock.Anything, "u1", "p1").Return(errors.New("already exists")).Times(1)

	err := suite.service.AddProduct(context.Background(), "u1", req)
	suite.NotNil(err)
}

// RemoveProduct
// =================================================================================================

func (suite *WishlistServiceTestSuite) TestRemoveProductSuccess() {
	suite.mockRepo.On("Remove", mock.Anything, "u1", "p1").Return(nil).Times(1)

	err := suite.service.RemoveProduct(context.Background(), "u1", "p1")
	suite.Nil(err)
}

func (suite *WishlistServiceTestSuite) TestRemoveProductFail() {
	suite.mockRepo.On("Remove", mock.Anything, "u1", "p1").Return(errors.New("not found")).Times(1)

	err := suite.service.RemoveProduct(context.Background(), "u1", "p1")
	suite.NotNil(err)
}
