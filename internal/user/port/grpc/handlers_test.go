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
	mockService *mocks.IUserService
	handler     *UserHandler
}

func (suite *UserHandlerTestSuite) SetupTest() {
	logger.Initialize(config.ProductionEnv)

	suite.mockService = mocks.NewIUserService(suite.T())
	suite.handler = NewUserHandler(suite.mockService)
}

func TestUserHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(UserHandlerTestSuite))
}

// Login
// =================================================================================================

func (suite *UserHandlerTestSuite) TestUserAPI_LoginSuccess() {
	req := &pb.LoginReq{
		Email:    "login@test.com",
		Password: "test123456",
	}

	suite.mockService.On("Login", mock.Anything, &dto.LoginReq{
		Email:    req.Email,
		Password: req.Password,
	}).Return(
		&model.User{
			Email:    "login@test.com",
			Password: "test123456",
		},
		"access-token",
		"refresh-token",
		nil,
	).Times(1)

	res, err := suite.handler.Login(context.Background(), req)

	suite.Nil(err)
	suite.Equal(req.Email, res.User.Email)
	suite.Equal("access-token", res.AccessToken)
	suite.Equal("refresh-token", res.RefreshToken)
}

func (suite *UserHandlerTestSuite) TestUserAPI_LoginFail() {
	req := &pb.LoginReq{
		Email:    "login@test.com",
		Password: "test123456",
	}

	suite.mockService.On("Login", mock.Anything, &dto.LoginReq{
		Email:    req.Email,
		Password: req.Password,
	}).Return(nil, "", "", errors.New("error")).Times(1)

	res, err := suite.handler.Login(context.Background(), req)
	suite.Nil(res)
	suite.NotNil(err)
}

// Register
// =================================================================================================

func (suite *UserHandlerTestSuite) TestUserAPI_RegisterSuccess() {
	req := &pb.RegisterReq{
		Email:    "register@test.com",
		Password: "test123456",
	}

	suite.mockService.On("Register", mock.Anything, &dto.RegisterReq{
		Email:    req.Email,
		Password: req.Password,
	}).Return(
		&model.User{
			Email:    "register@test.com",
			Password: "test123456",
		},
		nil,
	).Times(1)

	res, err := suite.handler.Register(context.Background(), req)

	suite.Nil(err)
	suite.Equal(req.Email, res.User.Email)
}

func (suite *UserHandlerTestSuite) TestUserAPI_RegisterFail() {
	req := &pb.RegisterReq{
		Email:    "register@test.com",
		Password: "test123456",
	}

	suite.mockService.On("Register", mock.Anything, &dto.RegisterReq{
		Email:    req.Email,
		Password: req.Password,
	}).Return(nil, errors.New("error")).Times(1)

	res, err := suite.handler.Register(context.Background(), req)
	suite.Nil(res)
	suite.NotNil(err)
}

//// GetMe
//// =================================================================================================

func (suite *UserHandlerTestSuite) TestUserAPI_GetMeSuccess() {
	userId := "123456"
	ctx := context.WithValue(context.Background(), "userId", userId)

	suite.mockService.On("GetUserByID", mock.Anything, userId).
		Return(
			&model.User{
				ID:       userId,
				Email:    "user@test.com",
				Password: "test123456",
			},
			nil,
		).Times(1)

	res, err := suite.handler.GetMe(ctx, &pb.GetMeReq{})
	suite.Nil(err)
	suite.Equal(userId, res.User.Id)
	suite.Equal("user@test.com", res.User.Email)
}

func (suite *UserHandlerTestSuite) TestUserAPI_GetMeUnauthorized() {
	res, err := suite.handler.GetMe(context.Background(), &pb.GetMeReq{})
	suite.Nil(res)
	suite.NotNil(err)
}

func (suite *UserHandlerTestSuite) TestUserAPI_GetMeFail() {
	userId := "123456"
	ctx := context.WithValue(context.Background(), "userId", userId)

	suite.mockService.On("GetUserByID", mock.Anything, "123456").
		Return(nil, errors.New("error")).Times(1)

	res, err := suite.handler.GetMe(ctx, &pb.GetMeReq{})
	suite.Nil(res)
	suite.NotNil(err)
}

//// Refresh Token
//// =================================================================================================

func (suite *UserHandlerTestSuite) TestUserAPI_RefreshTokenSuccess() {
	userId := "123456"
	ctx := context.WithValue(context.Background(), "userId", userId)

	suite.mockService.On("RefreshToken", mock.Anything, userId).
		Return("access-token", nil).Times(1)

	res, err := suite.handler.RefreshToken(ctx, &pb.RefreshTokenReq{})
	suite.Nil(err)
	suite.Equal("access-token", res.AccessToken)
}

func (suite *UserHandlerTestSuite) TestUserAPI_RefreshTokenUnauthorized() {
	res, err := suite.handler.RefreshToken(context.Background(), &pb.RefreshTokenReq{})
	suite.Nil(res)
	suite.NotNil(err)
}

func (suite *UserHandlerTestSuite) TestUserAPI_RefreshTokenFail() {
	userId := "123456"
	ctx := context.WithValue(context.Background(), "userId", userId)

	suite.mockService.On("RefreshToken", mock.Anything, "123456").
		Return("", errors.New("error")).Times(1)

	res, err := suite.handler.RefreshToken(ctx, &pb.RefreshTokenReq{})
	suite.Nil(res)
	suite.NotNil(err)
}

//// Change Password
//// =================================================================================================

func (suite *UserHandlerTestSuite) TestUserAPI_ChangePasswordSuccess() {
	req := &pb.ChangePasswordReq{
		Password:    "test123456",
		NewPassword: "new-test123456",
	}

	userId := "123456"
	ctx := context.WithValue(context.Background(), "userId", userId)

	suite.mockService.On("ChangePassword", mock.Anything, userId, &dto.ChangePasswordReq{
		Password:    req.Password,
		NewPassword: req.NewPassword,
	}).Return(nil).Times(1)

	res, err := suite.handler.ChangePassword(ctx, req)
	suite.NotNil(res)
	suite.Nil(err)
}

func (suite *UserHandlerTestSuite) TestUserAPI_ChangePasswordUnauthorized() {
	res, err := suite.handler.ChangePassword(context.Background(), &pb.ChangePasswordReq{})
	suite.Nil(res)
	suite.NotNil(err)
}

func (suite *UserHandlerTestSuite) TestUserAPI_ChangePasswordFail() {
	req := &pb.ChangePasswordReq{
		Password:    "test123456",
		NewPassword: "new-test123456",
	}

	userId := "123456"
	ctx := context.WithValue(context.Background(), "userId", userId)

	suite.mockService.On("ChangePassword", mock.Anything, "123456", &dto.ChangePasswordReq{
		Password:    req.Password,
		NewPassword: req.NewPassword,
	}).Return(errors.New("error")).Times(1)

	res, err := suite.handler.ChangePassword(ctx, req)
	suite.Nil(res)
	suite.NotNil(err)
}
