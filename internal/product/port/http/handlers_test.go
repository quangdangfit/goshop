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

	"goshop/internal/product/dto"
	"goshop/internal/product/model"
	srvMocks "goshop/internal/product/service/mocks"
	"goshop/pkg/config"
	"goshop/pkg/paging"
	redisMocks "goshop/pkg/redis/mocks"
	"goshop/pkg/response"
	"goshop/pkg/utils"
)

type ProductHandlerTestSuite struct {
	suite.Suite
	mockService *srvMocks.ProductService
	mockRedis   *redisMocks.Redis
	handler     *ProductHandler
}

func (suite *ProductHandlerTestSuite) SetupTest() {
	logger.Initialize(config.ProductionEnv)

	suite.mockService = srvMocks.NewProductService(suite.T())
	suite.mockRedis = redisMocks.NewRedis(suite.T())
	suite.handler = NewProductHandler(suite.mockRedis, suite.mockService)
}

func TestProductHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(ProductHandlerTestSuite))
}

func (suite *ProductHandlerTestSuite) prepareContext(path string, body any) (*gin.Context, *httptest.ResponseRecorder) {
	requestBody, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", path, bytes.NewBuffer(requestBody))
	c, _ := gin.CreateTestContext(w)
	c.Request = r

	return c, w
}

func (suite *ProductHandlerTestSuite) TestGetProductByID() {
	tests := []struct {
		name      string
		setup     func()
		expected  int
		checkBody func(writer *httptest.ResponseRecorder)
	}{
		{
			name: "SuccessFromDatabase",
			setup: func() {
				suite.mockRedis.On("Get", mock.Anything, &dto.Product{}).Return(errors.New("not found")).Times(1)
				suite.mockService.On("GetProductByID", mock.Anything, mock.Anything).
					Return(
						&model.Product{
							ID:          "123456",
							Name:        "product",
							Description: "description",
						},
						nil,
					).Times(1)
				suite.mockRedis.On("SetWithExpiration", mock.Anything, mock.Anything, mock.Anything).Return(nil).Times(1)
			},
			expected: http.StatusOK,
			checkBody: func(writer *httptest.ResponseRecorder) {
				var res response.Response
				var product dto.Product
				_ = json.Unmarshal(writer.Body.Bytes(), &res)
				utils.Copy(&product, &res.Result)
				suite.Equal("123456", product.ID)
				suite.Equal("product", product.Name)
				suite.Equal("description", product.Description)
			},
		},
		{
			name: "SuccessFromCache",
			setup: func() {
				suite.mockRedis.On("Get", mock.Anything, &dto.Product{}).Return(nil).Times(1)
			},
			expected: http.StatusOK,
			checkBody: func(writer *httptest.ResponseRecorder) {
				var res response.Response
				var product dto.Product
				_ = json.Unmarshal(writer.Body.Bytes(), &res)
				utils.Copy(&product, &res.Result)
				suite.Equal("", product.ID)
				suite.Equal("", product.Name)
				suite.Equal("", product.Description)
			},
		},
		{
			name: "Fail",
			setup: func() {
				suite.mockRedis.On("Get", mock.Anything, &dto.Product{}).Return(errors.New("not found")).Times(1)
				suite.mockService.On("GetProductByID", mock.Anything, mock.Anything).
					Return(nil, errors.New("error")).Times(1)
			},
			expected: http.StatusNotFound,
		},
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			ctx, writer := suite.prepareContext("/api/v1/products/123456", nil)
			tc.setup()
			suite.handler.GetProductByID(ctx)
			suite.Equal(tc.expected, writer.Code)
			if tc.checkBody != nil {
				tc.checkBody(writer)
			}
		})
	}
}

func (suite *ProductHandlerTestSuite) TestListProducts() {
	tests := []struct {
		name      string
		path      string
		setup     func()
		expected  int
		checkBody func(writer *httptest.ResponseRecorder)
	}{
		{
			name: "SuccessFromDatabase",
			path: "/api/v1/products",
			setup: func() {
				suite.mockRedis.On("Get", mock.Anything, &dto.ListProductRes{}).Return(errors.New("not found")).Times(1)
				suite.mockService.On("ListProducts", mock.Anything, mock.Anything).
					Return(
						[]*model.Product{
							{
								ID:          "123456",
								Name:        "product",
								Description: "description",
							},
						},
						&paging.Pagination{},
						nil,
					).Times(1)
				suite.mockRedis.On("SetWithExpiration", mock.Anything, mock.Anything, mock.Anything).Return(nil).Times(1)
			},
			expected: http.StatusOK,
			checkBody: func(writer *httptest.ResponseRecorder) {
				var res response.Response
				var products dto.ListProductRes
				_ = json.Unmarshal(writer.Body.Bytes(), &res)
				utils.Copy(&products, &res.Result)
				suite.Equal(1, len(products.Products))
				suite.Equal("123456", products.Products[0].ID)
				suite.Equal("product", products.Products[0].Name)
				suite.Equal("description", products.Products[0].Description)
			},
		},
		{
			name: "SuccessFromCache",
			path: "/api/v1/products",
			setup: func() {
				suite.mockRedis.On("Get", mock.Anything, &dto.ListProductRes{}).Return(nil).Times(1)
			},
			expected: http.StatusOK,
			checkBody: func(writer *httptest.ResponseRecorder) {
				var res response.Response
				var products []*dto.Product
				_ = json.Unmarshal(writer.Body.Bytes(), &res)
				utils.Copy(&products, &res.Result)
				suite.Equal(0, len(products))
			},
		},
		{
			name:     "InvalidQuery",
			path:     "/api/v1/products?page=a",
			setup:    func() {},
			expected: http.StatusBadRequest,
		},
		{
			name: "Fail",
			path: "/api/v1/products",
			setup: func() {
				suite.mockRedis.On("Get", mock.Anything, &dto.ListProductRes{}).Return(errors.New("not found")).Times(1)
				suite.mockService.On("ListProducts", mock.Anything, mock.Anything).
					Return(nil, nil, errors.New("error")).Times(1)
			},
			expected: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			ctx, writer := suite.prepareContext(tc.path, nil)
			tc.setup()
			suite.handler.ListProducts(ctx)
			suite.Equal(tc.expected, writer.Code)
			if tc.checkBody != nil {
				tc.checkBody(writer)
			}
		})
	}
}

