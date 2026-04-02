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
	t.Cleanup(func() { _ = sqlDB.Close() })
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

func (suite *AddressRepositoryTestSuite) TestListByUser() {
	tests := []struct {
		name    string
		setup   func()
		wantErr bool
		wantLen int
	}{
		{
			name: "Success",
			setup: func() {
				suite.mockDB.On("Find", mock.Anything, mock.Anything, mock.Anything).Return(nil).Times(1)
			},
			wantLen: 0,
		},
		{
			name: "DB error",
			setup: func() {
				suite.mockDB.On("Find", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("db error")).Times(1)
			},
			wantErr: true,
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.setup()
			addresses, err := suite.repo.ListByUser(context.Background(), "u1")
			if tc.wantErr {
				suite.NotNil(err)
				suite.Nil(addresses)
			} else {
				suite.Nil(err)
				suite.Equal(tc.wantLen, len(addresses))
			}
		})
	}
}

func (suite *AddressRepositoryTestSuite) TestGetByID() {
	tests := []struct {
		name      string
		addressID string
		setup     func()
		wantErr   bool
	}{
		{
			name:      "Success",
			addressID: "a1",
			setup: func() {
				suite.mockDB.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(nil).Times(1)
			},
		},
		{
			name:      "Not found",
			addressID: "notfound",
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
			address, err := suite.repo.GetByID(context.Background(), tc.addressID, "u1")
			if tc.wantErr {
				suite.NotNil(err)
				suite.Nil(address)
			} else {
				suite.Nil(err)
				suite.NotNil(address)
			}
		})
	}
}

func (suite *AddressRepositoryTestSuite) TestCreate() {
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
			address := &model.Address{UserID: "u1", Street: "123 Main St"}
			err := suite.repo.Create(context.Background(), address)
			if tc.wantErr {
				suite.NotNil(err)
			} else {
				suite.Nil(err)
			}
		})
	}
}

func (suite *AddressRepositoryTestSuite) TestUpdate() {
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
			address := &model.Address{ID: "a1", UserID: "u1", Street: "456 Oak Ave"}
			err := suite.repo.Update(context.Background(), address)
			if tc.wantErr {
				suite.NotNil(err)
			} else {
				suite.Nil(err)
			}
		})
	}
}

func (suite *AddressRepositoryTestSuite) TestDelete() {
	tests := []struct {
		name      string
		addressID string
		setup     func()
		wantErr   bool
	}{
		{
			name:      "Success",
			addressID: "a1",
			setup: func() {
				suite.mockDB.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(nil).Times(1)
				suite.mockDB.On("Delete", mock.Anything, mock.Anything, mock.Anything).Return(nil).Times(1)
			},
		},
		{
			name:      "GetByID fail",
			addressID: "notfound",
			setup: func() {
				suite.mockDB.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("not found")).Times(1)
			},
			wantErr: true,
		},
		{
			name:      "Delete fail",
			addressID: "a1",
			setup: func() {
				suite.mockDB.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(nil).Times(1)
				suite.mockDB.On("Delete", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("db error")).Times(1)
			},
			wantErr: true,
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.setup()
			err := suite.repo.Delete(context.Background(), tc.addressID, "u1")
			if tc.wantErr {
				suite.NotNil(err)
			} else {
				suite.Nil(err)
			}
		})
	}
}

func (suite *AddressRepositoryTestSuite) TestSetDefault() {
	tests := []struct {
		name      string
		addressID string
		setup     func()
		wantErr   bool
	}{
		{
			name:      "Success",
			addressID: "a1",
			setup: func() {
				gormDB, sqlMock := newAddressSQLMockGormDB(suite.T())
				sqlMock.ExpectExec(".*").WillReturnResult(sqlmock.NewResult(0, 1))

				suite.mockDB.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(nil).Times(1)
				suite.mockDB.On("GetDB").Return(gormDB).Maybe()
				suite.mockDB.On("Update", mock.Anything, mock.Anything).Return(nil).Maybe()
				suite.mockDB.On("WithTransaction", mock.Anything).
					Run(func(args mock.Arguments) {
						fn := args.Get(0).(func() error)
						_ = fn()
					}).
					Return(nil).Times(1)
			},
		},
		{
			name:      "GetByID fail",
			addressID: "notfound",
			setup: func() {
				suite.mockDB.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("not found")).Times(1)
			},
			wantErr: true,
		},
		{
			name:      "Transaction fail",
			addressID: "a1",
			setup: func() {
				suite.mockDB.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(nil).Times(1)
				suite.mockDB.On("WithTransaction", mock.Anything).Return(errors.New("tx error")).Times(1)
			},
			wantErr: true,
		},
		{
			name:      "DB update fail",
			addressID: "a1",
			setup: func() {
				gormDB, sqlMock := newAddressSQLMockGormDB(suite.T())
				sqlMock.ExpectExec(".*").WillReturnError(errors.New("db error"))

				suite.mockDB.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(nil).Times(1)
				suite.mockDB.On("GetDB").Return(gormDB).Maybe()
				suite.mockDB.On("WithTransaction", mock.Anything).
					Run(func(args mock.Arguments) {
						fn := args.Get(0).(func() error)
						_ = fn()
					}).
					Return(errors.New("db error")).Times(1)
			},
			wantErr: true,
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.setup()
			err := suite.repo.SetDefault(context.Background(), tc.addressID, "u1")
			if tc.wantErr {
				suite.NotNil(err)
			} else {
				suite.Nil(err)
			}
		})
	}
}
