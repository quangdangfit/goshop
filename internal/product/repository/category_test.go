package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/quangdangfit/gocommon/logger"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"goshop/internal/product/model"
	"goshop/pkg/config"
	dbsMocks "goshop/pkg/dbs/mocks"
)

type CategoryRepositoryTestSuite struct {
	suite.Suite
	mockDB *dbsMocks.Database
	repo   CategoryRepository
}

func (suite *CategoryRepositoryTestSuite) SetupTest() {
	logger.Initialize(config.ProductionEnv)
	suite.mockDB = dbsMocks.NewDatabase(suite.T())
	suite.repo = NewCategoryRepository(suite.mockDB)
}

func TestCategoryRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(CategoryRepositoryTestSuite))
}

// List
// =================================================================

func (suite *CategoryRepositoryTestSuite) TestListSuccess() {
	suite.mockDB.On("Find", mock.Anything, mock.Anything).Return(nil).Times(1)

	categories, err := suite.repo.List(context.Background())
	suite.Nil(err)
	suite.Equal(0, len(categories))
}

func (suite *CategoryRepositoryTestSuite) TestListFail() {
	suite.mockDB.On("Find", mock.Anything, mock.Anything).Return(errors.New("db error")).Times(1)

	categories, err := suite.repo.List(context.Background())
	suite.NotNil(err)
	suite.Nil(categories)
}

// GetByID
// =================================================================

func (suite *CategoryRepositoryTestSuite) TestGetByIDSuccess() {
	suite.mockDB.On("FindById", mock.Anything, "cat1", &model.Category{}).Return(nil).Times(1)

	category, err := suite.repo.GetByID(context.Background(), "cat1")
	suite.Nil(err)
	suite.NotNil(category)
}

func (suite *CategoryRepositoryTestSuite) TestGetByIDFail() {
	suite.mockDB.On("FindById", mock.Anything, "notfound", &model.Category{}).Return(errors.New("not found")).Times(1)

	category, err := suite.repo.GetByID(context.Background(), "notfound")
	suite.NotNil(err)
	suite.Nil(category)
}

// Create
// =================================================================

func (suite *CategoryRepositoryTestSuite) TestCreateSuccess() {
	category := &model.Category{Name: "Electronics", Slug: "electronics"}
	suite.mockDB.On("Create", mock.Anything, category).Return(nil).Times(1)

	err := suite.repo.Create(context.Background(), category)
	suite.Nil(err)
}

func (suite *CategoryRepositoryTestSuite) TestCreateFail() {
	category := &model.Category{Name: "Electronics", Slug: "electronics"}
	suite.mockDB.On("Create", mock.Anything, category).Return(errors.New("db error")).Times(1)

	err := suite.repo.Create(context.Background(), category)
	suite.NotNil(err)
}

// Update
// =================================================================

func (suite *CategoryRepositoryTestSuite) TestUpdateSuccess() {
	category := &model.Category{ID: "cat1", Name: "Updated Electronics"}
	suite.mockDB.On("Update", mock.Anything, category).Return(nil).Times(1)

	err := suite.repo.Update(context.Background(), category)
	suite.Nil(err)
}

func (suite *CategoryRepositoryTestSuite) TestUpdateFail() {
	category := &model.Category{ID: "cat1", Name: "Updated Electronics"}
	suite.mockDB.On("Update", mock.Anything, category).Return(errors.New("db error")).Times(1)

	err := suite.repo.Update(context.Background(), category)
	suite.NotNil(err)
}

// Delete
// =================================================================

func (suite *CategoryRepositoryTestSuite) TestDeleteSuccess() {
	suite.mockDB.On("Delete", mock.Anything, mock.Anything, mock.Anything).Return(nil).Times(1)

	err := suite.repo.Delete(context.Background(), "cat1")
	suite.Nil(err)
}

func (suite *CategoryRepositoryTestSuite) TestDeleteFail() {
	suite.mockDB.On("Delete", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("db error")).Times(1)

	err := suite.repo.Delete(context.Background(), "cat1")
	suite.NotNil(err)
}
