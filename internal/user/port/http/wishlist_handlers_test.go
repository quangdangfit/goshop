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

func (suite *WishlistHandlerTestSuite) TestGetWishlistSuccess() {
	ctx, writer := suite.prepareContext("GET", "/api/v1/wishlist", nil)
	ctx.Set("userId", "u1")

	suite.mockService.On("GetWishlist", mock.Anything, "u1").
		Return([]*model.Wishlist{
			{UserID: "u1", ProductID: "p1"},
			{UserID: "u1", ProductID: "p2"},
		}, nil).Times(1)

	suite.handler.GetWishlist(ctx)

	var res response.Response
	var items []*dto.WishlistItem
	_ = json.Unmarshal(writer.Body.Bytes(), &res)
	utils.Copy(&items, &res.Result)
	suite.Equal(http.StatusOK, writer.Code)
	suite.Equal(2, len(items))
}

func (suite *WishlistHandlerTestSuite) TestGetWishlistUnauthorized() {
	ctx, writer := suite.prepareContext("GET", "/api/v1/wishlist", nil)

	suite.handler.GetWishlist(ctx)

	suite.Equal(http.StatusUnauthorized, writer.Code)
}

func (suite *WishlistHandlerTestSuite) TestGetWishlistFail() {
	ctx, writer := suite.prepareContext("GET", "/api/v1/wishlist", nil)
	ctx.Set("userId", "u1")

	suite.mockService.On("GetWishlist", mock.Anything, "u1").
		Return(nil, errors.New("db error")).Times(1)

	suite.handler.GetWishlist(ctx)

	suite.Equal(http.StatusInternalServerError, writer.Code)
}

// AddProduct
// =================================================================================================

func (suite *WishlistHandlerTestSuite) TestAddProductSuccess() {
	req := &dto.AddToWishlistReq{ProductID: "p1"}
	ctx, writer := suite.prepareContext("POST", "/api/v1/wishlist", req)
	ctx.Set("userId", "u1")

	suite.mockService.On("AddProduct", mock.Anything, "u1", req).Return(nil).Times(1)

	suite.handler.AddProduct(ctx)

	suite.Equal(http.StatusOK, writer.Code)
}

func (suite *WishlistHandlerTestSuite) TestAddProductUnauthorized() {
	req := &dto.AddToWishlistReq{ProductID: "p1"}
	ctx, writer := suite.prepareContext("POST", "/api/v1/wishlist", req)

	suite.handler.AddProduct(ctx)

	suite.Equal(http.StatusUnauthorized, writer.Code)
}

func (suite *WishlistHandlerTestSuite) TestAddProductInvalidBody() {
	req := map[string]any{"product_id": 123}
	ctx, writer := suite.prepareContext("POST", "/api/v1/wishlist", req)
	ctx.Set("userId", "u1")

	suite.handler.AddProduct(ctx)

	suite.Equal(http.StatusBadRequest, writer.Code)
}

func (suite *WishlistHandlerTestSuite) TestAddProductFail() {
	req := &dto.AddToWishlistReq{ProductID: "p1"}
	ctx, writer := suite.prepareContext("POST", "/api/v1/wishlist", req)
	ctx.Set("userId", "u1")

	suite.mockService.On("AddProduct", mock.Anything, "u1", req).
		Return(errors.New("already exists")).Times(1)

	suite.handler.AddProduct(ctx)

	suite.Equal(http.StatusInternalServerError, writer.Code)
}

// RemoveProduct
// =================================================================================================

func (suite *WishlistHandlerTestSuite) TestRemoveProductSuccess() {
	ctx, writer := suite.prepareContext("DELETE", "/api/v1/wishlist/p1", nil)
	ctx.Set("userId", "u1")
	ctx.Params = gin.Params{{Key: "productId", Value: "p1"}}

	suite.mockService.On("RemoveProduct", mock.Anything, "u1", "p1").Return(nil).Times(1)

	suite.handler.RemoveProduct(ctx)

	suite.Equal(http.StatusOK, writer.Code)
}

func (suite *WishlistHandlerTestSuite) TestRemoveProductUnauthorized() {
	ctx, writer := suite.prepareContext("DELETE", "/api/v1/wishlist/p1", nil)

	suite.handler.RemoveProduct(ctx)

	suite.Equal(http.StatusUnauthorized, writer.Code)
}

func (suite *WishlistHandlerTestSuite) TestRemoveProductMissingID() {
	ctx, writer := suite.prepareContext("DELETE", "/api/v1/wishlist/", nil)
	ctx.Set("userId", "u1")
	// productId param is empty

	suite.handler.RemoveProduct(ctx)

	suite.Equal(http.StatusBadRequest, writer.Code)
}

func (suite *WishlistHandlerTestSuite) TestRemoveProductFail() {
	ctx, writer := suite.prepareContext("DELETE", "/api/v1/wishlist/p1", nil)
	ctx.Set("userId", "u1")
	ctx.Params = gin.Params{{Key: "productId", Value: "p1"}}

	suite.mockService.On("RemoveProduct", mock.Anything, "u1", "p1").
		Return(errors.New("not found")).Times(1)

	suite.handler.RemoveProduct(ctx)

	suite.Equal(http.StatusInternalServerError, writer.Code)
}
