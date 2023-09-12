package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/quangdangfit/gocommon/logger"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"goshop/internal/product/dto"
	"goshop/internal/product/model"
	"goshop/pkg/config"
	"goshop/pkg/dbs/mocks"
)

type ProductRepositoryTestSuite struct {
	suite.Suite
	mockDB *mocks.IDatabase
	repo   IProductRepository
}

func (suite *ProductRepositoryTestSuite) SetupTest() {
	logger.Initialize(config.ProductionEnv)

	suite.mockDB = mocks.NewIDatabase(suite.T())
	suite.repo = NewProductRepository(suite.mockDB)
}

func TestProductRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(ProductRepositoryTestSuite))
}

// Create
// =================================================================

func (suite *ProductRepositoryTestSuite) TestCreateProductSuccessfully() {
	product := &model.Product{
		Name:        "product name",
		Description: "product description",
		Price:       10.5,
	}
	suite.mockDB.On("Create", mock.Anything, product).
		Return(nil).Times(1)

	err := suite.repo.Create(context.Background(), product)
	suite.Nil(err)
}

func (suite *ProductRepositoryTestSuite) TestCreateProductFail() {
	product := &model.Product{
		Name:        "product name",
		Description: "product description",
		Price:       10.5,
	}
	suite.mockDB.On("Create", mock.Anything, product).
		Return(errors.New("error")).Times(1)

	err := suite.repo.Create(context.Background(), product)
	suite.NotNil(err)
}

// Update
// =================================================================

func (suite *ProductRepositoryTestSuite) TestUpdateProductSuccessfully() {
	product := &model.Product{
		ID:          "productId1",
		Name:        "product name",
		Description: "product description",
		Price:       10.5,
	}
	suite.mockDB.On("Update", mock.Anything, product).
		Return(nil).Times(1)

	err := suite.repo.Update(context.Background(), product)
	suite.Nil(err)
}

func (suite *ProductRepositoryTestSuite) TestUpdateProductFail() {
	product := &model.Product{
		ID:          "productId1",
		Name:        "product name",
		Description: "product description",
		Price:       10.5,
	}
	suite.mockDB.On("Update", mock.Anything, product).
		Return(errors.New("error")).Times(1)

	err := suite.repo.Update(context.Background(), product)
	suite.NotNil(err)
}

// GetProductByID
// =================================================================

func (suite *ProductRepositoryTestSuite) TestGetProductByIDSuccessfully() {
	suite.mockDB.On("FindById", mock.Anything, "productId1", &model.Product{}).
		Return(nil).Times(1)

	product, err := suite.repo.GetProductByID(context.Background(), "productId1")
	suite.Nil(err)
	suite.NotNil(product)
}

func (suite *ProductRepositoryTestSuite) TestGetProductByIDFail() {
	suite.mockDB.On("FindById", mock.Anything, "productId1", &model.Product{}).
		Return(errors.New("error")).Times(1)

	product, err := suite.repo.GetProductByID(context.Background(), "productId1")
	suite.NotNil(err)
	suite.Nil(product)
}

// ListProducts
// =================================================================

func (suite *ProductRepositoryTestSuite) TestListProductsSuccessfully() {
	req := &dto.ListProductReq{
		Name:      "name",
		Code:      "code",
		Page:      2,
		Limit:     10,
		OrderBy:   "name",
		OrderDesc: true,
	}

	suite.mockDB.On("Count", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(nil).Times(1)

	suite.mockDB.On("Find", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(nil).Times(1)

	products, pagination, err := suite.repo.ListProducts(context.Background(), req)
	suite.Nil(err)
	suite.Equal(0, len(products))
	suite.NotNil(pagination)
}

func (suite *ProductRepositoryTestSuite) TestListProductsCountFail() {
	suite.mockDB.On("Count", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(errors.New("error")).Times(1)

	product, pagination, err := suite.repo.ListProducts(context.Background(), &dto.ListProductReq{})
	suite.NotNil(err)
	suite.Nil(product)
	suite.Nil(pagination)
}

func (suite *ProductRepositoryTestSuite) TestListProductsFindFail() {
	suite.mockDB.On("Count", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(nil).Times(1)

	suite.mockDB.On("Find", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(errors.New("error")).Times(1)

	product, pagination, err := suite.repo.ListProducts(context.Background(), &dto.ListProductReq{})
	suite.NotNil(err)
	suite.Nil(product)
	suite.Nil(pagination)
}
