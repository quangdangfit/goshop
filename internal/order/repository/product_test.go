package repository

import (
	"context"
	"errors"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/quangdangfit/gocommon/logger"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"goshop/internal/order/model"
	"goshop/pkg/config"
	"goshop/pkg/dbs/mocks"
)

func newOrderProductSQLMockGormDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	t.Helper()
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sql mock: %v", err)
	}
	t.Cleanup(func() { sqlDB.Close() })
	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		t.Fatalf("failed to open gorm db: %v", err)
	}
	return gormDB, mock
}

type ProductRepositoryTestSuite struct {
	suite.Suite
	mockDB *mocks.Database
	repo   ProductRepository
}

func (suite *ProductRepositoryTestSuite) SetupTest() {
	logger.Initialize(config.ProductionEnv)

	suite.mockDB = mocks.NewDatabase(suite.T())
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

// DecrementStock
// =================================================================

func (suite *ProductRepositoryTestSuite) TestDecrementStockSuccess() {
	gormDB, sqlMock := newOrderProductSQLMockGormDB(suite.T())
	sqlMock.ExpectExec(".*").WillReturnResult(sqlmock.NewResult(0, 1))

	suite.mockDB.On("GetDB").Return(gormDB).Times(1)

	err := suite.repo.DecrementStock(context.Background(), "productId1", 2)
	suite.Nil(err)
}

func (suite *ProductRepositoryTestSuite) TestDecrementStockInsufficientStock() {
	gormDB, sqlMock := newOrderProductSQLMockGormDB(suite.T())
	sqlMock.ExpectExec(".*").WillReturnResult(sqlmock.NewResult(0, 0))

	suite.mockDB.On("GetDB").Return(gormDB).Times(1)

	err := suite.repo.DecrementStock(context.Background(), "productId1", 2)
	suite.NotNil(err)
	suite.Equal("insufficient stock", err.Error())
}

func (suite *ProductRepositoryTestSuite) TestDecrementStockFail() {
	gormDB, sqlMock := newOrderProductSQLMockGormDB(suite.T())
	sqlMock.ExpectExec(".*").WillReturnError(errors.New("db error"))

	suite.mockDB.On("GetDB").Return(gormDB).Times(1)

	err := suite.repo.DecrementStock(context.Background(), "productId1", 2)
	suite.NotNil(err)
}
