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
	mockRepo *mocks.UserRepository
	service  UserService
}

func (suite *UserServiceTestSuite) SetupTest() {
	logger.Initialize(config.ProductionEnv)
	validator := validation.New()
	suite.mockRepo = mocks.NewUserRepository(suite.T())
	suite.service = NewUserService(validator, suite.mockRepo)
}

func TestUserServiceTestSuite(t *testing.T) {
	suite.Run(t, new(UserServiceTestSuite))
}

func (suite *UserServiceTestSuite) TestLogin() {
	tests := []struct {
		name      string
		req       *dto.LoginReq
		setup     func()
		wantUser  bool
		wantToken bool
		wantErr   bool
	}{
		{
			name: "Success",
			req:  &dto.LoginReq{Email: "test@test.com", Password: "test123456"},
			setup: func() {
				suite.mockRepo.On("GetUserByEmail", mock.Anything, "test@test.com").
					Return(&model.User{
						Email:    "test@test.com",
						Password: utils.HashAndSalt([]byte("test123456")),
					}, nil).Times(1)
			},
			wantUser:  true,
			wantToken: true,
		},
		{
			name: "GetUserByEmail fail",
			req:  &dto.LoginReq{Email: "test@test.com", Password: "test123456"},
			setup: func() {
				suite.mockRepo.On("GetUserByEmail", mock.Anything, "test@test.com").
					Return(nil, errors.New("error")).Times(1)
			},
			wantErr: true,
		},
		{
			name:    "Invalid email format",
			req:     &dto.LoginReq{Email: "email", Password: "test123456"},
			setup:   func() {},
			wantErr: true,
		},
		{
			name: "Wrong password",
			req:  &dto.LoginReq{Email: "test@test.com", Password: "test123456"},
			setup: func() {
				suite.mockRepo.On("GetUserByEmail", mock.Anything, "test@test.com").
					Return(&model.User{Email: "test@test.com", Password: "password"}, nil).Times(1)
			},
			wantErr: true,
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.setup()
			user, accessToken, refreshToken, err := suite.service.Login(context.Background(), tc.req)
			if tc.wantErr {
				suite.Nil(user)
				suite.Empty(accessToken)
				suite.Empty(refreshToken)
				suite.NotNil(err)
			} else {
				suite.NotNil(user)
				suite.Equal(tc.req.Email, user.Email)
				if tc.wantToken {
					suite.NotEmpty(accessToken)
					suite.NotEmpty(refreshToken)
				}
				suite.Nil(err)
			}
		})
	}
}

func (suite *UserServiceTestSuite) TestRegister() {
	tests := []struct {
		name    string
		req     *dto.RegisterReq
		setup   func()
		wantErr bool
	}{
		{
			name: "Success",
			req:  &dto.RegisterReq{Email: "test@test.com", Password: "test123456"},
			setup: func() {
				suite.mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil).Times(1)
			},
		},
		{
			name: "Create fail",
			req:  &dto.RegisterReq{Email: "test@test.com", Password: "test123456"},
			setup: func() {
				suite.mockRepo.On("Create", mock.Anything, mock.Anything).Return(errors.New("error")).Times(1)
			},
			wantErr: true,
		},
		{
			name:    "Invalid email format",
			req:     &dto.RegisterReq{Email: "email", Password: "test123456"},
			setup:   func() {},
			wantErr: true,
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.setup()
			user, err := suite.service.Register(context.Background(), tc.req)
			if tc.wantErr {
				suite.Nil(user)
				suite.NotNil(err)
			} else {
				suite.NotNil(user)
				suite.Nil(err)
			}
		})
	}
}

