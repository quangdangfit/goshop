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

	"goshop/internal/cart/dto"
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

// GetCart
// =================================================================================================

func (suite *CartHandlerTestSuite) TestGetCartSuccess() {
	ctx, writer := suite.prepareContext(nil)
	ctx.Set("userId", "userID")

	suite.mockService.On("GetCartByUserID", mock.Anything, "userID").
		Return(&model.Cart{
			ID:     "cartId1",
			UserID: "userID",
			Lines: []*model.CartLine{
				{ProductID: "productId1", Quantity: 2},
			},
		}, nil).Times(1)

	suite.handler.GetCart(ctx)

	var res response.Response
	var cartRes dto.Cart
	_ = json.Unmarshal(writer.Body.Bytes(), &res)
	utils.Copy(&cartRes, &res.Result)
	suite.Equal(http.StatusOK, writer.Code)
	suite.Equal("cartId1", cartRes.ID)
}

func (suite *CartHandlerTestSuite) TestGetCartUnauthorized() {
	ctx, writer := suite.prepareContext(nil)

	suite.handler.GetCart(ctx)

	suite.Equal(http.StatusUnauthorized, writer.Code)
}

func (suite *CartHandlerTestSuite) TestGetCartFail() {
	ctx, writer := suite.prepareContext(nil)
	ctx.Set("userId", "userID")

	suite.mockService.On("GetCartByUserID", mock.Anything, "userID").
		Return(nil, errors.New("error")).Times(1)

	suite.handler.GetCart(ctx)

	suite.Equal(http.StatusInternalServerError, writer.Code)
}

// AddProduct
// =================================================================================================

func (suite *CartHandlerTestSuite) TestAddProductSuccess() {
	req := &dto.CartLineReq{
		ProductID: "productId1",
		Quantity:  2,
	}

	ctx, writer := suite.prepareContext(req)
	ctx.Set("userId", "userID")

	suite.mockService.On("AddProduct", mock.Anything, &dto.AddProductReq{
		UserID: "userID",
		Line:   req,
	}).Return(&model.Cart{
		ID:     "cartId1",
		UserID: "userID",
		Lines: []*model.CartLine{
			{ProductID: "productId1", Quantity: 2},
		},
	}, nil).Times(1)

	suite.handler.AddProduct(ctx)

	var res response.Response
	var cartRes dto.Cart
	_ = json.Unmarshal(writer.Body.Bytes(), &res)
	utils.Copy(&cartRes, &res.Result)
	suite.Equal(http.StatusOK, writer.Code)
	suite.Equal("cartId1", cartRes.ID)
}

func (suite *CartHandlerTestSuite) TestAddProductUnauthorized() {
	ctx, writer := suite.prepareContext(nil)

	suite.handler.AddProduct(ctx)

	suite.Equal(http.StatusUnauthorized, writer.Code)
}

func (suite *CartHandlerTestSuite) TestAddProductInvalidBody() {
	req := map[string]any{
		"product_id": 123,
		"quantity":   "invalid",
	}

	ctx, writer := suite.prepareContext(req)
	ctx.Set("userId", "userID")

	suite.handler.AddProduct(ctx)

	suite.Equal(http.StatusBadRequest, writer.Code)
}

func (suite *CartHandlerTestSuite) TestAddProductFail() {
	req := &dto.CartLineReq{
		ProductID: "productId1",
		Quantity:  2,
	}

	ctx, writer := suite.prepareContext(req)
	ctx.Set("userId", "userID")

	suite.mockService.On("AddProduct", mock.Anything, &dto.AddProductReq{
		UserID: "userID",
		Line:   req,
	}).Return(nil, errors.New("error")).Times(1)

	suite.handler.AddProduct(ctx)

	suite.Equal(http.StatusInternalServerError, writer.Code)
}

// RemoveProduct
// =================================================================================================

func (suite *CartHandlerTestSuite) TestRemoveProductSuccess() {
	ctx, writer := suite.prepareContext(nil)
	ctx.Set("userId", "userID")
	ctx.AddParam("productId", "productId1")

	suite.mockService.On("RemoveProduct", mock.Anything, &dto.RemoveProductReq{
		UserID:    "userID",
		ProductID: "productId1",
	}).Return(&model.Cart{
		ID:     "cartId1",
		UserID: "userID",
		Lines:  []*model.CartLine{},
	}, nil).Times(1)

	suite.handler.RemoveProduct(ctx)

	var res response.Response
	var cartRes dto.Cart
	_ = json.Unmarshal(writer.Body.Bytes(), &res)
	utils.Copy(&cartRes, &res.Result)
	suite.Equal(http.StatusOK, writer.Code)
	suite.Equal("cartId1", cartRes.ID)
}

func (suite *CartHandlerTestSuite) TestRemoveProductUnauthorized() {
	ctx, writer := suite.prepareContext(nil)
	ctx.AddParam("productId", "productId1")

	suite.handler.RemoveProduct(ctx)

	suite.Equal(http.StatusUnauthorized, writer.Code)
}

func (suite *CartHandlerTestSuite) TestRemoveProductMissProductID() {
	ctx, writer := suite.prepareContext(nil)
	ctx.Set("userId", "userID")

	suite.handler.RemoveProduct(ctx)

	suite.Equal(http.StatusBadRequest, writer.Code)
}

func (suite *CartHandlerTestSuite) TestRemoveProductFail() {
	ctx, writer := suite.prepareContext(nil)
	ctx.Set("userId", "userID")
	ctx.AddParam("productId", "productId1")

	suite.mockService.On("RemoveProduct", mock.Anything, &dto.RemoveProductReq{
		UserID:    "userID",
		ProductID: "productId1",
	}).Return(nil, errors.New("error")).Times(1)

	suite.handler.RemoveProduct(ctx)

	suite.Equal(http.StatusInternalServerError, writer.Code)
}
