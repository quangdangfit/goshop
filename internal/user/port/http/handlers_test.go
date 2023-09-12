package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/quangdangfit/gocommon/logger"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"goshop/internal/user/dto"
	"goshop/internal/user/model"
	"goshop/internal/user/service/mocks"
	"goshop/pkg/config"
	"goshop/pkg/response"
	"goshop/pkg/utils"
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

func (suite *UserHandlerTestSuite) prepareContext(body any) (*gin.Context, *httptest.ResponseRecorder) {
	requestBody, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("", "/", bytes.NewBuffer(requestBody))
	c, _ := gin.CreateTestContext(w)
	c.Request = r

	return c, w
}

// Login
// =================================================================================================

func (suite *UserHandlerTestSuite) TestLoginSuccess() {
	req := &dto.LoginReq{
		Email:    "login@test.com",
		Password: "test123456",
	}

	ctx, writer := suite.prepareContext(req)

	suite.mockService.On("Login", mock.Anything, req).
		Return(
			&model.User{
				Email:    "login@test.com",
				Password: "test123456",
			},
			"access-token",
			"refresh-token",
			nil,
		).Times(1)

	suite.handler.Login(ctx)

	var res response.Response
	var loginRes dto.LoginRes

	_ = json.Unmarshal(writer.Body.Bytes(), &res)
	utils.Copy(&loginRes, &res.Result)
	suite.Equal(http.StatusOK, writer.Code)
	suite.Equal(req.Email, loginRes.User.Email)
	suite.Equal("access-token", loginRes.AccessToken)
	suite.Equal("refresh-token", loginRes.RefreshToken)
}

func (suite *UserHandlerTestSuite) TestLoginInvalidEmailType() {
	req := map[string]interface{}{
		"email":    1,
		"password": "test123456",
	}

	ctx, writer := suite.prepareContext(req)

	suite.handler.Login(ctx)

	var res map[string]map[string]string
	_ = json.Unmarshal(writer.Body.Bytes(), &res)
	suite.Equal(http.StatusBadRequest, writer.Code)
	suite.Equal("Invalid parameters", res["error"]["message"])
}

func (suite *UserHandlerTestSuite) TestLoginInvalidPasswordType() {
	req := map[string]interface{}{
		"email":    "login@test.com",
		"password": 12345,
	}

	ctx, writer := suite.prepareContext(req)

	suite.handler.Login(ctx)

	var res map[string]map[string]string
	_ = json.Unmarshal(writer.Body.Bytes(), &res)
	suite.Equal(http.StatusBadRequest, writer.Code)
	suite.Equal("Invalid parameters", res["error"]["message"])
}

func (suite *UserHandlerTestSuite) TestLoginFail() {
	req := &dto.LoginReq{
		Email:    "login@test.com",
		Password: "test123456",
	}

	ctx, writer := suite.prepareContext(req)

	suite.mockService.On("Login", mock.Anything, req).
		Return(nil, "", "", errors.New("error")).Times(1)

	suite.handler.Login(ctx)

	var res map[string]map[string]string
	_ = json.Unmarshal(writer.Body.Bytes(), &res)
	suite.Equal(http.StatusInternalServerError, writer.Code)
	suite.Equal("Something went wrong", res["error"]["message"])
}

// Register
// =================================================================================================

func (suite *UserHandlerTestSuite) TestRegisterSuccess() {
	req := &dto.RegisterReq{
		Email:    "register@test.com",
		Password: "test123456",
	}

	ctx, writer := suite.prepareContext(req)

	suite.mockService.On("Register", mock.Anything, req).
		Return(
			&model.User{
				Email:    "register@test.com",
				Password: "test123456",
			},
			nil,
		).Times(1)

	suite.handler.Register(ctx)

	var res response.Response
	var registerRes dto.RegisterRes

	_ = json.Unmarshal(writer.Body.Bytes(), &res)
	utils.Copy(&registerRes, &res.Result)
	suite.Equal(http.StatusOK, writer.Code)
	suite.Equal(req.Email, registerRes.User.Email)
}

func (suite *UserHandlerTestSuite) TestRegisterInvalidEmailType() {
	req := map[string]interface{}{
		"email":    1,
		"password": "test123456",
	}

	ctx, writer := suite.prepareContext(req)

	suite.handler.Register(ctx)

	var res map[string]map[string]string
	_ = json.Unmarshal(writer.Body.Bytes(), &res)
	suite.Equal(http.StatusBadRequest, writer.Code)
	suite.Equal("Invalid parameters", res["error"]["message"])
}

func (suite *UserHandlerTestSuite) TestRegisterInvalidPasswordType() {
	req := map[string]interface{}{
		"email":    "login@test.com",
		"password": 12345,
	}

	ctx, writer := suite.prepareContext(req)

	suite.handler.Register(ctx)

	var res map[string]map[string]string
	_ = json.Unmarshal(writer.Body.Bytes(), &res)
	suite.Equal(http.StatusBadRequest, writer.Code)
	suite.Equal("Invalid parameters", res["error"]["message"])
}

func (suite *UserHandlerTestSuite) TestRegisterFail() {
	req := &dto.RegisterReq{
		Email:    "register@test.com",
		Password: "test123456",
	}

	ctx, writer := suite.prepareContext(req)

	suite.mockService.On("Register", mock.Anything, req).
		Return(nil, errors.New("error")).Times(1)

	suite.handler.Register(ctx)

	var res map[string]map[string]string
	_ = json.Unmarshal(writer.Body.Bytes(), &res)
	suite.Equal(http.StatusInternalServerError, writer.Code)
	suite.Equal("Something went wrong", res["error"]["message"])
}

