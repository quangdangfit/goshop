package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/quangdangfit/gocommon/logger"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"goshop/internal/order/model"
	"goshop/pkg/config"
	dbsMocks "goshop/pkg/dbs/mocks"
)

type UserRepositoryOrderTestSuite struct {
	suite.Suite
	mockDB *dbsMocks.Database
	repo   UserRepository
}

func (suite *UserRepositoryOrderTestSuite) SetupTest() {
	logger.Initialize(config.ProductionEnv)
	suite.mockDB = dbsMocks.NewDatabase(suite.T())
	suite.repo = NewUserRepository(suite.mockDB)
}

func TestUserRepositoryOrderTestSuite(t *testing.T) {
	suite.Run(t, new(UserRepositoryOrderTestSuite))
}

func (suite *UserRepositoryOrderTestSuite) TestGetUserByID() {
	tests := []struct {
		name    string
		id      string
		setup   func()
		wantErr bool
	}{
		{
			name: "Success",
			id:   "u1",
			setup: func() {
				suite.mockDB.On("FindById", mock.Anything, "u1", &model.User{}).Return(nil).Times(1)
			},
		},
		{
			name: "Not found",
			id:   "notfound",
			setup: func() {
				suite.mockDB.On("FindById", mock.Anything, "notfound", &model.User{}).Return(errors.New("not found")).Times(1)
			},
			wantErr: true,
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.setup()
			user, err := suite.repo.GetUserByID(context.Background(), tc.id)
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
