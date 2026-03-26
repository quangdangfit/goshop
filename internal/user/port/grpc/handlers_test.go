package grpc

import (
	"context"
	"errors"
	"testing"

	"github.com/quangdangfit/gocommon/logger"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"goshop/internal/user/dto"
	"goshop/internal/user/model"
	"goshop/internal/user/service/mocks"
	"goshop/pkg/config"
	pb "goshop/proto/gen/go/user"
)

type UserHandlerTestSuite struct {
	suite.Suite
	mockService *mocks.UserService
	handler     *UserHandler
}

func (suite *UserHandlerTestSuite) SetupTest() {
	logger.Initialize(config.ProductionEnv)

	suite.mockService = mocks.NewUserService(suite.T())
	suite.handler = NewUserHandler(suite.mockService)
}

func TestUserHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(UserHandlerTestSuite))
}

// Login
// =================================================================================================

func (suite *UserHandlerTestSuite) TestLogin() {
	tests := []struct {
		name      string
		setup     func()
		req       *pb.LoginReq
		expectNil bool
		expectErr bool
		validate  func(res *pb.LoginRes)
	}{
		{
			name: "Success",
			setup: func() {
				suite.mockService.On("Login", mock.Anything, &dto.LoginReq{
					Email:    "login@test.com",
					Password: "test123456",
				}).Return(
					&model.User{
						Email:    "login@test.com",
						Password: "test123456",
					},
					"access-token",
					"refresh-token",
					nil,
				).Times(1)
			},
			req: &pb.LoginReq{
				Email:    "login@test.com",
				Password: "test123456",
			},
			expectNil: false,
			expectErr: false,
			validate: func(res *pb.LoginRes) {
				suite.Equal("login@test.com", res.User.Email)
				suite.Equal("access-token", res.AccessToken)
				suite.Equal("refresh-token", res.RefreshToken)
			},
		},
		{
			name: "Fail",
			setup: func() {
				suite.mockService.On("Login", mock.Anything, &dto.LoginReq{
					Email:    "login@test.com",
					Password: "test123456",
				}).Return(nil, "", "", errors.New("error")).Times(1)
			},
			req: &pb.LoginReq{
				Email:    "login@test.com",
				Password: "test123456",
			},
			expectNil: true,
			expectErr: true,
		},
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.setup()

			res, err := suite.handler.Login(context.Background(), tc.req)

			if tc.expectErr {
				suite.NotNil(err)
			} else {
				suite.Nil(err)
			}
			if tc.expectNil {
				suite.Nil(res)
			} else {
				suite.NotNil(res)
				if tc.validate != nil {
					tc.validate(res)
				}
			}
		})
	}
}

// Register
// =================================================================================================

func (suite *UserHandlerTestSuite) TestRegister() {
	tests := []struct {
		name      string
		setup     func()
		req       *pb.RegisterReq
		expectNil bool
		expectErr bool
		validate  func(res *pb.RegisterRes)
	}{
		{
			name: "Success",
			setup: func() {
				suite.mockService.On("Register", mock.Anything, &dto.RegisterReq{
					Email:    "register@test.com",
					Password: "test123456",
				}).Return(
					&model.User{
						Email:    "register@test.com",
						Password: "test123456",
					},
					nil,
				).Times(1)
			},
			req: &pb.RegisterReq{
				Email:    "register@test.com",
				Password: "test123456",
			},
			expectNil: false,
			expectErr: false,
			validate: func(res *pb.RegisterRes) {
				suite.Equal("register@test.com", res.User.Email)
			},
		},
		{
			name: "Fail",
			setup: func() {
				suite.mockService.On("Register", mock.Anything, &dto.RegisterReq{
					Email:    "register@test.com",
					Password: "test123456",
				}).Return(nil, errors.New("error")).Times(1)
			},
			req: &pb.RegisterReq{
				Email:    "register@test.com",
				Password: "test123456",
			},
			expectNil: true,
			expectErr: true,
		},
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.setup()

			res, err := suite.handler.Register(context.Background(), tc.req)

			if tc.expectErr {
				suite.NotNil(err)
			} else {
				suite.Nil(err)
			}
			if tc.expectNil {
				suite.Nil(res)
			} else {
				suite.NotNil(res)
				if tc.validate != nil {
					tc.validate(res)
				}
			}
		})
	}
}

// GetMe
// =================================================================================================

