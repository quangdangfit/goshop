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

func (suite *UserServiceTestSuite) TestChangePasswordUserNotFound() {
	req := &serializers.ChangePasswordReq{
		Password:    "password",
		NewPassword: "newpassword",
	}
	suite.mockRepo.EXPECT().GetUserByID(gomock.Any(), gomock.Any()).Return(nil, errors.New("record not found")).Times(1)

	err := suite.service.ChangePassword(context.TODO(), "notfoundid", req)
	suite.NotNil(err)
}

func (suite *UserServiceTestSuite) TestChangePasswordUpdateFail() {
	req := &serializers.ChangePasswordReq{
		Password:    "password",
		NewPassword: "newpassword",
	}
	suite.mockRepo.EXPECT().GetUserByID(gomock.Any(), gomock.Any()).Return(&models.User{}, nil).Times(1)
	suite.mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(errors.New("update fail")).Times(1)

	err := suite.service.ChangePassword(context.TODO(), "notfoundid", req)
	suite.NotNil(err)
}
