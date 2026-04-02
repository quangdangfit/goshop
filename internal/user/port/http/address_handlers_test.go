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
	srvMocks "goshop/internal/user/service/mocks"
	"goshop/pkg/config"
	"goshop/pkg/response"
	"goshop/pkg/utils"
)

type AddressHandlerTestSuite struct {
	suite.Suite
	mockService *srvMocks.AddressService
	handler     *AddressHandler
}

func (suite *AddressHandlerTestSuite) SetupTest() {
	logger.Initialize(config.ProductionEnv)
	suite.mockService = srvMocks.NewAddressService(suite.T())
	suite.handler = NewAddressHandler(suite.mockService)
}

func TestAddressHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(AddressHandlerTestSuite))
}

func (suite *AddressHandlerTestSuite) prepareContext(method, path string, body any) (*gin.Context, *httptest.ResponseRecorder) {
	requestBody, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(method, path, bytes.NewBuffer(requestBody))
	c, _ := gin.CreateTestContext(w)
	c.Request = r
	return c, w
}

// ListAddresses
// =================================================================================================

func (suite *AddressHandlerTestSuite) TestListAddresses() {
	tests := []struct {
		name      string
		userId    string
		setup     func()
		expected  int
		checkBody func(writer *httptest.ResponseRecorder)
	}{
		{
			name:   "Success",
			userId: "u1",
			setup: func() {
				suite.mockService.On("ListAddresses", mock.Anything, "u1").
					Return([]*model.Address{
						{ID: "a1", UserID: "u1", Street: "123 Main St"},
						{ID: "a2", UserID: "u1", Street: "456 Oak Ave"},
					}, nil).Times(1)
			},
			expected: http.StatusOK,
			checkBody: func(writer *httptest.ResponseRecorder) {
				var res response.Response
				var addresses []*domain.Address
				_ = json.Unmarshal(writer.Body.Bytes(), &res)
				utils.Copy(&addresses, &res.Result)
				suite.Equal(2, len(addresses))
			},
		},
		{
			name:     "Unauthorized",
			userId:   "",
			setup:    func() {},
			expected: http.StatusUnauthorized,
		},
		{
			name:   "Fail",
			userId: "u1",
			setup: func() {
				suite.mockService.On("ListAddresses", mock.Anything, "u1").
					Return(nil, errors.New("db error")).Times(1)
			},
			expected: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			ctx, writer := suite.prepareContext("GET", "/api/v1/addresses", nil)
			if tc.userId != "" {
				ctx.Set("userId", tc.userId)
			}
			tc.setup()
			suite.handler.ListAddresses(ctx)
			suite.Equal(tc.expected, writer.Code)
			if tc.checkBody != nil {
				tc.checkBody(writer)
			}
		})
	}
}

// GetAddressByID
// =================================================================================================

func (suite *AddressHandlerTestSuite) TestGetAddressByID() {
	tests := []struct {
		name      string
		userId    string
		params    gin.Params
		setup     func()
		expected  int
		checkBody func(writer *httptest.ResponseRecorder)
	}{
		{
			name:   "Success",
			userId: "u1",
			params: gin.Params{{Key: "id", Value: "a1"}},
			setup: func() {
				suite.mockService.On("GetAddressByID", mock.Anything, "a1", "u1").
					Return(&model.Address{ID: "a1", UserID: "u1", Street: "123 Main St"}, nil).Times(1)
			},
			expected: http.StatusOK,
			checkBody: func(writer *httptest.ResponseRecorder) {
				var res response.Response
				var address domain.Address
				_ = json.Unmarshal(writer.Body.Bytes(), &res)
				utils.Copy(&address, &res.Result)
				suite.Equal("a1", address.ID)
			},
		},
		{
			name:     "Unauthorized",
			userId:   "",
			setup:    func() {},
			expected: http.StatusUnauthorized,
		},
		{
			name:     "MissingID",
			userId:   "u1",
			setup:    func() {},
			expected: http.StatusBadRequest,
		},
		{
			name:   "NotFound",
			userId: "u1",
			params: gin.Params{{Key: "id", Value: "notfound"}},
			setup: func() {
				suite.mockService.On("GetAddressByID", mock.Anything, "notfound", "u1").
					Return(nil, errors.New("not found")).Times(1)
			},
			expected: http.StatusNotFound,
		},
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			ctx, writer := suite.prepareContext("GET", "/api/v1/addresses/a1", nil)
			if tc.userId != "" {
				ctx.Set("userId", tc.userId)
			}
			if tc.params != nil {
				ctx.Params = tc.params
			}
			tc.setup()
			suite.handler.GetAddressByID(ctx)
			suite.Equal(tc.expected, writer.Code)
			if tc.checkBody != nil {
				tc.checkBody(writer)
			}
		})
	}
}

