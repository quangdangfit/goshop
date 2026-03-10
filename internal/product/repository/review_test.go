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
	t.Cleanup(func() { sqlDB.Close() })
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

// ListByProduct
// =================================================================

func (suite *ReviewRepositoryTestSuite) TestListByProductSuccess() {
	suite.mockDB.On("Count", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Times(1)
	suite.mockDB.On("Find", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Times(1)

	reviews, pg, err := suite.repo.ListByProduct(context.Background(), "p1", 1, 10)
	suite.Nil(err)
	suite.Equal(0, len(reviews))
	suite.NotNil(pg)
}

func (suite *ReviewRepositoryTestSuite) TestListByProductCountFail() {
	suite.mockDB.On("Count", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("db error")).Times(1)

	reviews, pg, err := suite.repo.ListByProduct(context.Background(), "p1", 1, 10)
	suite.NotNil(err)
	suite.Nil(reviews)
	suite.Nil(pg)
}

func (suite *ReviewRepositoryTestSuite) TestListByProductFindFail() {
	suite.mockDB.On("Count", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Times(1)
	suite.mockDB.On("Find", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("db error")).Times(1)

	reviews, pg, err := suite.repo.ListByProduct(context.Background(), "p1", 1, 10)
	suite.NotNil(err)
	suite.Nil(reviews)
	suite.Nil(pg)
}

// GetByID
// =================================================================

func (suite *ReviewRepositoryTestSuite) TestGetByIDSuccess() {
	suite.mockDB.On("FindById", mock.Anything, "r1", &model.Review{}).Return(nil).Times(1)

	review, err := suite.repo.GetByID(context.Background(), "r1")
	suite.Nil(err)
	suite.NotNil(review)
}

func (suite *ReviewRepositoryTestSuite) TestGetByIDFail() {
	suite.mockDB.On("FindById", mock.Anything, "notfound", &model.Review{}).Return(errors.New("not found")).Times(1)

	review, err := suite.repo.GetByID(context.Background(), "notfound")
	suite.NotNil(err)
	suite.Nil(review)
}

// Create
// =================================================================

func (suite *ReviewRepositoryTestSuite) TestCreateSuccess() {
	review := &model.Review{ProductID: "p1", UserID: "u1", Rating: 5}
	suite.mockDB.On("Create", mock.Anything, review).Return(nil).Times(1)

	err := suite.repo.Create(context.Background(), review)
	suite.Nil(err)
}

func (suite *ReviewRepositoryTestSuite) TestCreateFail() {
	review := &model.Review{ProductID: "p1", UserID: "u1", Rating: 5}
	suite.mockDB.On("Create", mock.Anything, review).Return(errors.New("db error")).Times(1)

	err := suite.repo.Create(context.Background(), review)
	suite.NotNil(err)
}

// Update
// =================================================================

func (suite *ReviewRepositoryTestSuite) TestUpdateSuccess() {
	review := &model.Review{ID: "r1", Rating: 4, Comment: "Updated"}
	suite.mockDB.On("Update", mock.Anything, review).Return(nil).Times(1)

	err := suite.repo.Update(context.Background(), review)
	suite.Nil(err)
}

func (suite *ReviewRepositoryTestSuite) TestUpdateFail() {
	review := &model.Review{ID: "r1", Rating: 4, Comment: "Updated"}
	suite.mockDB.On("Update", mock.Anything, review).Return(errors.New("db error")).Times(1)

	err := suite.repo.Update(context.Background(), review)
	suite.NotNil(err)
}

// Delete
// =================================================================

func (suite *ReviewRepositoryTestSuite) TestDeleteSuccess() {
	suite.mockDB.On("Delete", mock.Anything, mock.Anything, mock.Anything).Return(nil).Times(1)

	err := suite.repo.Delete(context.Background(), "r1", "u1")
	suite.Nil(err)
}

func (suite *ReviewRepositoryTestSuite) TestDeleteFail() {
	suite.mockDB.On("Delete", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("db error")).Times(1)

	err := suite.repo.Delete(context.Background(), "r1", "u1")
	suite.NotNil(err)
}

// GetAggregates
// =================================================================

func (suite *ReviewRepositoryTestSuite) TestGetAggregatesSuccess() {
	gormDB, sqlMock := newReviewSQLMockGormDB(suite.T())
	sqlMock.ExpectQuery(".*").
		WillReturnRows(sqlmock.NewRows([]string{"avg", "count"}).AddRow(4.5, 10))

	suite.mockDB.On("GetDB").Return(gormDB).Times(1)

	avg, count, err := suite.repo.GetAggregates(context.Background(), "p1")
	suite.Nil(err)
	suite.Equal(float64(4.5), avg)
	suite.Equal(10, count)
}

func (suite *ReviewRepositoryTestSuite) TestGetAggregatesFail() {
	gormDB, sqlMock := newReviewSQLMockGormDB(suite.T())
	sqlMock.ExpectQuery(".*").WillReturnError(errors.New("db error"))

	suite.mockDB.On("GetDB").Return(gormDB).Times(1)

	avg, count, err := suite.repo.GetAggregates(context.Background(), "p1")
	suite.NotNil(err)
	suite.Equal(float64(0), avg)
	suite.Equal(0, count)
}
