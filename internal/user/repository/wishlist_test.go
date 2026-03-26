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

func (suite *WishlistRepositoryTestSuite) TestGetWishlist() {
	tests := []struct {
		name    string
		setup   func()
		wantErr bool
	}{
		{
			name: "Success",
			setup: func() {
				suite.mockDB.On("Find", mock.Anything, mock.Anything, mock.Anything).Return(nil).Times(1)
			},
		},
		{
			name: "DB error",
			setup: func() {
				suite.mockDB.On("Find", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("db error")).Times(1)
			},
			wantErr: true,
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.setup()
			items, err := suite.repo.GetWishlist(context.Background(), "u1")
			if tc.wantErr {
				suite.NotNil(err)
				suite.Nil(items)
			} else {
				suite.Nil(err)
				suite.Equal(0, len(items))
			}
		})
	}
}

func (suite *WishlistRepositoryTestSuite) TestAdd() {
	tests := []struct {
		name    string
		setup   func()
		wantErr bool
	}{
		{
			name: "Success",
			setup: func() {
				suite.mockDB.On("Create", mock.Anything, &model.Wishlist{UserID: "u1", ProductID: "p1"}).Return(nil).Times(1)
			},
		},
		{
			name: "Duplicate",
			setup: func() {
				suite.mockDB.On("Create", mock.Anything, &model.Wishlist{UserID: "u1", ProductID: "p1"}).Return(errors.New("duplicate")).Times(1)
			},
			wantErr: true,
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.setup()
			err := suite.repo.Add(context.Background(), "u1", "p1")
			if tc.wantErr {
				suite.NotNil(err)
			} else {
				suite.Nil(err)
			}
		})
	}
}

func (suite *WishlistRepositoryTestSuite) TestRemove() {
	tests := []struct {
		name    string
		setup   func()
		wantErr bool
	}{
		{
			name: "Success",
			setup: func() {
				suite.mockDB.On("Delete", mock.Anything, mock.Anything, mock.Anything).Return(nil).Times(1)
			},
		},
		{
			name: "Not found",
			setup: func() {
				suite.mockDB.On("Delete", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("not found")).Times(1)
			},
			wantErr: true,
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.setup()
			err := suite.repo.Remove(context.Background(), "u1", "p1")
			if tc.wantErr {
				suite.NotNil(err)
			} else {
				suite.Nil(err)
			}
		})
	}
}
