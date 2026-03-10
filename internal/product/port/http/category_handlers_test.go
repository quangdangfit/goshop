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

func (suite *CategoryHandlerTestSuite) TestListCategoriesSuccess() {
	ctx, writer := suite.prepareContext("GET", "/api/v1/categories", nil)

	suite.mockService.On("ListCategories", mock.Anything).
		Return([]*model.Category{
			{ID: "c1", Name: "Electronics", Slug: "electronics"},
			{ID: "c2", Name: "Books", Slug: "books"},
		}, nil).Times(1)

	suite.handler.ListCategories(ctx)

	var res response.Response
	var categories []*dto.Category
	_ = json.Unmarshal(writer.Body.Bytes(), &res)
	utils.Copy(&categories, &res.Result)
	suite.Equal(http.StatusOK, writer.Code)
	suite.Equal(2, len(categories))
	suite.Equal("c1", categories[0].ID)
	suite.Equal("Electronics", categories[0].Name)
}

func (suite *CategoryHandlerTestSuite) TestListCategoriesFail() {
	ctx, writer := suite.prepareContext("GET", "/api/v1/categories", nil)

	suite.mockService.On("ListCategories", mock.Anything).
		Return(nil, errors.New("db error")).Times(1)

	suite.handler.ListCategories(ctx)

	suite.Equal(http.StatusInternalServerError, writer.Code)
}

// GetCategoryByID
// =================================================================================================

func (suite *CategoryHandlerTestSuite) TestGetCategoryByIDSuccess() {
	ctx, writer := suite.prepareContext("GET", "/api/v1/categories/c1", nil)
	ctx.Params = gin.Params{{Key: "id", Value: "c1"}}

	suite.mockService.On("GetCategoryByID", mock.Anything, "c1").
		Return(&model.Category{ID: "c1", Name: "Electronics", Slug: "electronics"}, nil).Times(1)

	suite.handler.GetCategoryByID(ctx)

	var res response.Response
	var category dto.Category
	_ = json.Unmarshal(writer.Body.Bytes(), &res)
	utils.Copy(&category, &res.Result)
	suite.Equal(http.StatusOK, writer.Code)
	suite.Equal("c1", category.ID)
	suite.Equal("Electronics", category.Name)
}

func (suite *CategoryHandlerTestSuite) TestGetCategoryByIDMissingID() {
	ctx, writer := suite.prepareContext("GET", "/api/v1/categories/", nil)
	// id param is empty string

	suite.handler.GetCategoryByID(ctx)

	suite.Equal(http.StatusBadRequest, writer.Code)
}

func (suite *CategoryHandlerTestSuite) TestGetCategoryByIDNotFound() {
	ctx, writer := suite.prepareContext("GET", "/api/v1/categories/notfound", nil)
	ctx.Params = gin.Params{{Key: "id", Value: "notfound"}}

	suite.mockService.On("GetCategoryByID", mock.Anything, "notfound").
		Return(nil, errors.New("not found")).Times(1)

	suite.handler.GetCategoryByID(ctx)

	suite.Equal(http.StatusNotFound, writer.Code)
}

// CreateCategory
// =================================================================================================

func (suite *CategoryHandlerTestSuite) TestCreateCategorySuccess() {
	req := &dto.CreateCategoryReq{
		Name:        "Electronics",
		Slug:        "electronics",
		Description: "Electronic devices",
	}
	ctx, writer := suite.prepareContext("POST", "/api/v1/categories", req)

	suite.mockService.On("Create", mock.Anything, req).
		Return(&model.Category{
			ID:          "c1",
			Name:        "Electronics",
			Slug:        "electronics",
			Description: "Electronic devices",
		}, nil).Times(1)

	suite.handler.CreateCategory(ctx)

	var res response.Response
	var category dto.Category
	_ = json.Unmarshal(writer.Body.Bytes(), &res)
	utils.Copy(&category, &res.Result)
	suite.Equal(http.StatusOK, writer.Code)
	suite.Equal("c1", category.ID)
	suite.Equal("Electronics", category.Name)
	suite.Equal("electronics", category.Slug)
}

