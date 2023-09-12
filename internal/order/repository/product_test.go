package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/quangdangfit/gocommon/logger"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"goshop/internal/order/model"
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
