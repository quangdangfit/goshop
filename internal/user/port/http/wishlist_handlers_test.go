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

type WishlistHandlerTestSuite struct {
	suite.Suite
	mockService *srvMocks.WishlistService
	handler     *WishlistHandler
}

func (suite *WishlistHandlerTestSuite) SetupTest() {
	logger.Initialize(config.ProductionEnv)
	suite.mockService = srvMocks.NewWishlistService(suite.T())
	suite.handler = NewWishlistHandler(suite.mockService)
}

func TestWishlistHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(WishlistHandlerTestSuite))
}

func (suite *WishlistHandlerTestSuite) prepareContext(method, path string, body any) (*gin.Context, *httptest.ResponseRecorder) {
	requestBody, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(method, path, bytes.NewBuffer(requestBody))
	c, _ := gin.CreateTestContext(w)
	c.Request = r
	return c, w
}

// GetWishlist
// =================================================================================================

func (suite *WishlistHandlerTestSuite) TestGetWishlist() {
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
			userId: "u1",
			setup: func() {
				suite.mockService.On("GetWishlist", mock.Anything, "u1").
					Return([]*model.Wishlist{
						{UserID: "u1", ProductID: "p1"},
						{UserID: "u1", ProductID: "p2"},
					}, nil).Times(1)
			},
			expected: http.StatusOK,
			checkBody: func(writer *httptest.ResponseRecorder) {
				var res response.Response
				var items []*domain.WishlistItem
				_ = json.Unmarshal(writer.Body.Bytes(), &res)
				utils.Copy(&items, &res.Result)
				suite.Equal(2, len(items))
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
				suite.mockService.On("GetWishlist", mock.Anything, "u1").
					Return(nil, errors.New("db error")).Times(1)
			},
			expected: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			ctx, writer := suite.prepareContext("GET", "/api/v1/wishlist", tc.body)
			if tc.userId != "" {
				ctx.Set("userId", tc.userId)
			}
			tc.setup()
			suite.handler.GetWishlist(ctx)
			suite.Equal(tc.expected, writer.Code)
			if tc.checkBody != nil {
				tc.checkBody(writer)
			}
		})
	}
}

// AddProduct
// =================================================================================================

func (suite *WishlistHandlerTestSuite) TestAddProduct() {
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
			body:   &domain.AddToWishlistReq{ProductID: "p1"},
			userId: "u1",
			setup: func() {
				suite.mockService.On("AddProduct", mock.Anything, "u1", &domain.AddToWishlistReq{ProductID: "p1"}).
					Return(nil).Times(1)
			},
			expected: http.StatusOK,
		},
		{
			name:     "Unauthorized",
			body:     &domain.AddToWishlistReq{ProductID: "p1"},
			userId:   "",
			setup:    func() {},
			expected: http.StatusUnauthorized,
		},
		{
			name:     "InvalidBody",
			body:     map[string]any{"product_id": 123},
			userId:   "u1",
			setup:    func() {},
			expected: http.StatusBadRequest,
		},
		{
			name:   "Fail",
			body:   &domain.AddToWishlistReq{ProductID: "p1"},
			userId: "u1",
			setup: func() {
				suite.mockService.On("AddProduct", mock.Anything, "u1", &domain.AddToWishlistReq{ProductID: "p1"}).
					Return(errors.New("already exists")).Times(1)
			},
			expected: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			ctx, writer := suite.prepareContext("POST", "/api/v1/wishlist", tc.body)
			if tc.userId != "" {
				ctx.Set("userId", tc.userId)
			}
			tc.setup()
			suite.handler.AddProduct(ctx)
			suite.Equal(tc.expected, writer.Code)
			if tc.checkBody != nil {
				tc.checkBody(writer)
			}
		})
	}
}

// RemoveProduct
// =================================================================================================

func (suite *WishlistHandlerTestSuite) TestRemoveProduct() {
	tests := []struct {
		name      string
		body      any
		params    gin.Params
		userId    string
		setup     func()
		expected  int
		checkBody func(writer *httptest.ResponseRecorder)
	}{
		{
			name:   "Success",
			userId: "u1",
			params: gin.Params{{Key: "productId", Value: "p1"}},
			setup: func() {
				suite.mockService.On("RemoveProduct", mock.Anything, "u1", "p1").
					Return(nil).Times(1)
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
			name:     "MissingID",
			userId:   "u1",
			setup:    func() {},
			expected: http.StatusBadRequest,
		},
		{
			name:   "Fail",
			userId: "u1",
			params: gin.Params{{Key: "productId", Value: "p1"}},
			setup: func() {
				suite.mockService.On("RemoveProduct", mock.Anything, "u1", "p1").
					Return(errors.New("not found")).Times(1)
			},
			expected: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			ctx, writer := suite.prepareContext("DELETE", "/api/v1/wishlist/p1", tc.body)
			if tc.userId != "" {
				ctx.Set("userId", tc.userId)
			}
			if tc.params != nil {
				ctx.Params = tc.params
			}
			tc.setup()
			suite.handler.RemoveProduct(ctx)
			suite.Equal(tc.expected, writer.Code)
			if tc.checkBody != nil {
				tc.checkBody(writer)
			}
		})
	}
}
