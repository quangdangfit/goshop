package service

import (
	"context"
	"errors"
	"testing"

	"github.com/quangdangfit/gocommon/logger"
	"github.com/quangdangfit/gocommon/validation"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"goshop/internal/user/dto"
	"goshop/internal/user/model"
	"goshop/internal/user/repository/mocks"
	"goshop/pkg/config"
)

type AddressServiceTestSuite struct {
	suite.Suite
	mockRepo *mocks.AddressRepository
	service  AddressService
}

func (suite *AddressServiceTestSuite) SetupTest() {
	logger.Initialize(config.ProductionEnv)
	validator := validation.New()
	suite.mockRepo = mocks.NewAddressRepository(suite.T())
	suite.service = NewAddressService(validator, suite.mockRepo)
}

func TestAddressServiceTestSuite(t *testing.T) {
	suite.Run(t, new(AddressServiceTestSuite))
}

func (suite *AddressServiceTestSuite) TestListAddresses() {
	tests := []struct {
		name    string
		setup   func()
		wantErr bool
		wantLen int
	}{
		{
			name: "Success",
			setup: func() {
				suite.mockRepo.On("ListByUser", mock.Anything, "u1").
					Return([]*model.Address{{ID: "a1", UserID: "u1", Street: "123 Main St"}}, nil).Times(1)
			},
			wantLen: 1,
		},
		{
			name: "DB error",
			setup: func() {
				suite.mockRepo.On("ListByUser", mock.Anything, "u1").
					Return(nil, errors.New("db error")).Times(1)
			},
			wantErr: true,
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.setup()
			addresses, err := suite.service.ListAddresses(context.Background(), "u1")
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

func (suite *AddressServiceTestSuite) TestGetAddressByID() {
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
				suite.mockRepo.On("GetByID", mock.Anything, "a1", "u1").
					Return(&model.Address{ID: "a1", UserID: "u1"}, nil).Times(1)
			},
		},
		{
			name:      "Not found",
			addressID: "notfound",
			setup: func() {
				suite.mockRepo.On("GetByID", mock.Anything, "notfound", "u1").
					Return(nil, errors.New("not found")).Times(1)
			},
			wantErr: true,
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.setup()
			addr, err := suite.service.GetAddressByID(context.Background(), tc.addressID, "u1")
			if tc.wantErr {
				suite.NotNil(err)
				suite.Nil(addr)
			} else {
				suite.Nil(err)
				suite.Equal(tc.addressID, addr.ID)
			}
		})
	}
}

func (suite *AddressServiceTestSuite) TestCreate() {
	tests := []struct {
		name    string
		req     *dto.CreateAddressReq
		setup   func()
		wantErr bool
	}{
		{
			name: "Success",
			req:  &dto.CreateAddressReq{Name: "Home", Phone: "0123456789", Street: "123 Main St", City: "Springfield", Country: "US"},
			setup: func() {
				suite.mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil).Times(1)
			},
		},
		{
			name:    "Validation fail",
			req:     &dto.CreateAddressReq{},
			setup:   func() {},
			wantErr: true,
		},
		{
			name: "DB fail",
			req:  &dto.CreateAddressReq{Name: "Home", Phone: "0123456789", Street: "123 Main St", City: "Springfield", Country: "US"},
			setup: func() {
				suite.mockRepo.On("Create", mock.Anything, mock.Anything).Return(errors.New("db error")).Times(1)
			},
			wantErr: true,
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.setup()
			addr, err := suite.service.Create(context.Background(), "u1", tc.req)
			if tc.wantErr {
				suite.NotNil(err)
				suite.Nil(addr)
			} else {
				suite.Nil(err)
				suite.NotNil(addr)
				suite.Equal("u1", addr.UserID)
				suite.Equal("123 Main St", addr.Street)
			}
		})
	}
}

func (suite *AddressServiceTestSuite) TestUpdate() {
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
				suite.mockRepo.On("GetByID", mock.Anything, "a1", "u1").
					Return(&model.Address{ID: "a1", UserID: "u1", Street: "123 Main St"}, nil).Times(1)
				suite.mockRepo.On("Update", mock.Anything, mock.Anything).Return(nil).Times(1)
			},
		},
		{
			name:      "GetByID fail",
			addressID: "notfound",
			setup: func() {
				suite.mockRepo.On("GetByID", mock.Anything, "notfound", "u1").
					Return(nil, errors.New("not found")).Times(1)
			},
			wantErr: true,
		},
		{
			name:      "DB fail",
			addressID: "a1",
			setup: func() {
				suite.mockRepo.On("GetByID", mock.Anything, "a1", "u1").
					Return(&model.Address{ID: "a1", UserID: "u1"}, nil).Times(1)
				suite.mockRepo.On("Update", mock.Anything, mock.Anything).Return(errors.New("db error")).Times(1)
			},
			wantErr: true,
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.setup()
			req := &dto.UpdateAddressReq{Street: "789 Elm St"}
			addr, err := suite.service.Update(context.Background(), tc.addressID, "u1", req)
			if tc.wantErr {
				suite.NotNil(err)
				suite.Nil(addr)
			} else {
				suite.Nil(err)
				suite.NotNil(addr)
				suite.Equal("789 Elm St", addr.Street)
			}
		})
	}
}

func (suite *AddressServiceTestSuite) TestDelete() {
	tests := []struct {
		name    string
		setup   func()
		wantErr bool
	}{
		{
			name: "Success",
			setup: func() {
				suite.mockRepo.On("Delete", mock.Anything, "a1", "u1").Return(nil).Times(1)
			},
		},
		{
			name: "Not found",
			setup: func() {
				suite.mockRepo.On("Delete", mock.Anything, "a1", "u1").Return(errors.New("not found")).Times(1)
			},
			wantErr: true,
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.setup()
			err := suite.service.Delete(context.Background(), "a1", "u1")
			if tc.wantErr {
				suite.NotNil(err)
			} else {
				suite.Nil(err)
			}
		})
	}
}

func (suite *AddressServiceTestSuite) TestSetDefault() {
	tests := []struct {
		name    string
		setup   func()
		wantErr bool
	}{
		{
			name: "Success",
			setup: func() {
				suite.mockRepo.On("SetDefault", mock.Anything, "a1", "u1").Return(nil).Times(1)
			},
		},
		{
			name: "Not found",
			setup: func() {
				suite.mockRepo.On("SetDefault", mock.Anything, "a1", "u1").Return(errors.New("not found")).Times(1)
			},
			wantErr: true,
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.setup()
			err := suite.service.SetDefault(context.Background(), "a1", "u1")
			if tc.wantErr {
				suite.NotNil(err)
			} else {
				suite.Nil(err)
			}
		})
	}
}