func (suite *UserHandlerTestSuite) TestGetMe() {
	tests := []struct {
		name      string
		setup     func()
		ctx       context.Context
		expectNil bool
		expectErr bool
		validate  func(res *pb.GetMeRes)
	}{
		{
			name: "Success",
			setup: func() {
				suite.mockService.On("GetUserByID", mock.Anything, "123456").
					Return(
						&model.User{
							ID:       "123456",
							Email:    "user@test.com",
							Password: "test123456",
						},
						nil,
					).Times(1)
			},
			ctx:       context.WithValue(context.Background(), "userId", "123456"),
			expectNil: false,
			expectErr: false,
			validate: func(res *pb.GetMeRes) {
				suite.Equal("123456", res.User.Id)
				suite.Equal("user@test.com", res.User.Email)
			},
		},
		{
			name:      "Unauthorized",
			setup:     func() {},
			ctx:       context.Background(),
			expectNil: true,
			expectErr: true,
		},
		{
			name: "Fail",
			setup: func() {
				suite.mockService.On("GetUserByID", mock.Anything, "123456").
					Return(nil, errors.New("error")).Times(1)
			},
			ctx:       context.WithValue(context.Background(), "userId", "123456"),
			expectNil: true,
			expectErr: true,
		},
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.setup()

			res, err := suite.handler.GetMe(tc.ctx, &pb.GetMeReq{})

			if tc.expectErr {
				suite.NotNil(err)
			} else {
				suite.Nil(err)
			}
			if tc.expectNil {
				suite.Nil(res)
			} else {
				suite.NotNil(res)
				if tc.validate != nil {
					tc.validate(res)
				}
			}
		})
	}
}

// RefreshToken
// =================================================================================================

func (suite *UserHandlerTestSuite) TestRefreshToken() {
	tests := []struct {
		name      string
		setup     func()
		ctx       context.Context
		expectNil bool
		expectErr bool
		validate  func(res *pb.RefreshTokenRes)
	}{
		{
			name: "Success",
			setup: func() {
				suite.mockService.On("RefreshToken", mock.Anything, "123456").
					Return("access-token", nil).Times(1)
			},
			ctx:       context.WithValue(context.Background(), "userId", "123456"),
			expectNil: false,
			expectErr: false,
			validate: func(res *pb.RefreshTokenRes) {
				suite.Equal("access-token", res.AccessToken)
			},
		},
		{
			name:      "Unauthorized",
			setup:     func() {},
			ctx:       context.Background(),
			expectNil: true,
			expectErr: true,
		},
		{
			name: "Fail",
			setup: func() {
				suite.mockService.On("RefreshToken", mock.Anything, "123456").
					Return("", errors.New("error")).Times(1)
			},
			ctx:       context.WithValue(context.Background(), "userId", "123456"),
			expectNil: true,
			expectErr: true,
		},
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.setup()

			res, err := suite.handler.RefreshToken(tc.ctx, &pb.RefreshTokenReq{})

			if tc.expectErr {
				suite.NotNil(err)
			} else {
				suite.Nil(err)
			}
			if tc.expectNil {
				suite.Nil(res)
			} else {
				suite.NotNil(res)
				if tc.validate != nil {
					tc.validate(res)
				}
			}
		})
	}
}

// ChangePassword
// =================================================================================================

func (suite *UserHandlerTestSuite) TestChangePassword() {
	tests := []struct {
		name      string
		setup     func()
		req       *pb.ChangePasswordReq
		ctx       context.Context
		expectNil bool
		expectErr bool
	}{
		{
			name: "Success",
			setup: func() {
				suite.mockService.On("ChangePassword", mock.Anything, "123456", &dto.ChangePasswordReq{
					Password:    "test123456",
					NewPassword: "new-test123456",
				}).Return(nil).Times(1)
			},
			req: &pb.ChangePasswordReq{
				Password:    "test123456",
				NewPassword: "new-test123456",
			},
			ctx:       context.WithValue(context.Background(), "userId", "123456"),
			expectNil: false,
			expectErr: false,
		},
		{
			name:      "Unauthorized",
			setup:     func() {},
			req:       &pb.ChangePasswordReq{},
			ctx:       context.Background(),
			expectNil: true,
			expectErr: true,
		},
		{
			name: "Fail",
			setup: func() {
				suite.mockService.On("ChangePassword", mock.Anything, "123456", &dto.ChangePasswordReq{
					Password:    "test123456",
					NewPassword: "new-test123456",
				}).Return(errors.New("error")).Times(1)
			},
			req: &pb.ChangePasswordReq{
				Password:    "test123456",
				NewPassword: "new-test123456",
			},
			ctx:       context.WithValue(context.Background(), "userId", "123456"),
			expectNil: true,
			expectErr: true,
		},
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.setup()

			res, err := suite.handler.ChangePassword(tc.ctx, tc.req)

			if tc.expectErr {
				suite.NotNil(err)
			} else {
				suite.Nil(err)
			}
			if tc.expectNil {
				suite.Nil(res)
			} else {
				suite.NotNil(res)
			}
		})
	}
}
