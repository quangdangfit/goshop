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

	"goshop/internal/user/model"
	"goshop/pkg/config"
	dbsMocks "goshop/pkg/dbs/mocks"
)

func newAddressSQLMockGormDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
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

type AddressRepositoryTestSuite struct {
	suite.Suite
	mockDB *dbsMocks.Database
	repo   AddressRepository
}

func (suite *AddressRepositoryTestSuite) SetupTest() {
	logger.Initialize(config.ProductionEnv)
	suite.mockDB = dbsMocks.NewDatabase(suite.T())
	suite.repo = NewAddressRepository(suite.mockDB)
}

func TestAddressRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(AddressRepositoryTestSuite))
}

// ListByUser
// =================================================================

func (suite *AddressRepositoryTestSuite) TestListByUserSuccess() {
	suite.mockDB.On("Find", mock.Anything, mock.Anything, mock.Anything).Return(nil).Times(1)

	addresses, err := suite.repo.ListByUser(context.Background(), "u1")
	suite.Nil(err)
	suite.Equal(0, len(addresses))
}

func (suite *AddressRepositoryTestSuite) TestListByUserFail() {
	suite.mockDB.On("Find", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("db error")).Times(1)

	addresses, err := suite.repo.ListByUser(context.Background(), "u1")
	suite.NotNil(err)
	suite.Nil(addresses)
}

// GetByID
// =================================================================

func (suite *AddressRepositoryTestSuite) TestGetByIDSuccess() {
	suite.mockDB.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(nil).Times(1)

	address, err := suite.repo.GetByID(context.Background(), "a1", "u1")
	suite.Nil(err)
	suite.NotNil(address)
}

func (suite *AddressRepositoryTestSuite) TestGetByIDFail() {
	suite.mockDB.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("not found")).Times(1)

	address, err := suite.repo.GetByID(context.Background(), "notfound", "u1")
	suite.NotNil(err)
	suite.Nil(address)
}

// Create
// =================================================================

func (suite *AddressRepositoryTestSuite) TestCreateSuccess() {
	address := &model.Address{UserID: "u1", Street: "123 Main St"}
	suite.mockDB.On("Create", mock.Anything, address).Return(nil).Times(1)

	err := suite.repo.Create(context.Background(), address)
	suite.Nil(err)
}

func (suite *AddressRepositoryTestSuite) TestCreateFail() {
	address := &model.Address{UserID: "u1", Street: "123 Main St"}
	suite.mockDB.On("Create", mock.Anything, address).Return(errors.New("db error")).Times(1)

	err := suite.repo.Create(context.Background(), address)
	suite.NotNil(err)
}

// Update
// =================================================================

func (suite *AddressRepositoryTestSuite) TestUpdateSuccess() {
	address := &model.Address{ID: "a1", UserID: "u1", Street: "456 Oak Ave"}
	suite.mockDB.On("Update", mock.Anything, address).Return(nil).Times(1)

	err := suite.repo.Update(context.Background(), address)
	suite.Nil(err)
}

func (suite *AddressRepositoryTestSuite) TestUpdateFail() {
	address := &model.Address{ID: "a1", UserID: "u1", Street: "456 Oak Ave"}
	suite.mockDB.On("Update", mock.Anything, address).Return(errors.New("db error")).Times(1)

	err := suite.repo.Update(context.Background(), address)
	suite.NotNil(err)
}

// Delete
// =================================================================

func (suite *AddressRepositoryTestSuite) TestDeleteSuccess() {
	suite.mockDB.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(nil).Times(1)
	suite.mockDB.On("Delete", mock.Anything, mock.Anything, mock.Anything).Return(nil).Times(1)

	err := suite.repo.Delete(context.Background(), "a1", "u1")
	suite.Nil(err)
}

func (suite *AddressRepositoryTestSuite) TestDeleteGetByIDFail() {
	suite.mockDB.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("not found")).Times(1)

	err := suite.repo.Delete(context.Background(), "notfound", "u1")
	suite.NotNil(err)
}

func (suite *AddressRepositoryTestSuite) TestDeleteFail() {
	suite.mockDB.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(nil).Times(1)
	suite.mockDB.On("Delete", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("db error")).Times(1)

	err := suite.repo.Delete(context.Background(), "a1", "u1")
	suite.NotNil(err)
}

// SetDefault
// =================================================================

func (suite *AddressRepositoryTestSuite) TestSetDefaultSuccess() {
	gormDB, sqlMock := newAddressSQLMockGormDB(suite.T())
	sqlMock.ExpectExec(".*").WillReturnResult(sqlmock.NewResult(0, 1))

	suite.mockDB.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(nil).Times(1)
	suite.mockDB.On("GetDB").Return(gormDB).Maybe()
	suite.mockDB.On("Update", mock.Anything, mock.Anything).Return(nil).Maybe()
	suite.mockDB.On("WithTransaction", mock.Anything).
		Run(func(args mock.Arguments) {
			fn := args.Get(0).(func() error)
			fn()
		}).
		Return(nil).Times(1)

	err := suite.repo.SetDefault(context.Background(), "a1", "u1")
	suite.Nil(err)
}

func (suite *AddressRepositoryTestSuite) TestSetDefaultGetByIDFail() {
	suite.mockDB.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("not found")).Times(1)

	err := suite.repo.SetDefault(context.Background(), "notfound", "u1")
	suite.NotNil(err)
}

func (suite *AddressRepositoryTestSuite) TestSetDefaultTransactionFail() {
	suite.mockDB.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(nil).Times(1)
	suite.mockDB.On("WithTransaction", mock.Anything).Return(errors.New("tx error")).Times(1)

	err := suite.repo.SetDefault(context.Background(), "a1", "u1")
	suite.NotNil(err)
}

func (suite *AddressRepositoryTestSuite) TestSetDefaultDBUpdateFail() {
	gormDB, sqlMock := newAddressSQLMockGormDB(suite.T())
	sqlMock.ExpectExec(".*").WillReturnError(errors.New("db error"))

	suite.mockDB.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(nil).Times(1)
	suite.mockDB.On("GetDB").Return(gormDB).Maybe()
	suite.mockDB.On("WithTransaction", mock.Anything).
		Run(func(args mock.Arguments) {
			fn := args.Get(0).(func() error)
			fn()
		}).
		Return(errors.New("db error")).Times(1)

	err := suite.repo.SetDefault(context.Background(), "a1", "u1")
	suite.NotNil(err)
}
