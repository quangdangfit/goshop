package services

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/quangdangfit/gocommon/logger"
	"github.com/stretchr/testify/suite"

	"goshop/app/models"
	"goshop/app/serializers"
	"goshop/config"
	"goshop/mocks"
	"goshop/pkg/utils"
)

type UserServiceTestSuite struct {
	suite.Suite
	mockRepo *mocks.MockIUserRepository
	service  IUserService
}

func (suite *UserServiceTestSuite) SetupTest() {
	logger.Initialize(config.TestEnv)

	mockCtrl := gomock.NewController(suite.T())
	defer mockCtrl.Finish()
	suite.mockRepo = mocks.NewMockIUserRepository(mockCtrl)
	suite.service = NewUserService(suite.mockRepo)
}

func TestUserServiceTestSuite(t *testing.T) {
	suite.Run(t, new(UserServiceTestSuite))
}

// Login
// =================================================================

func (suite *UserServiceTestSuite) TestLoginGetUserByEmailFail() {
	req := &serializers.LoginReq{
		Email:    "test@test.com",
		Password: "test123456",
	}
	suite.mockRepo.EXPECT().GetUserByEmail(gomock.Any(), req.Email).Return(nil, errors.New("error")).Times(1)

	user, accessToken, refreshToken, err := suite.service.Login(context.Background(), req)
	suite.Nil(user)
	suite.Empty(accessToken)
	suite.Empty(refreshToken)
	suite.NotNil(err)
}

func (suite *UserServiceTestSuite) TestLoginSuccess() {
	req := &serializers.LoginReq{
		Email:    "test@test.com",
		Password: "test123456",
	}
	suite.mockRepo.EXPECT().GetUserByEmail(gomock.Any(), req.Email).Return(&models.User{
		Email:    "test@test.com",
		Password: utils.HashAndSalt([]byte("test123456")),
	}, nil).Times(1)

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
	req := &serializers.RegisterReq{
		Email:    "test@test.com",
		Password: "test123456",
	}
	suite.mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil).Times(1)

	user, err := suite.service.Register(context.Background(), req)
	suite.NotNil(user)
	suite.Nil(err)
}

func (suite *UserServiceTestSuite) TestRegisterCreateUserFail() {
	req := &serializers.RegisterReq{
		Email:    "test@test.com",
		Password: "test123456",
	}
	suite.mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(errors.New("error")).Times(1)

	user, err := suite.service.Register(context.Background(), req)
	suite.Nil(user)
	suite.NotNil(err)
}

// GetUserByID
// =================================================================

func (suite *UserServiceTestSuite) TestGetUserByIDSuccess() {
	userID := "userID"
	suite.mockRepo.EXPECT().GetUserByID(gomock.Any(), userID).Return(&models.User{
		Base: models.Base{
			ID: userID,
		},
		Email: "test@test.com",
	}, nil).Times(1)

	user, err := suite.service.GetUserByID(context.Background(), userID)
	suite.NotNil(user)
	suite.Equal(userID, user.ID)
	suite.Equal("test@test.com", user.Email)
	suite.Nil(err)
}

func (suite *UserServiceTestSuite) TestGetUserByIDFail() {
	userID := "userID"
	suite.mockRepo.EXPECT().GetUserByID(gomock.Any(), userID).Return(nil, errors.New("error")).Times(1)

	user, err := suite.service.GetUserByID(context.Background(), userID)
	suite.Nil(user)
	suite.NotNil(err)
}

// RefreshToken
// =================================================================

func (suite *UserServiceTestSuite) TestRefreshTokenSuccess() {
	userID := "userID"
	suite.mockRepo.EXPECT().GetUserByID(gomock.Any(), userID).Return(&models.User{
		Base: models.Base{
			ID: userID,
		},
		Email: "test@test.com",
	}, nil).Times(1)

	refreshToken, err := suite.service.GetUserByID(context.Background(), userID)
	suite.NotEmpty(refreshToken)
	suite.Nil(err)
}

func (suite *UserServiceTestSuite) TestRefreshTokenGetUserByIDFail() {
	userID := "userID"
	suite.mockRepo.EXPECT().GetUserByID(gomock.Any(), userID).Return(nil, errors.New("error")).Times(1)

	refreshToken, err := suite.service.GetUserByID(context.Background(), userID)
	suite.Empty(refreshToken)
	suite.NotNil(err)
}

// ChangePassword
// =================================================================

func (suite *UserServiceTestSuite) TestChangePasswordSuccess() {
	userID := "userID"
	req := &serializers.ChangePasswordReq{
		Password:    "password",
		NewPassword: "newPassword",
	}

	suite.mockRepo.EXPECT().GetUserByID(gomock.Any(), userID).Return(&models.User{
		Base: models.Base{
			ID: userID,
		},
		Email: "test@test.com",
	}, nil).Times(1)

	suite.mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil).Times(1)

	err := suite.service.ChangePassword(context.Background(), userID, req)
	suite.Nil(err)
}

func (suite *UserServiceTestSuite) TestChangePasswordGetUserByIDFail() {
	userID := "userID"
	req := &serializers.ChangePasswordReq{
		Password:    "password",
		NewPassword: "newPassword",
	}

	suite.mockRepo.EXPECT().GetUserByID(gomock.Any(), userID).Return(nil, errors.New("error")).Times(1)

	err := suite.service.ChangePassword(context.Background(), userID, req)
	suite.NotNil(err)
}

func (suite *UserServiceTestSuite) TestChangePasswordUpdateUserFail() {
	userID := "userID"
	req := &serializers.ChangePasswordReq{
		Password:    "password",
		NewPassword: "newPassword",
	}

	suite.mockRepo.EXPECT().GetUserByID(gomock.Any(), userID).Return(&models.User{
		Base: models.Base{
			ID: userID,
		},
		Email: "test@test.com",
	}, nil).Times(1)

	suite.mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(errors.New("error")).Times(1)

	err := suite.service.ChangePassword(context.Background(), userID, req)
	suite.NotNil(err)
}
