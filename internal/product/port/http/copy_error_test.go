package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"math"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/quangdangfit/gocommon/logger"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"goshop/internal/product/domain"
	"goshop/internal/product/model"
	"goshop/pkg/config"
	"goshop/pkg/paging"
)

var errCacheMiss = errors.New("cache miss")

// nanProduct returns a Product whose Price is NaN; json.Marshal rejects NaN, which
// is the only practical way to drive utils.Copy into its error branch in handler tests.
func nanProduct() *model.Product {
	return &model.Product{ID: "p1", Price: math.NaN()}
}

func newProductHandlerSuite(t *testing.T) *ProductHandlerTestSuite {
	logger.Initialize(config.ProductionEnv)
	gin.SetMode(gin.TestMode)
	s := &ProductHandlerTestSuite{}
	s.SetT(t)
	s.SetupTest()
	return s
}

func TestGetProductByID_CopyError(t *testing.T) {
	s := newProductHandlerSuite(t)
	s.mockRedis.On("Get", mock.Anything, mock.Anything).Return(errCacheMiss).Once()
	s.mockService.On("GetProductByID", mock.Anything, "p1").Return(nanProduct(), nil).Once()

	ctx, w := s.prepareContext("/api/v1/products/p1", nil)
	ctx.Params = gin.Params{{Key: "id", Value: "p1"}}
	s.handler.GetProductByID(ctx)
	require.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestCreateProduct_CopyError(t *testing.T) {
	s := newProductHandlerSuite(t)
	s.mockService.On("Create", mock.Anything, mock.Anything).Return(nanProduct(), nil).Once()
	s.mockRedis.On("RemovePattern", mock.Anything).Return(nil).Maybe()

	body, _ := json.Marshal(domain.CreateProductReq{
		Name: "x", Description: "x", Price: 1, StockQuantity: 1,
	})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))
	s.handler.CreateProduct(c)
	require.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestUpdateProduct_CopyError(t *testing.T) {
	s := newProductHandlerSuite(t)
	s.mockService.On("Update", mock.Anything, "p1", mock.Anything).Return(nanProduct(), nil).Once()
	s.mockRedis.On("RemovePattern", mock.Anything).Return(nil).Maybe()

	body, _ := json.Marshal(domain.UpdateProductReq{Name: "x"})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "p1"}}
	c.Request = httptest.NewRequest(http.MethodPut, "/", bytes.NewBuffer(body))
	s.handler.UpdateProduct(c)
	require.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestAddStock_CopyError(t *testing.T) {
	s := newProductHandlerSuite(t)
	s.mockService.On("AddStock", mock.Anything, "p1", 1, mock.Anything).Return(nanProduct(), nil).Once()
	s.mockRedis.On("RemovePattern", mock.Anything).Return(nil).Maybe()

	body, _ := json.Marshal(map[string]any{"quantity": 1})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("userId", "admin")
	c.Params = gin.Params{{Key: "id", Value: "p1"}}
	c.Request = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))
	s.handler.AddStock(c)
	require.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestListProducts_CopyError(t *testing.T) {
	s := newProductHandlerSuite(t)
	s.mockRedis.On("Get", mock.Anything, mock.Anything).Return(errCacheMiss).Once()
	s.mockService.On("ListProducts", mock.Anything, mock.Anything).
		Return([]*model.Product{nanProduct()}, &paging.Pagination{}, nil).Once()

	ctx, w := s.prepareContext("/api/v1/products", nil)
	s.handler.ListProducts(ctx)
	require.Equal(t, http.StatusInternalServerError, w.Code)
}
