package service

import (
	"context"
	"errors"
	"testing"

	"github.com/quangdangfit/gocommon/logger"
	"github.com/quangdangfit/gocommon/validation"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"goshop/internal/user/dto"
	"goshop/internal/user/model"
	"goshop/internal/user/repository/mocks"
	"goshop/pkg/config"
	"goshop/pkg/utils"
)

type UserServiceTestSuite struct {
	suite.Suite
	mockRepo *mocks.IUserRepository
	service  IUserService
}

func (suite *UserServiceTestSuite) SetupTest() {
	logger.Initialize(config.ProductionEnv)

	validator := validation.New()
	suite.mockRepo = mocks.NewIUserRepository(suite.T())
	suite.service = NewUserService(validator, suite.mockRepo)
}

func TestUserServiceTestSuite(t *testing.T) {
	suite.Run(t, new(UserServiceTestSuite))
}

// Login
// =================================================================

func (suite *UserServiceTestSuite) TestLoginGetUserByEmailFail() {
	req := &dto.LoginReq{
		Email:    "test@test.com",
		Password: "test123456",
	}
	suite.mockRepo.On("GetUserByEmail", mock.Anything, req.Email).
		Return(nil, errors.New("error")).Times(1)

	user, accessToken, refreshToken, err := suite.service.Login(context.Background(), req)
	suite.Nil(user)
	suite.Empty(accessToken)
	suite.Empty(refreshToken)
	suite.NotNil(err)
}

func (suite *UserServiceTestSuite) TestLoginInvalidEmailFormat() {
	req := &dto.LoginReq{
		Email:    "email",
		Password: "test123456",
	}

	user, accessToken, refreshToken, err := suite.service.Login(context.Background(), req)
	suite.Nil(user)
	suite.Empty(accessToken)
	suite.Empty(refreshToken)
	suite.NotNil(err)
}

func (suite *UserServiceTestSuite) TestLoginWrongPassword() {
	req := &dto.LoginReq{
		Email:    "test@test.com",
		Password: "test123456",
	}

	suite.mockRepo.On("GetUserByEmail", mock.Anything, req.Email).
		Return(&model.User{
			Email:    "test@test.com",
			Password: "password",
		}, nil).Times(1)

	user, accessToken, refreshToken, err := suite.service.Login(context.Background(), req)
	suite.Nil(user)
	suite.Empty(accessToken)
	suite.Empty(refreshToken)
	suite.NotNil(err)
}

func (suite *UserServiceTestSuite) TestLoginSuccess() {
	req := &dto.LoginReq{
		Email:    "test@test.com",
		Password: "test123456",
	}
	suite.mockRepo.On("GetUserByEmail", mock.Anything, req.Email).
		Return(
			&model.User{
				Email:    "test@test.com",
				Password: utils.HashAndSalt([]byte("test123456")),
			},
			nil,
		).Times(1)

	user, accessToken, refreshToken, err := suite.service.Login(context.Background(), req)
	suite.NotNil(user)
	suite.Equal(req.Email, user.Email)
	suite.NotEmpty(accessToken)
	suite.NotEmpty(refreshToken)
	suite.Nil(err)
}

// Register
// =================================================================

func (suite *UserServiceTestSuite) TestRegisterSuccess() {
	req := &dto.RegisterReq{
		Email:    "test@test.com",
		Password: "test123456",
	}
	suite.mockRepo.On("Create", mock.Anything, mock.Anything).
		Return(nil).Times(1)

	user, err := suite.service.Register(context.Background(), req)
	suite.NotNil(user)
	suite.Nil(err)
}

func (suite *UserServiceTestSuite) TestRegisterCreateUserFail() {
	req := &dto.RegisterReq{
		Email:    "test@test.com",
		Password: "test123456",
	}
	suite.mockRepo.On("Create", mock.Anything, mock.Anything).
		Return(errors.New("error")).Times(1)

	user, err := suite.service.Register(context.Background(), req)
	suite.Nil(user)
	suite.NotNil(err)
}

func (suite *UserServiceTestSuite) TestRegisterInvalidEmailFormat() {
	req := &dto.RegisterReq{
		Email:    "email",
		Password: "test123456",
	}
	user, err := suite.service.Register(context.Background(), req)
	suite.Nil(user)
	suite.NotNil(err)
}

// GetUserByID
// =================================================================

