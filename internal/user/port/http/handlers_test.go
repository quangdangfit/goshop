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

	domain "goshop/internal/user/domain"
	"goshop/internal/user/model"
	"goshop/internal/user/service/mocks"
	"goshop/pkg/config"
	"goshop/pkg/response"
	"goshop/pkg/utils"
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

func (suite *UserHandlerTestSuite) TestLogin() {
	tests := []struct {
		name      string
		body      any
		setup     func()
		expected  int
		checkBody func(writer *httptest.ResponseRecorder)
	}{
		{
			name: "Success",
			body: &domain.LoginReq{
				Email:    "login@test.com",
				Password: "test123456",
			},
			setup: func() {
				suite.mockService.On("Login", mock.Anything, &domain.LoginReq{
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
			expected: http.StatusOK,
			checkBody: func(writer *httptest.ResponseRecorder) {
				var res response.Response
				var loginRes domain.LoginRes

				_ = json.Unmarshal(writer.Body.Bytes(), &res)
				utils.Copy(&loginRes, &res.Result)
				suite.Equal("login@test.com", loginRes.User.Email)
				suite.Equal("access-token", loginRes.AccessToken)
				suite.Equal("refresh-token", loginRes.RefreshToken)
			},
		},
		{
			name: "InvalidEmailType",
			body: map[string]interface{}{
				"email":    1,
				"password": "test123456",
			},
			setup:    func() {},
			expected: http.StatusBadRequest,
			checkBody: func(writer *httptest.ResponseRecorder) {
				var res map[string]map[string]string
				_ = json.Unmarshal(writer.Body.Bytes(), &res)
				suite.Equal("Invalid request parameters", res["error"]["message"])
			},
		},
		{
			name: "InvalidPasswordType",
			body: map[string]interface{}{
				"email":    "login@test.com",
				"password": 12345,
			},
			setup:    func() {},
			expected: http.StatusBadRequest,
			checkBody: func(writer *httptest.ResponseRecorder) {
				var res map[string]map[string]string
				_ = json.Unmarshal(writer.Body.Bytes(), &res)
				suite.Equal("Invalid request parameters", res["error"]["message"])
			},
		},
		{
			name: "Fail",
			body: &domain.LoginReq{
				Email:    "login@test.com",
				Password: "test123456",
			},
			setup: func() {
				suite.mockService.On("Login", mock.Anything, &domain.LoginReq{
					Email:    "login@test.com",
					Password: "test123456",
				}).Return(nil, "", "", errors.New("error")).Times(1)
			},
			expected: http.StatusInternalServerError,
			checkBody: func(writer *httptest.ResponseRecorder) {
				var res map[string]map[string]string
				_ = json.Unmarshal(writer.Body.Bytes(), &res)
				suite.Equal("Something went wrong", res["error"]["message"])
			},
		},
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			ctx, writer := suite.prepareContext(tc.body)
			tc.setup()
			suite.handler.Login(ctx)
			suite.Equal(tc.expected, writer.Code)
			if tc.checkBody != nil {
				tc.checkBody(writer)
			}
		})
	}
}

// Register
// =================================================================================================

func (suite *UserHandlerTestSuite) TestRegister() {
	tests := []struct {
		name      string
		body      any
		setup     func()
		expected  int
		checkBody func(writer *httptest.ResponseRecorder)
	}{
		{
			name: "Success",
			body: &domain.RegisterReq{
				Email:    "register@test.com",
				Password: "test123456",
			},
			setup: func() {
				suite.mockService.On("Register", mock.Anything, &domain.RegisterReq{
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
			expected: http.StatusOK,
			checkBody: func(writer *httptest.ResponseRecorder) {
				var res response.Response
				var registerRes domain.RegisterRes

				_ = json.Unmarshal(writer.Body.Bytes(), &res)
				utils.Copy(&registerRes, &res.Result)
				suite.Equal("register@test.com", registerRes.User.Email)
			},
		},
		{
			name: "InvalidEmailType",
			body: map[string]interface{}{
				"email":    1,
				"password": "test123456",
			},
			setup:    func() {},
			expected: http.StatusBadRequest,
			checkBody: func(writer *httptest.ResponseRecorder) {
				var res map[string]map[string]string
				_ = json.Unmarshal(writer.Body.Bytes(), &res)
				suite.Equal("Invalid request parameters", res["error"]["message"])
			},
		},
		{
			name: "InvalidPasswordType",
			body: map[string]interface{}{
				"email":    "login@test.com",
				"password": 12345,
			},
			setup:    func() {},
			expected: http.StatusBadRequest,
			checkBody: func(writer *httptest.ResponseRecorder) {
				var res map[string]map[string]string
				_ = json.Unmarshal(writer.Body.Bytes(), &res)
				suite.Equal("Invalid request parameters", res["error"]["message"])
			},
		},
		{
			name: "Fail",
			body: &domain.RegisterReq{
				Email:    "register@test.com",
				Password: "test123456",
			},
			setup: func() {
				suite.mockService.On("Register", mock.Anything, &domain.RegisterReq{
					Email:    "register@test.com",
					Password: "test123456",
				}).Return(nil, errors.New("error")).Times(1)
			},
			expected: http.StatusInternalServerError,
			checkBody: func(writer *httptest.ResponseRecorder) {
				var res map[string]map[string]string
				_ = json.Unmarshal(writer.Body.Bytes(), &res)
				suite.Equal("Something went wrong", res["error"]["message"])
			},
		},
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			ctx, writer := suite.prepareContext(tc.body)
			tc.setup()
			suite.handler.Register(ctx)
			suite.Equal(tc.expected, writer.Code)
			if tc.checkBody != nil {
				tc.checkBody(writer)
			}
		})
	}
}

// GetMe
// =================================================================================================

func (suite *UserHandlerTestSuite) TestGetMe() {
	tests := []struct {
		name      string
		userId    string
		setup     func()
		expected  int
		checkBody func(writer *httptest.ResponseRecorder)
	}{
		{
			name:   "Success",
			userId: "123456",
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
			expected: http.StatusOK,
			checkBody: func(writer *httptest.ResponseRecorder) {
				var res response.Response
				var getMeRes domain.User

				_ = json.Unmarshal(writer.Body.Bytes(), &res)
				utils.Copy(&getMeRes, &res.Result)
				suite.Equal("123456", getMeRes.ID)
				suite.Equal("user@test.com", getMeRes.Email)
			},
		},
		{
			name:     "Unauthorized",
			userId:   "",
			setup:    func() {},
			expected: http.StatusUnauthorized,
			checkBody: func(writer *httptest.ResponseRecorder) {
				var res map[string]map[string]string
				_ = json.Unmarshal(writer.Body.Bytes(), &res)
				suite.Equal("Unauthorized", res["error"]["message"])
			},
		},
		{
			name:   "Fail",
			userId: "123456",
			setup: func() {
				suite.mockService.On("GetUserByID", mock.Anything, "123456").
					Return(nil, errors.New("error")).Times(1)
			},
			expected: http.StatusInternalServerError,
			checkBody: func(writer *httptest.ResponseRecorder) {
				var res map[string]map[string]string
				_ = json.Unmarshal(writer.Body.Bytes(), &res)
				suite.Equal("Something went wrong", res["error"]["message"])
			},
		},
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			ctx, writer := suite.prepareContext(nil)
			if tc.userId != "" {
				ctx.Set("userId", tc.userId)
			}
			tc.setup()
			suite.handler.GetMe(ctx)
			suite.Equal(tc.expected, writer.Code)
			if tc.checkBody != nil {
				tc.checkBody(writer)
			}
		})
	}
}