func (suite *UserServiceTestSuite) TestGetUserByID() {
	tests := []struct {
		name    string
		setup   func()
		wantErr bool
	}{
		{
			name: "Success",
			setup: func() {
				suite.mockRepo.On("GetUserByID", mock.Anything, "userID").
					Return(&model.User{ID: "userID", Email: "test@test.com"}, nil).Times(1)
			},
		},
		{
			name: "Not found",
			setup: func() {
				suite.mockRepo.On("GetUserByID", mock.Anything, "userID").
					Return(nil, errors.New("error")).Times(1)
			},
			wantErr: true,
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.setup()
			user, err := suite.service.GetUserByID(context.Background(), "userID")
			if tc.wantErr {
				suite.Nil(user)
				suite.NotNil(err)
			} else {
				suite.NotNil(user)
				suite.Equal("userID", user.ID)
				suite.Equal("test@test.com", user.Email)
				suite.Nil(err)
			}
		})
	}
}

func (suite *UserServiceTestSuite) TestRefreshToken() {
	tests := []struct {
		name    string
		setup   func()
		wantErr bool
	}{
		{
			name: "Success",
			setup: func() {
				suite.mockRepo.On("GetUserByID", mock.Anything, "userID").
					Return(&model.User{ID: "userID", Email: "test@test.com"}, nil).Times(1)
			},
		},
		{
			name: "GetUserByID fail",
			setup: func() {
				suite.mockRepo.On("GetUserByID", mock.Anything, "userID").
					Return(nil, errors.New("error")).Times(1)
			},
			wantErr: true,
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.setup()
			refreshToken, err := suite.service.RefreshToken(context.Background(), "userID")
			if tc.wantErr {
				suite.Empty(refreshToken)
				suite.NotNil(err)
			} else {
				suite.NotEmpty(refreshToken)
				suite.Nil(err)
			}
		})
	}
}

func (suite *UserServiceTestSuite) TestChangePassword() {
	tests := []struct {
		name    string
		req     *dto.ChangePasswordReq
		setup   func()
		wantErr bool
	}{
		{
			name: "Success",
			req:  &dto.ChangePasswordReq{Password: "password", NewPassword: "newPassword"},
			setup: func() {
				suite.mockRepo.On("GetUserByID", mock.Anything, "userID").
					Return(&model.User{
						ID: "userID", Email: "test@test.com",
						Password: utils.HashAndSalt([]byte("password")),
					}, nil).Times(1)
				suite.mockRepo.On("Update", mock.Anything, mock.Anything).Return(nil).Times(1)
			},
		},
		{
			name: "GetUserByID fail",
			req:  &dto.ChangePasswordReq{Password: "password", NewPassword: "newPassword"},
			setup: func() {
				suite.mockRepo.On("GetUserByID", mock.Anything, "userID").
					Return(nil, errors.New("error")).Times(1)
			},
			wantErr: true,
		},
		{
			name:    "Missing required field",
			req:     &dto.ChangePasswordReq{Password: "password", NewPassword: ""},
			setup:   func() {},
			wantErr: true,
		},
		{
			name: "Wrong current password",
			req:  &dto.ChangePasswordReq{Password: "password1", NewPassword: "newPassword"},
			setup: func() {
				suite.mockRepo.On("GetUserByID", mock.Anything, "userID").
					Return(&model.User{
						ID: "userID", Email: "test@test.com",
						Password: utils.HashAndSalt([]byte("password")),
					}, nil).Times(1)
			},
			wantErr: true,
		},
		{
			name: "Update fail",
			req:  &dto.ChangePasswordReq{Password: "password", NewPassword: "newPassword"},
			setup: func() {
				suite.mockRepo.On("GetUserByID", mock.Anything, "userID").
					Return(&model.User{
						ID: "userID", Email: "test@test.com",
						Password: utils.HashAndSalt([]byte("password")),
					}, nil).Times(1)
				suite.mockRepo.On("Update", mock.Anything, mock.Anything).Return(errors.New("error")).Times(1)
			},
			wantErr: true,
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.setup()
			err := suite.service.ChangePassword(context.Background(), "userID", tc.req)
			if tc.wantErr {
				suite.NotNil(err)
			} else {
				suite.Nil(err)
			}
		})
	}
}