// GetMe
// =================================================================================================

func (suite *UserHandlerTestSuite) TestGetMeSuccess() {
	ctx, writer := suite.prepareContext(nil)
	ctx.Set("userId", "123456")

	suite.mockService.On("GetUserByID", mock.Anything, "123456").
		Return(
			&model.User{
				ID:       "123456",
				Email:    "user@test.com",
				Password: "test123456",
			},
			nil,
		).Times(1)

	suite.handler.GetMe(ctx)

	var res response.Response
	var getMeRes dto.User

	_ = json.Unmarshal(writer.Body.Bytes(), &res)
	utils.Copy(&getMeRes, &res.Result)
	suite.Equal(http.StatusOK, writer.Code)
	suite.Equal("123456", getMeRes.ID)
	suite.Equal("user@test.com", getMeRes.Email)
}

func (suite *UserHandlerTestSuite) TestGetMeUnauthorized() {
	ctx, writer := suite.prepareContext(nil)
	suite.handler.GetMe(ctx)

	var res map[string]map[string]string
	_ = json.Unmarshal(writer.Body.Bytes(), &res)
	suite.Equal(http.StatusUnauthorized, writer.Code)
	suite.Equal("Unauthorized", res["error"]["message"])
}

func (suite *UserHandlerTestSuite) TestGetMeFail() {
	ctx, writer := suite.prepareContext(nil)
	ctx.Set("userId", "123456")

	suite.mockService.On("GetUserByID", mock.Anything, "123456").
		Return(nil, errors.New("error")).Times(1)

	suite.handler.GetMe(ctx)

	var res map[string]map[string]string
	_ = json.Unmarshal(writer.Body.Bytes(), &res)
	suite.Equal(http.StatusInternalServerError, writer.Code)
	suite.Equal("Something went wrong", res["error"]["message"])
}

// Refresh Token
// =================================================================================================

func (suite *UserHandlerTestSuite) TestRefreshTokenSuccess() {
	ctx, writer := suite.prepareContext(nil)
	ctx.Set("userId", "123456")

	suite.mockService.On("RefreshToken", mock.Anything, "123456").
		Return("access-token", nil).Times(1)

	suite.handler.RefreshToken(ctx)

	var res response.Response
	var getMeRes dto.RefreshTokenRes

	_ = json.Unmarshal(writer.Body.Bytes(), &res)
	utils.Copy(&getMeRes, &res.Result)
	suite.Equal(http.StatusOK, writer.Code)
	suite.Equal("access-token", getMeRes.AccessToken)
}

func (suite *UserHandlerTestSuite) TestRefreshTokenUnauthorized() {
	ctx, writer := suite.prepareContext(nil)
	suite.handler.RefreshToken(ctx)

	var res map[string]map[string]string
	_ = json.Unmarshal(writer.Body.Bytes(), &res)
	suite.Equal(http.StatusUnauthorized, writer.Code)
	suite.Equal("Unauthorized", res["error"]["message"])
}

func (suite *UserHandlerTestSuite) TestRefreshTokenFail() {
	ctx, writer := suite.prepareContext(nil)
	ctx.Set("userId", "123456")

	suite.mockService.On("RefreshToken", mock.Anything, "123456").
		Return("", errors.New("error")).Times(1)

	suite.handler.RefreshToken(ctx)

	var res map[string]map[string]string
	_ = json.Unmarshal(writer.Body.Bytes(), &res)
	suite.Equal(http.StatusInternalServerError, writer.Code)
	suite.Equal("Something went wrong", res["error"]["message"])
}

// Change Password
// =================================================================================================

func (suite *UserHandlerTestSuite) TestChangePasswordSuccess() {
	req := &dto.ChangePasswordReq{
		Password:    "test123456",
		NewPassword: "new-test123456",
	}

	ctx, writer := suite.prepareContext(req)
	ctx.Set("userId", "123456")

	suite.mockService.On("ChangePassword", mock.Anything, "123456", req).
		Return(nil).Times(1)

	suite.handler.ChangePassword(ctx)

	suite.Equal(http.StatusOK, writer.Code)
}

func (suite *UserHandlerTestSuite) TestChangePasswordInvalidPasswordType() {
	req := map[string]interface{}{
		"password":     12345,
		"new_password": 12345,
	}

	ctx, writer := suite.prepareContext(req)
	suite.handler.ChangePassword(ctx)

	var res map[string]map[string]string
	_ = json.Unmarshal(writer.Body.Bytes(), &res)
	suite.Equal(http.StatusBadRequest, writer.Code)
	suite.Equal("Invalid parameters", res["error"]["message"])
}

func (suite *UserHandlerTestSuite) TestChangePasswordFail() {
	req := &dto.ChangePasswordReq{
		Password:    "test123456",
		NewPassword: "new-test123456",
	}

	ctx, writer := suite.prepareContext(req)
	ctx.Set("userId", "123456")

	suite.mockService.On("ChangePassword", mock.Anything, "123456", req).
		Return(errors.New("error")).Times(1)

	suite.handler.ChangePassword(ctx)

	var res map[string]map[string]string
	_ = json.Unmarshal(writer.Body.Bytes(), &res)
	suite.Equal(http.StatusInternalServerError, writer.Code)
	suite.Equal("Something went wrong", res["error"]["message"])
}
