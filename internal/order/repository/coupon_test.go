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
	dbsMocks "goshop/pkg/dbs/mocks"
)

func newCouponSQLMockGormDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
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

type CouponRepositoryTestSuite struct {
	suite.Suite
	mockDB *dbsMocks.Database
	repo   CouponRepository
}

func (suite *CouponRepositoryTestSuite) SetupTest() {
	logger.Initialize(config.ProductionEnv)
	suite.mockDB = dbsMocks.NewDatabase(suite.T())
	suite.repo = NewCouponRepository(suite.mockDB)
}

func TestCouponRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(CouponRepositoryTestSuite))
}

// GetByCode
// =================================================================

func (suite *CouponRepositoryTestSuite) TestGetByCodeSuccess() {
	suite.mockDB.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(nil).Times(1)

	coupon, err := suite.repo.GetByCode(context.Background(), "SAVE10")
	suite.Nil(err)
	suite.NotNil(coupon)
}

func (suite *CouponRepositoryTestSuite) TestGetByCodeFail() {
	suite.mockDB.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("not found")).Times(1)

	coupon, err := suite.repo.GetByCode(context.Background(), "INVALID")
	suite.NotNil(err)
	suite.Nil(coupon)
}

// Create
// =================================================================

func (suite *CouponRepositoryTestSuite) TestCreateSuccess() {
	coupon := &model.Coupon{Code: "SAVE10"}
	suite.mockDB.On("Create", mock.Anything, coupon).Return(nil).Times(1)

	err := suite.repo.Create(context.Background(), coupon)
	suite.Nil(err)
}

func (suite *CouponRepositoryTestSuite) TestCreateFail() {
	coupon := &model.Coupon{Code: "SAVE10"}
	suite.mockDB.On("Create", mock.Anything, coupon).Return(errors.New("duplicate")).Times(1)

	err := suite.repo.Create(context.Background(), coupon)
	suite.NotNil(err)
}

// IncrUsedCount
// =================================================================

func (suite *CouponRepositoryTestSuite) TestIncrUsedCountSuccess() {
	gormDB, sqlMock := newCouponSQLMockGormDB(suite.T())
	sqlMock.ExpectExec(".*").WillReturnResult(sqlmock.NewResult(0, 1))

	suite.mockDB.On("GetDB").Return(gormDB).Times(1)

	err := suite.repo.IncrUsedCount(context.Background(), "c1")
	suite.Nil(err)
}

func (suite *CouponRepositoryTestSuite) TestIncrUsedCountFail() {
	gormDB, sqlMock := newCouponSQLMockGormDB(suite.T())
	sqlMock.ExpectExec(".*").WillReturnError(errors.New("db error"))

	suite.mockDB.On("GetDB").Return(gormDB).Times(1)

	err := suite.repo.IncrUsedCount(context.Background(), "c1")
	suite.NotNil(err)
}
