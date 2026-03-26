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
	"goshop/pkg/dbs/mocks"
)

type UserRepositoryTestSuite struct {
	suite.Suite
	mockDB *mocks.Database
	repo   UserRepository
}

func (suite *UserRepositoryTestSuite) SetupTest() {
	logger.Initialize(config.ProductionEnv)
	suite.mockDB = mocks.NewDatabase(suite.T())
	suite.repo = NewUserRepository(suite.mockDB)
}

func TestUserRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(UserRepositoryTestSuite))
}

func (suite *UserRepositoryTestSuite) TestCreate() {
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
			user := &model.User{Email: "test@test.com", Password: "test123456"}
			err := suite.repo.Create(context.Background(), user)
			if tc.wantErr {
				suite.NotNil(err)
			} else {
				suite.Nil(err)
			}
		})
	}
}

func (suite *UserRepositoryTestSuite) TestUpdate() {
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
			user := &model.User{ID: "userId1", Email: "test@test.com", Password: "test123456"}
			err := suite.repo.Update(context.Background(), user)
			if tc.wantErr {
				suite.NotNil(err)
			} else {
				suite.Nil(err)
			}
		})
	}
}

func (suite *UserRepositoryTestSuite) TestGetUserByID() {
	tests := []struct {
		name    string
		setup   func()
		wantErr bool
	}{
		{
			name: "Success",
			setup: func() {
				suite.mockDB.On("FindById", mock.Anything, "userId1", &model.User{}).Return(nil).Times(1)
			},
		},
		{
			name: "Not found",
			setup: func() {
				suite.mockDB.On("FindById", mock.Anything, "userId1", &model.User{}).Return(errors.New("error")).Times(1)
			},
			wantErr: true,
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.setup()
			user, err := suite.repo.GetUserByID(context.Background(), "userId1")
			if tc.wantErr {
				suite.NotNil(err)
				suite.Nil(user)
			} else {
				suite.Nil(err)
				suite.NotNil(user)
			}
		})
	}
}

func (suite *UserRepositoryTestSuite) TestGetUserByEmail() {
	tests := []struct {
		name    string
		setup   func()
		wantErr bool
	}{
		{
			name: "Success",
			setup: func() {
				suite.mockDB.On("FindOne", mock.Anything, &model.User{}, mock.Anything).Return(nil).Times(1)
			},
		},
		{
			name: "Not found",
			setup: func() {
				suite.mockDB.On("FindOne", mock.Anything, &model.User{}, mock.Anything).Return(errors.New("error")).Times(1)
			},
			wantErr: true,
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.setup()
			user, err := suite.repo.GetUserByEmail(context.Background(), "email@test.com")
			if tc.wantErr {
				suite.NotNil(err)
				suite.Nil(user)
			} else {
				suite.Nil(err)
				suite.NotNil(user)
			}
		})
	}
}
