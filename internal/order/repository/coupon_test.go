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
	t.Cleanup(func() { _ = sqlDB.Close() })
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

func (suite *CouponRepositoryTestSuite) TestGetByCode() {
	tests := []struct {
		name    string
		code    string
		setup   func()
		wantErr bool
	}{
		{
			name: "Success",
			code: "SAVE10",
			setup: func() {
				suite.mockDB.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(nil).Times(1)
			},
		},
		{
			name: "Not found",
			code: "INVALID",
			setup: func() {
				suite.mockDB.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("not found")).Times(1)
			},
			wantErr: true,
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.setup()
			coupon, err := suite.repo.GetByCode(context.Background(), tc.code)
			if tc.wantErr {
				suite.NotNil(err)
				suite.Nil(coupon)
			} else {
				suite.Nil(err)
				suite.NotNil(coupon)
			}
		})
	}
}

func (suite *CouponRepositoryTestSuite) TestCreate() {
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
			name: "Duplicate",
			setup: func() {
				suite.mockDB.On("Create", mock.Anything, mock.Anything).Return(errors.New("duplicate")).Times(1)
			},
			wantErr: true,
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.setup()
			coupon := &model.Coupon{Code: "SAVE10"}
			err := suite.repo.Create(context.Background(), coupon)
			if tc.wantErr {
				suite.NotNil(err)
			} else {
				suite.Nil(err)
			}
		})
	}
}

func (suite *CouponRepositoryTestSuite) TestIncrUsedCount() {
	tests := []struct {
		name    string
		setup   func()
		wantErr bool
	}{
		{
			name: "Success",
			setup: func() {
				gormDB, sqlMock := newCouponSQLMockGormDB(suite.T())
				sqlMock.ExpectExec(".*").WillReturnResult(sqlmock.NewResult(0, 1))
				suite.mockDB.On("GetDB").Return(gormDB).Times(1)
			},
		},
		{
			name: "DB error",
			setup: func() {
				gormDB, sqlMock := newCouponSQLMockGormDB(suite.T())
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
			err := suite.repo.IncrUsedCount(context.Background(), "c1")
			if tc.wantErr {
				suite.NotNil(err)
			} else {
				suite.Nil(err)
			}
		})
	}
}
