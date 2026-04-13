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

	"goshop/internal/cart/domain"
	"goshop/internal/cart/model"
	"goshop/internal/cart/service/mocks"
	"goshop/pkg/config"
	"goshop/pkg/response"
	"goshop/pkg/utils"
)

type CartHandlerTestSuite struct {
	suite.Suite
	mockService *mocks.CartService
	handler     *CartHandler
}

func (suite *CartHandlerTestSuite) SetupTest() {
	logger.Initialize(config.ProductionEnv)

	suite.mockService = mocks.NewCartService(suite.T())
	suite.handler = NewCartHandler(suite.mockService)
}

func TestCartHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(CartHandlerTestSuite))
}

func (suite *CartHandlerTestSuite) prepareContext(body any) (*gin.Context, *httptest.ResponseRecorder) {
	requestBody, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("", "/", bytes.NewBuffer(requestBody))
	c, _ := gin.CreateTestContext(w)
	c.Request = r
	return c, w
}

func (suite *CartHandlerTestSuite) TestGetCart() {
	tests := []struct {
		name      string
		setup     func(ctx *gin.Context)
		expected  int
		checkBody func(writer *httptest.ResponseRecorder)
	}{
		{
			name: "Success",
			setup: func(ctx *gin.Context) {
				ctx.Set("userId", "userID")
				suite.mockService.On("GetCartByUserID", mock.Anything, "userID").
					Return(&model.Cart{
						ID:     "cartId1",
						UserID: "userID",
						Lines: []*model.CartLine{
							{ProductID: "productId1", Quantity: 2},
						},
					}, nil).Times(1)
			},
			expected: http.StatusOK,
			checkBody: func(writer *httptest.ResponseRecorder) {
				var res response.Response
				var cartRes domain.Cart
				_ = json.Unmarshal(writer.Body.Bytes(), &res)
				_ = utils.Copy(&cartRes, &res.Result)
				suite.Equal("cartId1", cartRes.ID)
			},
		},
		{
			name:     "Unauthorized",
			setup:    func(ctx *gin.Context) {},
			expected: http.StatusUnauthorized,
		},
		{
			name: "Fail",
			setup: func(ctx *gin.Context) {
				ctx.Set("userId", "userID")
				suite.mockService.On("GetCartByUserID", mock.Anything, "userID").
					Return(nil, errors.New("error")).Times(1)
			},
			expected: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			ctx, writer := suite.prepareContext(nil)
			tc.setup(ctx)
			suite.handler.GetCart(ctx)
			suite.Equal(tc.expected, writer.Code)
			if tc.checkBody != nil {
				tc.checkBody(writer)
			}
		})
	}
}

func (suite *CartHandlerTestSuite) TestAddProduct() {
	tests := []struct {
		name      string
		body      any
		setup     func(ctx *gin.Context)
		expected  int
		checkBody func(writer *httptest.ResponseRecorder)
	}{
		{
			name: "Success",
			body: &domain.CartLineReq{
				ProductID: "productId1",
				Quantity:  2,
			},
			setup: func(ctx *gin.Context) {
				ctx.Set("userId", "userID")
				suite.mockService.On("AddProduct", mock.Anything, &domain.AddProductReq{
					UserID: "userID",
					Line: &domain.CartLineReq{
						ProductID: "productId1",
						Quantity:  2,
					},
				}).Return(&model.Cart{
					ID:     "cartId1",
					UserID: "userID",
					Lines: []*model.CartLine{
						{ProductID: "productId1", Quantity: 2},
					},
				}, nil).Times(1)
			},
			expected: http.StatusOK,
			checkBody: func(writer *httptest.ResponseRecorder) {
				var res response.Response
				var cartRes domain.Cart
				_ = json.Unmarshal(writer.Body.Bytes(), &res)
				_ = utils.Copy(&cartRes, &res.Result)
				suite.Equal("cartId1", cartRes.ID)
			},
		},
		{
			name:     "Unauthorized",
			body:     nil,
			setup:    func(ctx *gin.Context) {},
			expected: http.StatusUnauthorized,
		},
		{
			name: "InvalidBody",
			body: map[string]any{
				"product_id": 123,
				"quantity":   "invalid",
			},
			setup: func(ctx *gin.Context) {
				ctx.Set("userId", "userID")
			},
			expected: http.StatusBadRequest,
		},
		{
			name: "Fail",
			body: &domain.CartLineReq{
				ProductID: "productId1",
				Quantity:  2,
			},
			setup: func(ctx *gin.Context) {
				ctx.Set("userId", "userID")
				suite.mockService.On("AddProduct", mock.Anything, &domain.AddProductReq{
					UserID: "userID",
					Line: &domain.CartLineReq{
						ProductID: "productId1",
						Quantity:  2,
					},
				}).Return(nil, errors.New("error")).Times(1)
			},
			expected: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			ctx, writer := suite.prepareContext(tc.body)
			tc.setup(ctx)
			suite.handler.AddProduct(ctx)
			suite.Equal(tc.expected, writer.Code)
			if tc.checkBody != nil {
				tc.checkBody(writer)
			}
		})
	}
}

func (suite *CartHandlerTestSuite) TestRemoveProduct() {
	tests := []struct {
		name      string
		setup     func(ctx *gin.Context)
		expected  int
		checkBody func(writer *httptest.ResponseRecorder)
	}{
		{
			name: "Success",
			setup: func(ctx *gin.Context) {
				ctx.Set("userId", "userID")
				ctx.AddParam("productId", "productId1")
				suite.mockService.On("RemoveProduct", mock.Anything, &domain.RemoveProductReq{
					UserID:    "userID",
					ProductID: "productId1",
				}).Return(&model.Cart{
					ID:     "cartId1",
					UserID: "userID",
					Lines:  []*model.CartLine{},
				}, nil).Times(1)
			},
			expected: http.StatusOK,
			checkBody: func(writer *httptest.ResponseRecorder) {
				var res response.Response
				var cartRes domain.Cart
				_ = json.Unmarshal(writer.Body.Bytes(), &res)
				_ = utils.Copy(&cartRes, &res.Result)
				suite.Equal("cartId1", cartRes.ID)
			},
		},
		{
			name: "Unauthorized",
			setup: func(ctx *gin.Context) {
				ctx.AddParam("productId", "productId1")
			},
			expected: http.StatusUnauthorized,
		},
		{
			name: "MissProductID",
			setup: func(ctx *gin.Context) {
				ctx.Set("userId", "userID")
			},
			expected: http.StatusBadRequest,
		},
		{
			name: "Fail",
			setup: func(ctx *gin.Context) {
				ctx.Set("userId", "userID")
				ctx.AddParam("productId", "productId1")
				suite.mockService.On("RemoveProduct", mock.Anything, &domain.RemoveProductReq{
					UserID:    "userID",
					ProductID: "productId1",
				}).Return(nil, errors.New("error")).Times(1)
			},
			expected: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			ctx, writer := suite.prepareContext(nil)
			tc.setup(ctx)
			suite.handler.RemoveProduct(ctx)
			suite.Equal(tc.expected, writer.Code)
			if tc.checkBody != nil {
				tc.checkBody(writer)
			}
		})
	}
}