func (suite *CategoryHandlerTestSuite) TestCreateCategoryInvalidBody() {
	req := map[string]any{
		"name": 123,
		"slug": "electronics",
	}
	ctx, writer := suite.prepareContext("POST", "/api/v1/categories", req)

	suite.handler.CreateCategory(ctx)

	suite.Equal(http.StatusBadRequest, writer.Code)
}

func (suite *CategoryHandlerTestSuite) TestCreateCategoryFail() {
	req := &dto.CreateCategoryReq{
		Name: "Electronics",
		Slug: "electronics",
	}
	ctx, writer := suite.prepareContext("POST", "/api/v1/categories", req)

	suite.mockService.On("Create", mock.Anything, req).
		Return(nil, errors.New("duplicate key")).Times(1)

	suite.handler.CreateCategory(ctx)

	suite.Equal(http.StatusInternalServerError, writer.Code)
}

// UpdateCategory
// =================================================================================================

func (suite *CategoryHandlerTestSuite) TestUpdateCategorySuccess() {
	req := &dto.UpdateCategoryReq{
		Name:        "Updated Electronics",
		Description: "Updated description",
	}
	ctx, writer := suite.prepareContext("PUT", "/api/v1/categories/c1", req)
	ctx.Params = gin.Params{{Key: "id", Value: "c1"}}

	suite.mockService.On("Update", mock.Anything, "c1", req).
		Return(&model.Category{
			ID:          "c1",
			Name:        "Updated Electronics",
			Slug:        "electronics",
			Description: "Updated description",
		}, nil).Times(1)

	suite.handler.UpdateCategory(ctx)

	var res response.Response
	var category dto.Category
	_ = json.Unmarshal(writer.Body.Bytes(), &res)
	utils.Copy(&category, &res.Result)
	suite.Equal(http.StatusOK, writer.Code)
	suite.Equal("c1", category.ID)
	suite.Equal("Updated Electronics", category.Name)
}

func (suite *CategoryHandlerTestSuite) TestUpdateCategoryInvalidBody() {
	req := map[string]any{
		"name": 123,
	}
	ctx, writer := suite.prepareContext("PUT", "/api/v1/categories/c1", req)
	ctx.Params = gin.Params{{Key: "id", Value: "c1"}}

	suite.handler.UpdateCategory(ctx)

	suite.Equal(http.StatusBadRequest, writer.Code)
}

func (suite *CategoryHandlerTestSuite) TestUpdateCategoryFail() {
	req := &dto.UpdateCategoryReq{
		Name: "Updated Electronics",
	}
	ctx, writer := suite.prepareContext("PUT", "/api/v1/categories/c1", req)
	ctx.Params = gin.Params{{Key: "id", Value: "c1"}}

	suite.mockService.On("Update", mock.Anything, "c1", req).
		Return(nil, errors.New("not found")).Times(1)

	suite.handler.UpdateCategory(ctx)

	suite.Equal(http.StatusInternalServerError, writer.Code)
}

// DeleteCategory
// =================================================================================================

func (suite *CategoryHandlerTestSuite) TestDeleteCategorySuccess() {
	ctx, writer := suite.prepareContext("DELETE", "/api/v1/categories/c1", nil)
	ctx.Params = gin.Params{{Key: "id", Value: "c1"}}

	suite.mockService.On("Delete", mock.Anything, "c1").Return(nil).Times(1)

	suite.handler.DeleteCategory(ctx)

	suite.Equal(http.StatusOK, writer.Code)
}

func (suite *CategoryHandlerTestSuite) TestDeleteCategoryFail() {
	ctx, writer := suite.prepareContext("DELETE", "/api/v1/categories/c1", nil)
	ctx.Params = gin.Params{{Key: "id", Value: "c1"}}

	suite.mockService.On("Delete", mock.Anything, "c1").Return(errors.New("not found")).Times(1)

	suite.handler.DeleteCategory(ctx)

	suite.Equal(http.StatusInternalServerError, writer.Code)
}