// CreateAddress
// =================================================================================================

func (suite *AddressHandlerTestSuite) TestCreateAddress() {
	tests := []struct {
		name      string
		body      any
		userId    string
		setup     func()
		expected  int
		checkBody func(writer *httptest.ResponseRecorder)
	}{
		{
			name:   "Success",
			body:   &domain.CreateAddressReq{Street: "123 Main St", City: "Springfield"},
			userId: "u1",
			setup: func() {
				suite.mockService.On("Create", mock.Anything, "u1", &domain.CreateAddressReq{Street: "123 Main St", City: "Springfield"}).
					Return(&model.Address{ID: "a1", UserID: "u1", Street: "123 Main St"}, nil).Times(1)
			},
			expected: http.StatusOK,
			checkBody: func(writer *httptest.ResponseRecorder) {
				var res response.Response
				var address domain.Address
				_ = json.Unmarshal(writer.Body.Bytes(), &res)
				utils.Copy(&address, &res.Result)
				suite.Equal("a1", address.ID)
			},
		},
		{
			name:     "Unauthorized",
			body:     &domain.CreateAddressReq{Street: "123 Main St"},
			userId:   "",
			setup:    func() {},
			expected: http.StatusUnauthorized,
		},
		{
			name:     "InvalidBody",
			body:     map[string]any{"street": 123},
			userId:   "u1",
			setup:    func() {},
			expected: http.StatusBadRequest,
		},
		{
			name:   "Fail",
			body:   &domain.CreateAddressReq{Street: "123 Main St", City: "Springfield"},
			userId: "u1",
			setup: func() {
				suite.mockService.On("Create", mock.Anything, "u1", &domain.CreateAddressReq{Street: "123 Main St", City: "Springfield"}).
					Return(nil, errors.New("db error")).Times(1)
			},
			expected: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			ctx, writer := suite.prepareContext("POST", "/api/v1/addresses", tc.body)
			if tc.userId != "" {
				ctx.Set("userId", tc.userId)
			}
			tc.setup()
			suite.handler.CreateAddress(ctx)
			suite.Equal(tc.expected, writer.Code)
			if tc.checkBody != nil {
				tc.checkBody(writer)
			}
		})
	}
}

// UpdateAddress
// =================================================================================================