func (suite *ProductHandlerTestSuite) TestCreateProduct() {
	tests := []struct {
		name      string
		body      any
		setup     func()
		expected  int
		checkBody func(writer *httptest.ResponseRecorder)
	}{
		{
			name: "Success",
			body: &dto.CreateProductReq{
				Name:        "product",
				Description: "description",
				Price:       10.5,
			},
			setup: func() {
				suite.mockService.On("Create", mock.Anything, &dto.CreateProductReq{
					Name:        "product",
					Description: "description",
					Price:       10.5,
				}).Return(
					&model.Product{
						Name:        "product",
						Description: "description",
						Price:       10.5,
					},
					nil,
				).Times(1)
				suite.mockRedis.On("RemovePattern", "*product*").Return(nil).Times(1)
			},
			expected: http.StatusOK,
			checkBody: func(writer *httptest.ResponseRecorder) {
				var res response.Response
				var resData dto.Product
				_ = json.Unmarshal(writer.Body.Bytes(), &res)
				utils.Copy(&resData, &res.Result)
				suite.Equal("product", resData.Name)
				suite.Equal("description", resData.Description)
				suite.Equal(10.5, resData.Price)
			},
		},
		{
			name: "InvalidPriceType",
			body: map[string]any{
				"name":        "product",
				"description": "description",
				"price":       "10.5",
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
			body: &dto.CreateProductReq{
				Name:        "product",
				Description: "description",
				Price:       10.5,
			},
			setup: func() {
				suite.mockService.On("Create", mock.Anything, &dto.CreateProductReq{
					Name:        "product",
					Description: "description",
					Price:       10.5,
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
			ctx, writer := suite.prepareContext("/api/v1/products", tc.body)
			tc.setup()
			suite.handler.CreateProduct(ctx)
			suite.Equal(tc.expected, writer.Code)
			if tc.checkBody != nil {
				tc.checkBody(writer)
			}
		})
	}
}

func (suite *ProductHandlerTestSuite) TestUpdateProduct() {
	tests := []struct {
		name      string
		body      any
		setup     func()
		expected  int
		checkBody func(writer *httptest.ResponseRecorder)
	}{
		{
			name: "Success",
			body: &dto.UpdateProductReq{
				Name:        "product",
				Description: "description",
				Price:       10.5,
			},
			setup: func() {
				suite.mockService.On("Update", mock.Anything, mock.Anything, &dto.UpdateProductReq{
					Name:        "product",
					Description: "description",
					Price:       10.5,
				}).Return(
					&model.Product{
						ID:          "123456",
						Name:        "product",
						Description: "description",
						Price:       10.5,
					},
					nil,
				).Times(1)
				suite.mockRedis.On("RemovePattern", "*product*").Return(nil).Times(1)
			},
			expected: http.StatusOK,
			checkBody: func(writer *httptest.ResponseRecorder) {
				var res response.Response
				var resData dto.Product
				_ = json.Unmarshal(writer.Body.Bytes(), &res)
				utils.Copy(&resData, &res.Result)
				suite.Equal("123456", resData.ID)
				suite.Equal("product", resData.Name)
				suite.Equal("description", resData.Description)
				suite.Equal(10.5, resData.Price)
			},
		},
		{
			name: "InvalidPriceType",
			body: map[string]any{
				"name":        "product",
				"description": "description",
				"price":       "10.5",
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
			body: &dto.UpdateProductReq{
				Name:        "product",
				Description: "description",
				Price:       10.5,
			},
			setup: func() {
				suite.mockService.On("Update", mock.Anything, mock.Anything, &dto.UpdateProductReq{
					Name:        "product",
					Description: "description",
					Price:       10.5,
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
			ctx, writer := suite.prepareContext("/api/v1/products/123456", tc.body)
			tc.setup()
			suite.handler.UpdateProduct(ctx)
			suite.Equal(tc.expected, writer.Code)
			if tc.checkBody != nil {
				tc.checkBody(writer)
			}
		})
	}
}
