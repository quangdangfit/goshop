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
	"goshop/pkg/response"
	"goshop/pkg/utils"
)

type CategoryHandlerTestSuite struct {
	suite.Suite
	mockService *srvMocks.CategoryService
	handler     *CategoryHandler
}

func (suite *CategoryHandlerTestSuite) SetupTest() {
	logger.Initialize(config.ProductionEnv)
	suite.mockService = srvMocks.NewCategoryService(suite.T())
	suite.handler = NewCategoryHandler(suite.mockService)
}

func TestCategoryHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(CategoryHandlerTestSuite))
}

func (suite *CategoryHandlerTestSuite) prepareContext(method, path string, body any) (*gin.Context, *httptest.ResponseRecorder) {
	requestBody, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(method, path, bytes.NewBuffer(requestBody))
	c, _ := gin.CreateTestContext(w)
	c.Request = r
	return c, w
}

// ListCategories
// =================================================================================================

func (suite *CategoryHandlerTestSuite) TestListCategories() {
	tests := []struct {
		name      string
		body      any
		params    gin.Params
		setup     func()
		expected  int
		checkBody func(writer *httptest.ResponseRecorder)
	}{
		{
			name: "Success",
			setup: func() {
				suite.mockService.On("ListCategories", mock.Anything).
					Return([]*model.Category{
						{ID: "c1", Name: "Electronics", Slug: "electronics"},
						{ID: "c2", Name: "Books", Slug: "books"},
					}, nil).Times(1)
			},
			expected: http.StatusOK,
			checkBody: func(writer *httptest.ResponseRecorder) {
				var res response.Response
				var categories []*dto.Category
				_ = json.Unmarshal(writer.Body.Bytes(), &res)
				utils.Copy(&categories, &res.Result)
				suite.Equal(2, len(categories))
				suite.Equal("c1", categories[0].ID)
				suite.Equal("Electronics", categories[0].Name)
			},
		},
		{
			name: "Fail",
			setup: func() {
				suite.mockService.On("ListCategories", mock.Anything).
					Return(nil, errors.New("db error")).Times(1)
			},
			expected: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			ctx, writer := suite.prepareContext("GET", "/api/v1/categories", tc.body)
			if tc.params != nil {
				ctx.Params = tc.params
			}
			tc.setup()
			suite.handler.ListCategories(ctx)
			suite.Equal(tc.expected, writer.Code)
			if tc.checkBody != nil {
				tc.checkBody(writer)
			}
		})
	}
}

// GetCategoryByID
// =================================================================================================

func (suite *CategoryHandlerTestSuite) TestGetCategoryByID() {
	tests := []struct {
		name      string
		body      any
		params    gin.Params
		setup     func()
		expected  int
		checkBody func(writer *httptest.ResponseRecorder)
	}{
		{
			name:   "Success",
			params: gin.Params{{Key: "id", Value: "c1"}},
			setup: func() {
				suite.mockService.On("GetCategoryByID", mock.Anything, "c1").
					Return(&model.Category{ID: "c1", Name: "Electronics", Slug: "electronics"}, nil).Times(1)
			},
			expected: http.StatusOK,
			checkBody: func(writer *httptest.ResponseRecorder) {
				var res response.Response
				var category dto.Category
				_ = json.Unmarshal(writer.Body.Bytes(), &res)
				utils.Copy(&category, &res.Result)
				suite.Equal("c1", category.ID)
				suite.Equal("Electronics", category.Name)
			},
		},
		{
			name:     "MissingID",
			setup:    func() {},
			expected: http.StatusBadRequest,
		},
		{
			name:   "NotFound",
			params: gin.Params{{Key: "id", Value: "notfound"}},
			setup: func() {
				suite.mockService.On("GetCategoryByID", mock.Anything, "notfound").
					Return(nil, errors.New("not found")).Times(1)
			},
			expected: http.StatusNotFound,
		},
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			ctx, writer := suite.prepareContext("GET", "/api/v1/categories/c1", tc.body)
			if tc.params != nil {
				ctx.Params = tc.params
			}
			tc.setup()
			suite.handler.GetCategoryByID(ctx)
			suite.Equal(tc.expected, writer.Code)
			if tc.checkBody != nil {
				tc.checkBody(writer)
			}
		})
	}
}

// CreateCategory
// =================================================================================================

