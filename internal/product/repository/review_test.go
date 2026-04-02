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

	"goshop/internal/product/model"
	"goshop/pkg/config"
	dbsMocks "goshop/pkg/dbs/mocks"
)

func newReviewSQLMockGormDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
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

type ReviewRepositoryTestSuite struct {
	suite.Suite
	mockDB *dbsMocks.Database
	repo   ReviewRepository
}

func (suite *ReviewRepositoryTestSuite) SetupTest() {
	logger.Initialize(config.ProductionEnv)
	suite.mockDB = dbsMocks.NewDatabase(suite.T())
	suite.repo = NewReviewRepository(suite.mockDB)
}

func TestReviewRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(ReviewRepositoryTestSuite))
}

func (suite *ReviewRepositoryTestSuite) TestListByProduct() {
	tests := []struct {
		name    string
		setup   func()
		wantErr bool
	}{
		{
			name: "Success",
			setup: func() {
				suite.mockDB.On("Count", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Times(1)
				suite.mockDB.On("Find", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Times(1)
			},
		},
		{
			name: "Count fail",
			setup: func() {
				suite.mockDB.On("Count", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("db error")).Times(1)
			},
			wantErr: true,
		},
		{
			name: "Find fail",
			setup: func() {
				suite.mockDB.On("Count", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Times(1)
				suite.mockDB.On("Find", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("db error")).Times(1)
			},
			wantErr: true,
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.setup()
			reviews, pg, err := suite.repo.ListByProduct(context.Background(), "p1", 1, 10)
			if tc.wantErr {
				suite.NotNil(err)
				suite.Nil(reviews)
				suite.Nil(pg)
			} else {
				suite.Nil(err)
				suite.Equal(0, len(reviews))
				suite.NotNil(pg)
			}
		})
	}
}

func (suite *ReviewRepositoryTestSuite) TestGetByID() {
	tests := []struct {
		name    string
		id      string
		setup   func()
		wantErr bool
	}{
		{
			name: "Success",
			id:   "r1",
			setup: func() {
				suite.mockDB.On("FindById", mock.Anything, "r1", &model.Review{}).Return(nil).Times(1)
			},
		},
		{
			name: "Not found",
			id:   "notfound",
			setup: func() {
				suite.mockDB.On("FindById", mock.Anything, "notfound", &model.Review{}).Return(errors.New("not found")).Times(1)
			},
			wantErr: true,
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.setup()
			review, err := suite.repo.GetByID(context.Background(), tc.id)
			if tc.wantErr {
				suite.NotNil(err)
				suite.Nil(review)
			} else {
				suite.Nil(err)
				suite.NotNil(review)
			}
		})
	}
}

func (suite *ReviewRepositoryTestSuite) TestCreate() {
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
				suite.mockDB.On("Create", mock.Anything, mock.Anything).Return(errors.New("db error")).Times(1)
			},
			wantErr: true,
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.setup()
			review := &model.Review{ProductID: "p1", UserID: "u1", Rating: 5}
			err := suite.repo.Create(context.Background(), review)
			if tc.wantErr {
				suite.NotNil(err)
			} else {
				suite.Nil(err)
			}
		})
	}
}

func (suite *ReviewRepositoryTestSuite) TestUpdate() {
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
				suite.mockDB.On("Update", mock.Anything, mock.Anything).Return(errors.New("db error")).Times(1)
			},
			wantErr: true,
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.setup()
			review := &model.Review{ID: "r1", Rating: 4, Comment: "Updated"}
			err := suite.repo.Update(context.Background(), review)
			if tc.wantErr {
				suite.NotNil(err)
			} else {
				suite.Nil(err)
			}
		})
	}
}

func (suite *ReviewRepositoryTestSuite) TestDelete() {
	tests := []struct {
		name    string
		setup   func()
		wantErr bool
	}{
		{
			name: "Success",
			setup: func() {
				suite.mockDB.On("Delete", mock.Anything, mock.Anything, mock.Anything).Return(nil).Times(1)
			},
		},
		{
			name: "DB error",
			setup: func() {
				suite.mockDB.On("Delete", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("db error")).Times(1)
			},
			wantErr: true,
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.setup()
			err := suite.repo.Delete(context.Background(), "r1", "u1")
			if tc.wantErr {
				suite.NotNil(err)
			} else {
				suite.Nil(err)
			}
		})
	}
}

func (suite *ReviewRepositoryTestSuite) TestGetAggregates() {
	tests := []struct {
		name      string
		setup     func()
		wantErr   bool
		wantAvg   float64
		wantCount int
	}{
		{
			name: "Success",
			setup: func() {
				gormDB, sqlMock := newReviewSQLMockGormDB(suite.T())
				sqlMock.ExpectQuery(".*").
					WillReturnRows(sqlmock.NewRows([]string{"avg", "count"}).AddRow(4.5, 10))
				suite.mockDB.On("GetDB").Return(gormDB).Times(1)
			},
			wantAvg:   4.5,
			wantCount: 10,
		},
		{
			name: "DB error",
			setup: func() {
				gormDB, sqlMock := newReviewSQLMockGormDB(suite.T())
				sqlMock.ExpectQuery(".*").WillReturnError(errors.New("db error"))
				suite.mockDB.On("GetDB").Return(gormDB).Times(1)
			},
			wantErr:   true,
			wantAvg:   0,
			wantCount: 0,
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.setup()
			avg, count, err := suite.repo.GetAggregates(context.Background(), "p1")
			if tc.wantErr {
				suite.NotNil(err)
			} else {
				suite.Nil(err)
			}
			suite.Equal(tc.wantAvg, avg)
			suite.Equal(tc.wantCount, count)
		})
	}
}
