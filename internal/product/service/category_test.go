package service

import (
	"context"
	"errors"
	"testing"

	"github.com/quangdangfit/gocommon/logger"
	"github.com/quangdangfit/gocommon/validation"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"goshop/internal/product/dto"
	"goshop/internal/product/model"
	"goshop/internal/product/repository/mocks"
	"goshop/pkg/config"
)

type CategoryServiceTestSuite struct {
	suite.Suite
	mockRepo *mocks.CategoryRepository
	service  CategoryService
}

func (suite *CategoryServiceTestSuite) SetupTest() {
	logger.Initialize(config.ProductionEnv)
	validator := validation.New()
	suite.mockRepo = mocks.NewCategoryRepository(suite.T())
	suite.service = NewCategoryService(validator, suite.mockRepo)
}

func TestCategoryServiceTestSuite(t *testing.T) {
	suite.Run(t, new(CategoryServiceTestSuite))
}

// ListCategories
// =================================================================================================

func (suite *CategoryServiceTestSuite) TestListCategoriesSuccess() {
	suite.mockRepo.On("List", mock.Anything).
		Return([]*model.Category{
			{ID: "cat1", Name: "Electronics"},
			{ID: "cat2", Name: "Clothing"},
		}, nil).Times(1)

	categories, err := suite.service.ListCategories(context.Background())
	suite.Nil(err)
	suite.Equal(2, len(categories))
}

func (suite *CategoryServiceTestSuite) TestListCategoriesFail() {
	suite.mockRepo.On("List", mock.Anything).
		Return(nil, errors.New("db error")).Times(1)

	categories, err := suite.service.ListCategories(context.Background())
	suite.NotNil(err)
	suite.Nil(categories)
}

// GetCategoryByID
// =================================================================================================

func (suite *CategoryServiceTestSuite) TestGetCategoryByIDSuccess() {
	suite.mockRepo.On("GetByID", mock.Anything, "cat1").
		Return(&model.Category{ID: "cat1", Name: "Electronics"}, nil).Times(1)

	category, err := suite.service.GetCategoryByID(context.Background(), "cat1")
	suite.Nil(err)
	suite.Equal("cat1", category.ID)
}

func (suite *CategoryServiceTestSuite) TestGetCategoryByIDFail() {
	suite.mockRepo.On("GetByID", mock.Anything, "notfound").
		Return(nil, errors.New("not found")).Times(1)

	category, err := suite.service.GetCategoryByID(context.Background(), "notfound")
	suite.NotNil(err)
	suite.Nil(category)
}

// Create
// =================================================================================================

func (suite *CategoryServiceTestSuite) TestCreateSuccess() {
	req := &dto.CreateCategoryReq{
		Name: "Electronics",
		Slug: "electronics",
	}
	suite.mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil).Times(1)

	category, err := suite.service.Create(context.Background(), req)
	suite.Nil(err)
	suite.NotNil(category)
	suite.Equal("Electronics", category.Name)
}

func (suite *CategoryServiceTestSuite) TestCreateValidationFail() {
	req := &dto.CreateCategoryReq{} // missing required fields

	category, err := suite.service.Create(context.Background(), req)
	suite.NotNil(err)
	suite.Nil(category)
}

func (suite *CategoryServiceTestSuite) TestCreateDBFail() {
	req := &dto.CreateCategoryReq{
		Name: "Electronics",
		Slug: "electronics",
	}
	suite.mockRepo.On("Create", mock.Anything, mock.Anything).Return(errors.New("db error")).Times(1)

	category, err := suite.service.Create(context.Background(), req)
	suite.NotNil(err)
	suite.Nil(category)
}

// Update
// =================================================================================================

func (suite *CategoryServiceTestSuite) TestUpdateSuccess() {
	req := &dto.UpdateCategoryReq{Name: "Updated Electronics"}

	suite.mockRepo.On("GetByID", mock.Anything, "cat1").
		Return(&model.Category{ID: "cat1", Name: "Electronics"}, nil).Times(1)
	suite.mockRepo.On("Update", mock.Anything, mock.Anything).Return(nil).Times(1)

	category, err := suite.service.Update(context.Background(), "cat1", req)
	suite.Nil(err)
	suite.NotNil(category)
	suite.Equal("Updated Electronics", category.Name)
}

func (suite *CategoryServiceTestSuite) TestUpdateGetByIDFail() {
	req := &dto.UpdateCategoryReq{Name: "Updated"}

	suite.mockRepo.On("GetByID", mock.Anything, "notfound").
		Return(nil, errors.New("not found")).Times(1)

	category, err := suite.service.Update(context.Background(), "notfound", req)
	suite.NotNil(err)
	suite.Nil(category)
}

func (suite *CategoryServiceTestSuite) TestUpdateDBFail() {
	req := &dto.UpdateCategoryReq{Name: "Updated Electronics"}

	suite.mockRepo.On("GetByID", mock.Anything, "cat1").
		Return(&model.Category{ID: "cat1", Name: "Electronics"}, nil).Times(1)
	suite.mockRepo.On("Update", mock.Anything, mock.Anything).Return(errors.New("db error")).Times(1)

	category, err := suite.service.Update(context.Background(), "cat1", req)
	suite.NotNil(err)
	suite.Nil(category)
}

// Delete
// =================================================================================================

func (suite *CategoryServiceTestSuite) TestDeleteSuccess() {
	suite.mockRepo.On("Delete", mock.Anything, "cat1").Return(nil).Times(1)

	err := suite.service.Delete(context.Background(), "cat1")
	suite.Nil(err)
}

func (suite *CategoryServiceTestSuite) TestDeleteFail() {
	suite.mockRepo.On("Delete", mock.Anything, "cat1").Return(errors.New("not found")).Times(1)

	err := suite.service.Delete(context.Background(), "cat1")
	suite.NotNil(err)
}
