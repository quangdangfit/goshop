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

func (suite *WishlistServiceTestSuite) TestGetWishlist() {
	tests := []struct {
		name    string
		setup   func()
		wantErr bool
		wantLen int
	}{
		{
			name: "Success",
			setup: func() {
				suite.mockRepo.On("GetWishlist", mock.Anything, "u1").
					Return([]*model.Wishlist{
						{UserID: "u1", ProductID: "p1"},
						{UserID: "u1", ProductID: "p2"},
					}, nil).Times(1)
			},
			wantLen: 2,
		},
		{
			name: "DB error",
			setup: func() {
				suite.mockRepo.On("GetWishlist", mock.Anything, "u1").
					Return(nil, errors.New("db error")).Times(1)
			},
			wantErr: true,
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.setup()
			items, err := suite.service.GetWishlist(context.Background(), "u1")
			if tc.wantErr {
				suite.NotNil(err)
				suite.Nil(items)
			} else {
				suite.Nil(err)
				suite.Equal(tc.wantLen, len(items))
			}
		})
	}
}

func (suite *WishlistServiceTestSuite) TestAddProduct() {
	tests := []struct {
		name    string
		setup   func()
		wantErr bool
	}{
		{
			name: "Success",
			setup: func() {
				suite.mockRepo.On("Add", mock.Anything, "u1", "p1").Return(nil).Times(1)
			},
		},
		{
			name: "Already exists",
			setup: func() {
				suite.mockRepo.On("Add", mock.Anything, "u1", "p1").Return(errors.New("already exists")).Times(1)
			},
			wantErr: true,
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.setup()
			req := &dto.AddToWishlistReq{ProductID: "p1"}
			err := suite.service.AddProduct(context.Background(), "u1", req)
			if tc.wantErr {
				suite.NotNil(err)
			} else {
				suite.Nil(err)
			}
		})
	}
}

func (suite *WishlistServiceTestSuite) TestRemoveProduct() {
	tests := []struct {
		name    string
		setup   func()
		wantErr bool
	}{
		{
			name: "Success",
			setup: func() {
				suite.mockRepo.On("Remove", mock.Anything, "u1", "p1").Return(nil).Times(1)
			},
		},
		{
			name: "Not found",
			setup: func() {
				suite.mockRepo.On("Remove", mock.Anything, "u1", "p1").Return(errors.New("not found")).Times(1)
			},
			wantErr: true,
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.setup()
			err := suite.service.RemoveProduct(context.Background(), "u1", "p1")
			if tc.wantErr {
				suite.NotNil(err)
			} else {
				suite.Nil(err)
			}
		})
	}
}