func (suite *CategoryHandlerTestSuite) TestCreateCategory() {
	tests := []struct {
		name      string
		body      any
		params    gin.Params
		setup     func()
		expected  int
		checkBody func(writer *httptest.ResponseRecorder)
	}{
		{
			name: "Success",
			body: &dto.CreateCategoryReq{
				Name:        "Electronics",
				Slug:        "electronics",
				Description: "Electronic devices",
			},
			setup: func() {
				suite.mockService.On("Create", mock.Anything, &dto.CreateCategoryReq{
					Name:        "Electronics",
					Slug:        "electronics",
					Description: "Electronic devices",
				}).
					Return(&model.Category{
						ID:          "c1",
						Name:        "Electronics",
						Slug:        "electronics",
						Description: "Electronic devices",
					}, nil).Times(1)
			},
			expected: http.StatusOK,
			checkBody: func(writer *httptest.ResponseRecorder) {
				var res response.Response
				var category dto.Category
				_ = json.Unmarshal(writer.Body.Bytes(), &res)
				utils.Copy(&category, &res.Result)
				suite.Equal("c1", category.ID)
				suite.Equal("Electronics", category.Name)
				suite.Equal("electronics", category.Slug)
			},
		},
		{
			name: "InvalidBody",
			body: map[string]any{
				"name": 123,
				"slug": "electronics",
			},
			setup:    func() {},
			expected: http.StatusBadRequest,
		},
		{
			name: "Fail",
			body: &dto.CreateCategoryReq{
				Name: "Electronics",
				Slug: "electronics",
			},
			setup: func() {
				suite.mockService.On("Create", mock.Anything, &dto.CreateCategoryReq{
					Name: "Electronics",
					Slug: "electronics",
				}).
					Return(nil, errors.New("duplicate key")).Times(1)
			},
			expected: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			ctx, writer := suite.prepareContext("POST", "/api/v1/categories", tc.body)
			if tc.params != nil {
				ctx.Params = tc.params
			}
			tc.setup()
			suite.handler.CreateCategory(ctx)
			suite.Equal(tc.expected, writer.Code)
			if tc.checkBody != nil {
				tc.checkBody(writer)
			}
		})
	}
}

// UpdateCategory
// =================================================================================================

func (suite *CategoryHandlerTestSuite) TestUpdateCategory() {
	tests := []struct {
		name      string
		body      any
		params    gin.Params
		setup     func()
		expected  int
		checkBody func(writer *httptest.ResponseRecorder)
	}{
		{
			name: "Success",
			body: &dto.UpdateCategoryReq{
				Name:        "Updated Electronics",
				Description: "Updated description",
			},
			params: gin.Params{{Key: "id", Value: "c1"}},
			setup: func() {
				suite.mockService.On("Update", mock.Anything, "c1", &dto.UpdateCategoryReq{
					Name:        "Updated Electronics",
					Description: "Updated description",
				}).
					Return(&model.Category{
						ID:          "c1",
						Name:        "Updated Electronics",
						Slug:        "electronics",
						Description: "Updated description",
					}, nil).Times(1)
			},
			expected: http.StatusOK,
			checkBody: func(writer *httptest.ResponseRecorder) {
				var res response.Response
				var category dto.Category
				_ = json.Unmarshal(writer.Body.Bytes(), &res)
				utils.Copy(&category, &res.Result)
				suite.Equal("c1", category.ID)
				suite.Equal("Updated Electronics", category.Name)
			},
		},
		{
			name: "InvalidBody",
			body: map[string]any{
				"name": 123,
			},
			params:   gin.Params{{Key: "id", Value: "c1"}},
			setup:    func() {},
			expected: http.StatusBadRequest,
		},
		{
			name: "Fail",
			body: &dto.UpdateCategoryReq{
				Name: "Updated Electronics",
			},
			params: gin.Params{{Key: "id", Value: "c1"}},
			setup: func() {
				suite.mockService.On("Update", mock.Anything, "c1", &dto.UpdateCategoryReq{
					Name: "Updated Electronics",
				}).
					Return(nil, errors.New("not found")).Times(1)
			},
			expected: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			ctx, writer := suite.prepareContext("PUT", "/api/v1/categories/c1", tc.body)
			if tc.params != nil {
				ctx.Params = tc.params
			}
			tc.setup()
			suite.handler.UpdateCategory(ctx)
			suite.Equal(tc.expected, writer.Code)
			if tc.checkBody != nil {
				tc.checkBody(writer)
			}
		})
	}
}

// DeleteCategory
// =================================================================================================

func (suite *CategoryHandlerTestSuite) TestDeleteCategory() {
	tests := []struct {
		name      string
		body      any
		params    gin.Params
		setup     func()
		expected  int
		checkBody func(writer *httptest.ResponseRecorder)
	}{
		{
			name:   "Success",
			params: gin.Params{{Key: "id", Value: "c1"}},
			setup: func() {
				suite.mockService.On("Delete", mock.Anything, "c1").Return(nil).Times(1)
			},
			expected: http.StatusOK,
		},
		{
			name:   "Fail",
			params: gin.Params{{Key: "id", Value: "c1"}},
			setup: func() {
				suite.mockService.On("Delete", mock.Anything, "c1").Return(errors.New("not found")).Times(1)
			},
			expected: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			ctx, writer := suite.prepareContext("DELETE", "/api/v1/categories/c1", tc.body)
			if tc.params != nil {
				ctx.Params = tc.params
			}
			tc.setup()
			suite.handler.DeleteCategory(ctx)
			suite.Equal(tc.expected, writer.Code)
			if tc.checkBody != nil {
				tc.checkBody(writer)
			}
		})
	}
}
