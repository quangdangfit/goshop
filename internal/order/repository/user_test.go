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

// GetUserByID
// =================================================================

func (suite *UserRepositoryOrderTestSuite) TestGetUserByIDSuccess() {
	suite.mockDB.On("FindById", mock.Anything, "u1", &model.User{}).Return(nil).Times(1)

	user, err := suite.repo.GetUserByID(context.Background(), "u1")
	suite.Nil(err)
	suite.NotNil(user)
}

func (suite *UserRepositoryOrderTestSuite) TestGetUserByIDFail() {
	suite.mockDB.On("FindById", mock.Anything, "notfound", &model.User{}).Return(errors.New("not found")).Times(1)

	user, err := suite.repo.GetUserByID(context.Background(), "notfound")
	suite.NotNil(err)
	suite.Nil(user)
}
