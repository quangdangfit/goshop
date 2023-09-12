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
	mockDB *mocks.IDatabase
	repo   IUserRepository
}

func (suite *UserRepositoryTestSuite) SetupTest() {
	logger.Initialize(config.ProductionEnv)

	suite.mockDB = mocks.NewIDatabase(suite.T())
	suite.repo = NewUserRepository(suite.mockDB)
}

func TestUserRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(UserRepositoryTestSuite))
}

// Create
// =================================================================

func (suite *UserRepositoryTestSuite) TestCreateUserSuccessfully() {
	user := &model.User{
		Email:    "test@test.com",
		Password: "test123456",
	}
	suite.mockDB.On("Create", mock.Anything, user).
		Return(nil).Times(1)

	err := suite.repo.Create(context.Background(), user)
	suite.Nil(err)
}

func (suite *UserRepositoryTestSuite) TestCreateUserFail() {
	user := &model.User{
		Email:    "test@test.com",
		Password: "test123456",
	}
	suite.mockDB.On("Create", mock.Anything, user).
		Return(errors.New("error")).Times(1)

	err := suite.repo.Create(context.Background(), user)
	suite.NotNil(err)
}

// Update
// =================================================================

func (suite *UserRepositoryTestSuite) TestUpdateUserSuccessfully() {
	user := &model.User{
		ID:       "userId1",
		Email:    "test@test.com",
		Password: "test123456",
	}
	suite.mockDB.On("Update", mock.Anything, user).
		Return(nil).Times(1)

	err := suite.repo.Update(context.Background(), user)
	suite.Nil(err)
}

func (suite *UserRepositoryTestSuite) TestUpdateUserFail() {
	user := &model.User{
		ID:       "userId1",
		Email:    "test@test.com",
		Password: "test123456",
	}
	suite.mockDB.On("Update", mock.Anything, user).
		Return(errors.New("error")).Times(1)

	err := suite.repo.Update(context.Background(), user)
	suite.NotNil(err)
}

// GetUserByID
// =================================================================

func (suite *UserRepositoryTestSuite) TestFindByIdSuccessfully() {
	suite.mockDB.On("FindById", mock.Anything, "userId1", &model.User{}).
		Return(nil).Times(1)

	user, err := suite.repo.GetUserByID(context.Background(), "userId1")
	suite.Nil(err)
	suite.NotNil(user)
}

func (suite *UserRepositoryTestSuite) TestFindByIdFail() {
	suite.mockDB.On("FindById", mock.Anything, "userId1", &model.User{}).
		Return(errors.New("error")).Times(1)

	user, err := suite.repo.GetUserByID(context.Background(), "userId1")
	suite.NotNil(err)
	suite.Nil(user)
}

// GetUserByEmail
// =================================================================

func (suite *UserRepositoryTestSuite) TestGetUserByEmailSuccessfully() {
	suite.mockDB.On("FindOne", mock.Anything, &model.User{}, mock.AnythingOfType("dbs.optionFn")).
		Return(nil).Times(1)

	user, err := suite.repo.GetUserByEmail(context.Background(), "email@test.com")
	suite.Nil(err)
	suite.NotNil(user)
}

func (suite *UserRepositoryTestSuite) TestGetUserByEmailFail() {
	suite.mockDB.On("FindOne", mock.Anything, &model.User{}, mock.AnythingOfType("dbs.optionFn")).
		Return(errors.New("error")).Times(1)

	user, err := suite.repo.GetUserByEmail(context.Background(), "email@test.com")
	suite.NotNil(err)
	suite.Nil(user)
}