// Refresh Token
// =================================================================================================

func (suite *UserHandlerTestSuite) TestRefreshToken() {
	tests := []struct {
		name      string
		userId    string
		setup     func()
		expected  int
		checkBody func(writer *httptest.ResponseRecorder)
	}{
		{
			name:   "Success",
			userId: "123456",
			setup: func() {
				suite.mockService.On("RefreshToken", mock.Anything, "123456").
					Return("access-token", nil).Times(1)
			},
			expected: http.StatusOK,
			checkBody: func(writer *httptest.ResponseRecorder) {
				var res response.Response
				var refreshRes domain.RefreshTokenRes

				_ = json.Unmarshal(writer.Body.Bytes(), &res)
				utils.Copy(&refreshRes, &res.Result)
				suite.Equal("access-token", refreshRes.AccessToken)
			},
		},
		{
			name:     "Unauthorized",
			userId:   "",
			setup:    func() {},
			expected: http.StatusUnauthorized,
			checkBody: func(writer *httptest.ResponseRecorder) {
				var res map[string]map[string]string
				_ = json.Unmarshal(writer.Body.Bytes(), &res)
				suite.Equal("Unauthorized", res["error"]["message"])
			},
		},
		{
			name:   "Fail",
			userId: "123456",
			setup: func() {
				suite.mockService.On("RefreshToken", mock.Anything, "123456").
					Return("", errors.New("error")).Times(1)
			},
			expected: http.StatusInternalServerError,
			checkBody: func(writer *httptest.ResponseRecorder) {
				var res map[string]map[string]string
				_ = json.Unmarshal(writer.Body.Bytes(), &res)
				suite.Equal("Something went wrong", res["error"]["message"])
			},
		},
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			ctx, writer := suite.prepareContext(nil)
			if tc.userId != "" {
				ctx.Set("userId", tc.userId)
			}
			tc.setup()
			suite.handler.RefreshToken(ctx)
			suite.Equal(tc.expected, writer.Code)
			if tc.checkBody != nil {
				tc.checkBody(writer)
			}
		})
	}
}

