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

	"goshop/internal/product/domain"
	"goshop/internal/product/model"
	"goshop/pkg/config"
	"goshop/pkg/dbs/mocks"
)

func newProductSQLMockGormDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	t.Helper()
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sql mock: %v", err)
	}
	t.Cleanup(func() { _ = sqlDB.Close() })
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

func (suite *ProductRepositoryTestSuite) TestCreate() {
	tests := []struct {
		name    string
		setup   func()
		wantErr bool
	}{
		{
			name: "Success",
			setup: func() {
				suite.mockDB.On("Create", mock.Anything, mock.Anything).Return(nil).Times(1)
			},
		},
		{
			name: "DB error",
			setup: func() {
				suite.mockDB.On("Create", mock.Anything, mock.Anything).Return(errors.New("error")).Times(1)
			},
			wantErr: true,
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.setup()
			product := &model.Product{Name: "product name", Description: "product description", Price: 10.5}
			err := suite.repo.Create(context.Background(), product)
			if tc.wantErr {
				suite.NotNil(err)
			} else {
				suite.Nil(err)
			}
		})
	}
}

func (suite *ProductRepositoryTestSuite) TestUpdate() {
	tests := []struct {
		name    string
		setup   func()
		wantErr bool
	}{
		{
			name: "Success",
			setup: func() {
				suite.mockDB.On("Update", mock.Anything, mock.Anything).Return(nil).Times(1)
			},
		},
		{
			name: "DB error",
			setup: func() {
				suite.mockDB.On("Update", mock.Anything, mock.Anything).Return(errors.New("error")).Times(1)
			},
			wantErr: true,
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.setup()
			product := &model.Product{ID: "productId1", Name: "product name", Description: "product description", Price: 10.5}
			err := suite.repo.Update(context.Background(), product)
			if tc.wantErr {
				suite.NotNil(err)
			} else {
				suite.Nil(err)
			}
		})
	}
}

func (suite *ProductRepositoryTestSuite) TestGetProductByID() {
	tests := []struct {
		name    string
		setup   func()
		wantErr bool
	}{
		{
			name: "Success",
			setup: func() {
				suite.mockDB.On("FindById", mock.Anything, "productId1", &model.Product{}).Return(nil).Times(1)
			},
		},
		{
			name: "Not found",
			setup: func() {
				suite.mockDB.On("FindById", mock.Anything, "productId1", &model.Product{}).Return(errors.New("error")).Times(1)
			},
			wantErr: true,
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.setup()
			product, err := suite.repo.GetProductByID(context.Background(), "productId1")
			if tc.wantErr {
				suite.NotNil(err)
				suite.Nil(product)
			} else {
				suite.Nil(err)
				suite.NotNil(product)
			}
		})
	}
}

func (suite *ProductRepositoryTestSuite) TestListProducts() {
	tests := []struct {
		name    string
		req     *domain.ListProductReq
		setup   func()
		wantErr bool
	}{
		{
			name: "Success",
			req:  &domain.ListProductReq{Name: "name", Code: "code", Page: 2, Limit: 10, OrderBy: "name", OrderDesc: true},
			setup: func() {
				suite.mockDB.On("Count", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Times(1)
				suite.mockDB.On("Find", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Times(1)
			},
		},
		{
			name: "With CategoryID",
			req:  &domain.ListProductReq{CategoryID: "cat1"},
			setup: func() {
				suite.mockDB.On("Count", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Times(1)
				suite.mockDB.On("Find", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Times(1)
			},
		},
		{
			name: "Count fail",
			req:  &domain.ListProductReq{},
			setup: func() {
				suite.mockDB.On("Count", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("error")).Times(1)
			},
			wantErr: true,
		},
		{
			name: "Find fail",
			req:  &domain.ListProductReq{},
			setup: func() {
				suite.mockDB.On("Count", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Times(1)
				suite.mockDB.On("Find", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("error")).Times(1)
			},
			wantErr: true,
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.setup()
			products, pagination, err := suite.repo.ListProducts(context.Background(), tc.req)
			if tc.wantErr {
				suite.NotNil(err)
				suite.Nil(products)
				suite.Nil(pagination)
			} else {
				suite.Nil(err)
				suite.Equal(0, len(products))
				suite.NotNil(pagination)
			}
		})
	}
}

func (suite *ProductRepositoryTestSuite) TestDecrementStock() {
	tests := []struct {
		name    string
		setup   func()
		wantErr bool
		errMsg  string
	}{
		{
			name: "Success",
			setup: func() {
				gormDB, sqlMock := newProductSQLMockGormDB(suite.T())
				sqlMock.ExpectExec(".*").WillReturnResult(sqlmock.NewResult(0, 1))
				suite.mockDB.On("GetDB").Return(gormDB).Times(1)
			},
		},
		{
			name: "Insufficient stock",
			setup: func() {
				gormDB, sqlMock := newProductSQLMockGormDB(suite.T())
				sqlMock.ExpectExec(".*").WillReturnResult(sqlmock.NewResult(0, 0))
				suite.mockDB.On("GetDB").Return(gormDB).Times(1)
			},
			wantErr: true,
			errMsg:  "insufficient stock",
		},
		{
			name: "DB error",
			setup: func() {
				gormDB, sqlMock := newProductSQLMockGormDB(suite.T())
				sqlMock.ExpectExec(".*").WillReturnError(errors.New("db error"))
				suite.mockDB.On("GetDB").Return(gormDB).Times(1)
			},
			wantErr: true,
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.setup()
			err := suite.repo.DecrementStock(context.Background(), "productId1", 2)
			if tc.wantErr {
				suite.NotNil(err)
				if tc.errMsg != "" {
					suite.Equal(tc.errMsg, err.Error())
				}
			} else {
				suite.Nil(err)
			}
		})
	}
}

func (suite *ProductRepositoryTestSuite) TestUpdateRating() {
	tests := []struct {
		name    string
		setup   func()
		wantErr bool
	}{
		{
			name: "Success",
			setup: func() {
				gormDB, sqlMock := newProductSQLMockGormDB(suite.T())
				sqlMock.ExpectExec(".*").WillReturnResult(sqlmock.NewResult(0, 1))
				suite.mockDB.On("GetDB").Return(gormDB).Times(1)
			},
		},
		{
			name: "DB error",
			setup: func() {
				gormDB, sqlMock := newProductSQLMockGormDB(suite.T())
				sqlMock.ExpectExec(".*").WillReturnError(errors.New("db error"))
				suite.mockDB.On("GetDB").Return(gormDB).Times(1)
			},
			wantErr: true,
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.setup()
			err := suite.repo.UpdateRating(context.Background(), "productId1", 4.5, 10)
			if tc.wantErr {
				suite.NotNil(err)
			} else {
				suite.Nil(err)
			}
		})
	}
}