func (suite *AddressHandlerTestSuite) TestUpdateAddress() {
	tests := []struct {
		name      string
		body      any
		userId    string
		params    gin.Params
		setup     func()
		expected  int
		checkBody func(writer *httptest.ResponseRecorder)
	}{
		{
			name:   "Success",
			body:   &domain.UpdateAddressReq{Street: "789 Elm St"},
			userId: "u1",
			params: gin.Params{{Key: "id", Value: "a1"}},
			setup: func() {
				suite.mockService.On("Update", mock.Anything, "a1", "u1", &domain.UpdateAddressReq{Street: "789 Elm St"}).
					Return(&model.Address{ID: "a1", UserID: "u1", Street: "789 Elm St"}, nil).Times(1)
			},
			expected: http.StatusOK,
			checkBody: func(writer *httptest.ResponseRecorder) {
				var res response.Response
				var address domain.Address
				_ = json.Unmarshal(writer.Body.Bytes(), &res)
				utils.Copy(&address, &res.Result)
				suite.Equal("a1", address.ID)
			},
		},
		{
			name:     "Unauthorized",
			body:     &domain.UpdateAddressReq{Street: "789 Elm St"},
			userId:   "",
			setup:    func() {},
			expected: http.StatusUnauthorized,
		},
		{
			name:     "InvalidBody",
			body:     map[string]any{"street": 123},
			userId:   "u1",
			params:   gin.Params{{Key: "id", Value: "a1"}},
			setup:    func() {},
			expected: http.StatusBadRequest,
		},
		{
			name:   "Fail",
			body:   &domain.UpdateAddressReq{Street: "789 Elm St"},
			userId: "u1",
			params: gin.Params{{Key: "id", Value: "a1"}},
			setup: func() {
				suite.mockService.On("Update", mock.Anything, "a1", "u1", &domain.UpdateAddressReq{Street: "789 Elm St"}).
					Return(nil, errors.New("not found")).Times(1)
			},
			expected: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			ctx, writer := suite.prepareContext("PUT", "/api/v1/addresses/a1", tc.body)
			if tc.userId != "" {
				ctx.Set("userId", tc.userId)
			}
			if tc.params != nil {
				ctx.Params = tc.params
			}
			tc.setup()
			suite.handler.UpdateAddress(ctx)
			suite.Equal(tc.expected, writer.Code)
			if tc.checkBody != nil {
				tc.checkBody(writer)
			}
		})
	}
}

// DeleteAddress
// =================================================================================================

func (suite *AddressHandlerTestSuite) TestDeleteAddress() {
	tests := []struct {
		name     string
		userId   string
		params   gin.Params
		setup    func()
		expected int
	}{
		{
			name:   "Success",
			userId: "u1",
			params: gin.Params{{Key: "id", Value: "a1"}},
			setup: func() {
				suite.mockService.On("Delete", mock.Anything, "a1", "u1").Return(nil).Times(1)
			},
			expected: http.StatusOK,
		},
		{
			name:     "Unauthorized",
			userId:   "",
			setup:    func() {},
			expected: http.StatusUnauthorized,
		},
		{
			name:   "Fail",
			userId: "u1",
			params: gin.Params{{Key: "id", Value: "a1"}},
			setup: func() {
				suite.mockService.On("Delete", mock.Anything, "a1", "u1").Return(errors.New("not found")).Times(1)
			},
			expected: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			ctx, writer := suite.prepareContext("DELETE", "/api/v1/addresses/a1", nil)
			if tc.userId != "" {
				ctx.Set("userId", tc.userId)
			}
			if tc.params != nil {
				ctx.Params = tc.params
			}
			tc.setup()
			suite.handler.DeleteAddress(ctx)
			suite.Equal(tc.expected, writer.Code)
		})
	}
}

// SetDefaultAddress
// =================================================================================================

func (suite *AddressHandlerTestSuite) TestSetDefaultAddress() {
	tests := []struct {
		name     string
		userId   string
		params   gin.Params
		setup    func()
		expected int
	}{
		{
			name:   "Success",
			userId: "u1",
			params: gin.Params{{Key: "id", Value: "a1"}},
			setup: func() {
				suite.mockService.On("SetDefault", mock.Anything, "a1", "u1").Return(nil).Times(1)
			},
			expected: http.StatusOK,
		},
		{
			name:     "Unauthorized",
			userId:   "",
			setup:    func() {},
			expected: http.StatusUnauthorized,
		},
		{
			name:   "Fail",
			userId: "u1",
			params: gin.Params{{Key: "id", Value: "a1"}},
			setup: func() {
				suite.mockService.On("SetDefault", mock.Anything, "a1", "u1").Return(errors.New("not found")).Times(1)
			},
			expected: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			ctx, writer := suite.prepareContext("PUT", "/api/v1/addresses/a1/default", nil)
			if tc.userId != "" {
				ctx.Set("userId", tc.userId)
			}
			if tc.params != nil {
				ctx.Params = tc.params
			}
			tc.setup()
			suite.handler.SetDefaultAddress(ctx)
			suite.Equal(tc.expected, writer.Code)
		})
	}
}