// Change Password
// =================================================================================================

func (suite *UserHandlerTestSuite) TestChangePassword() {
	tests := []struct {
		name      string
		body      any
		userId    string
		setup     func()
		expected  int
		checkBody func(writer *httptest.ResponseRecorder)
	}{
		{
			name: "Success",
			body: &domain.ChangePasswordReq{
				Password:    "test123456",
				NewPassword: "new-test123456",
			},
			userId: "123456",
			setup: func() {
				suite.mockService.On("ChangePassword", mock.Anything, "123456", &domain.ChangePasswordReq{
					Password:    "test123456",
					NewPassword: "new-test123456",
				}).Return(nil).Times(1)
			},
			expected:  http.StatusOK,
			checkBody: nil,
		},
		{
			name: "InvalidPasswordType",
			body: map[string]interface{}{
				"password":     12345,
				"new_password": 12345,
			},
			userId:   "",
			setup:    func() {},
			expected: http.StatusBadRequest,
			checkBody: func(writer *httptest.ResponseRecorder) {
				var res map[string]map[string]string
				_ = json.Unmarshal(writer.Body.Bytes(), &res)
				suite.Equal("Invalid request parameters", res["error"]["message"])
			},
		},
		{
			name: "Fail",
			body: &domain.ChangePasswordReq{
				Password:    "test123456",
				NewPassword: "new-test123456",
			},
			userId: "123456",
			setup: func() {
				suite.mockService.On("ChangePassword", mock.Anything, "123456", &domain.ChangePasswordReq{
					Password:    "test123456",
					NewPassword: "new-test123456",
				}).Return(errors.New("error")).Times(1)
			},
			expected: http.StatusInternalServerError,
			checkBody: func(writer *httptest.ResponseRecorder) {
				var res map[string]map[string]string
				_ = json.Unmarshal(writer.Body.Bytes(), &res)
				suite.Equal("Something went wrong", res["error"]["message"])
			},
		},
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			ctx, writer := suite.prepareContext(tc.body)
			if tc.userId != "" {
				ctx.Set("userId", tc.userId)
			}
			tc.setup()
			suite.handler.ChangePassword(ctx)
			suite.Equal(tc.expected, writer.Code)
			if tc.checkBody != nil {
				tc.checkBody(writer)
			}
		})
	}
}