func (suite *UserServiceTestSuite) TestGetUserByIDSuccess() {
	userID := "userID"

	suite.mockRepo.On("GetUserByID", mock.Anything, userID).
		Return(
			&model.User{
				ID:    userID,
				Email: "test@test.com",
			},
			nil,
		).Times(1)

	user, err := suite.service.GetUserByID(context.Background(), userID)
	suite.NotNil(user)
	suite.Equal(userID, user.ID)
	suite.Equal("test@test.com", user.Email)
	suite.Nil(err)
}

func (suite *UserServiceTestSuite) TestGetUserByIDFail() {
	userID := "userID"
	suite.mockRepo.On("GetUserByID", mock.Anything, userID).
		Return(nil, errors.New("error")).Times(1)

	user, err := suite.service.GetUserByID(context.Background(), userID)
	suite.Nil(user)
	suite.NotNil(err)
}

// RefreshToken
// =================================================================

func (suite *UserServiceTestSuite) TestRefreshTokenSuccess() {
	userID := "userID"
	suite.mockRepo.On("GetUserByID", mock.Anything, userID).
		Return(
			&model.User{
				ID:    userID,
				Email: "test@test.com",
			}, nil,
		).Times(1)

	refreshToken, err := suite.service.RefreshToken(context.Background(), userID)
	suite.NotEmpty(refreshToken)
	suite.Nil(err)
}

func (suite *UserServiceTestSuite) TestRefreshTokenGetUserByIDFail() {
	userID := "userID"
	suite.mockRepo.On("GetUserByID", mock.Anything, userID).
		Return(nil, errors.New("error")).Times(1)

	refreshToken, err := suite.service.RefreshToken(context.Background(), userID)
	suite.Empty(refreshToken)
	suite.NotNil(err)
}

// ChangePassword
// =================================================================

func (suite *UserServiceTestSuite) TestChangePasswordSuccess() {
	userID := "userID"
	req := &dto.ChangePasswordReq{
		Password:    "password",
		NewPassword: "newPassword",
	}

	suite.mockRepo.On("GetUserByID", mock.Anything, userID).
		Return(
			&model.User{
				ID:       userID,
				Email:    "test@test.com",
				Password: utils.HashAndSalt([]byte("password")),
			}, nil,
		).Times(1)
	suite.mockRepo.On("Update", mock.Anything, mock.Anything).
		Return(nil).Times(1)

	err := suite.service.ChangePassword(context.Background(), userID, req)
	suite.Nil(err)
}

func (suite *UserServiceTestSuite) TestChangePasswordGetUserByIDFail() {
	userID := "userID"
	req := &dto.ChangePasswordReq{
		Password:    "password",
		NewPassword: "newPassword",
	}

	suite.mockRepo.On("GetUserByID", mock.Anything, userID).
		Return(nil, errors.New("error")).Times(1)

	err := suite.service.ChangePassword(context.Background(), userID, req)
	suite.NotNil(err)
}

func (suite *UserServiceTestSuite) TestChangePasswordMissRequiredField() {
	userID := "userID"
	req := &dto.ChangePasswordReq{
		Password:    "password",
		NewPassword: "",
	}

	err := suite.service.ChangePassword(context.Background(), userID, req)
	suite.NotNil(err)
}

func (suite *UserServiceTestSuite) TestChangePasswordWrongCurrentPassword() {
	userID := "userID"
	req := &dto.ChangePasswordReq{
		Password:    "password1",
		NewPassword: "newPassword",
	}

	suite.mockRepo.On("GetUserByID", mock.Anything, userID).
		Return(
			&model.User{
				ID:       userID,
				Email:    "test@test.com",
				Password: utils.HashAndSalt([]byte("password")),
			}, nil,
		).Times(1)

	err := suite.service.ChangePassword(context.Background(), userID, req)
	suite.NotNil(err)
}

func (suite *UserServiceTestSuite) TestChangePasswordUpdateUserFail() {
	userID := "userID"
	req := &dto.ChangePasswordReq{
		Password:    "password",
		NewPassword: "newPassword",
	}

	suite.mockRepo.On("GetUserByID", mock.Anything, userID).
		Return(
			&model.User{
				ID:       userID,
				Email:    "test@test.com",
				Password: utils.HashAndSalt([]byte("password")),
			}, nil,
		).Times(1)
	suite.mockRepo.On("Update", mock.Anything, mock.Anything).
		Return(errors.New("error")).Times(1)

	err := suite.service.ChangePassword(context.Background(), userID, req)
	suite.NotNil(err)
}
