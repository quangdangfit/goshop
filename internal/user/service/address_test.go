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

// ListAddresses
// =================================================================================================

func (suite *AddressServiceTestSuite) TestListAddressesSuccess() {
	suite.mockRepo.On("ListByUser", mock.Anything, "u1").
		Return([]*model.Address{
			{ID: "a1", UserID: "u1", Street: "123 Main St"},
		}, nil).Times(1)

	addresses, err := suite.service.ListAddresses(context.Background(), "u1")
	suite.Nil(err)
	suite.Equal(1, len(addresses))
	suite.Equal("a1", addresses[0].ID)
}

func (suite *AddressServiceTestSuite) TestListAddressesFail() {
	suite.mockRepo.On("ListByUser", mock.Anything, "u1").
		Return(nil, errors.New("db error")).Times(1)

	addresses, err := suite.service.ListAddresses(context.Background(), "u1")
	suite.NotNil(err)
	suite.Nil(addresses)
}

// GetAddressByID
// =================================================================================================

func (suite *AddressServiceTestSuite) TestGetAddressByIDSuccess() {
	suite.mockRepo.On("GetByID", mock.Anything, "a1", "u1").
		Return(&model.Address{ID: "a1", UserID: "u1"}, nil).Times(1)

	addr, err := suite.service.GetAddressByID(context.Background(), "a1", "u1")
	suite.Nil(err)
	suite.Equal("a1", addr.ID)
}

func (suite *AddressServiceTestSuite) TestGetAddressByIDFail() {
	suite.mockRepo.On("GetByID", mock.Anything, "notfound", "u1").
		Return(nil, errors.New("not found")).Times(1)

	addr, err := suite.service.GetAddressByID(context.Background(), "notfound", "u1")
	suite.NotNil(err)
	suite.Nil(addr)
}

// Create
// =================================================================================================

func (suite *AddressServiceTestSuite) TestCreateSuccess() {
	req := &dto.CreateAddressReq{
		Name:    "Home",
		Phone:   "0123456789",
		Street:  "123 Main St",
		City:    "Springfield",
		Country: "US",
	}
	suite.mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil).Times(1)

	addr, err := suite.service.Create(context.Background(), "u1", req)
	suite.Nil(err)
	suite.NotNil(addr)
	suite.Equal("u1", addr.UserID)
	suite.Equal("123 Main St", addr.Street)
}

func (suite *AddressServiceTestSuite) TestCreateValidationFail() {
	req := &dto.CreateAddressReq{} // missing required fields

	addr, err := suite.service.Create(context.Background(), "u1", req)
	suite.NotNil(err)
	suite.Nil(addr)
}

func (suite *AddressServiceTestSuite) TestCreateDBFail() {
	req := &dto.CreateAddressReq{
		Name:    "Home",
		Phone:   "0123456789",
		Street:  "123 Main St",
		City:    "Springfield",
		Country: "US",
	}
	suite.mockRepo.On("Create", mock.Anything, mock.Anything).Return(errors.New("db error")).Times(1)

	addr, err := suite.service.Create(context.Background(), "u1", req)
	suite.NotNil(err)
	suite.Nil(addr)
}

// Update
// =================================================================================================

func (suite *AddressServiceTestSuite) TestUpdateSuccess() {
	req := &dto.UpdateAddressReq{Street: "789 Elm St"}

	suite.mockRepo.On("GetByID", mock.Anything, "a1", "u1").
		Return(&model.Address{ID: "a1", UserID: "u1", Street: "123 Main St"}, nil).Times(1)
	suite.mockRepo.On("Update", mock.Anything, mock.Anything).Return(nil).Times(1)

	addr, err := suite.service.Update(context.Background(), "a1", "u1", req)
	suite.Nil(err)
	suite.NotNil(addr)
	suite.Equal("789 Elm St", addr.Street)
}

func (suite *AddressServiceTestSuite) TestUpdateGetByIDFail() {
	req := &dto.UpdateAddressReq{Street: "789 Elm St"}

	suite.mockRepo.On("GetByID", mock.Anything, "notfound", "u1").
		Return(nil, errors.New("not found")).Times(1)

	addr, err := suite.service.Update(context.Background(), "notfound", "u1", req)
	suite.NotNil(err)
	suite.Nil(addr)
}

func (suite *AddressServiceTestSuite) TestUpdateDBFail() {
	req := &dto.UpdateAddressReq{Street: "789 Elm St"}

	suite.mockRepo.On("GetByID", mock.Anything, "a1", "u1").
		Return(&model.Address{ID: "a1", UserID: "u1"}, nil).Times(1)
	suite.mockRepo.On("Update", mock.Anything, mock.Anything).Return(errors.New("db error")).Times(1)

	addr, err := suite.service.Update(context.Background(), "a1", "u1", req)
	suite.NotNil(err)
	suite.Nil(addr)
}

// Delete
// =================================================================================================

func (suite *AddressServiceTestSuite) TestDeleteSuccess() {
	suite.mockRepo.On("Delete", mock.Anything, "a1", "u1").Return(nil).Times(1)

	err := suite.service.Delete(context.Background(), "a1", "u1")
	suite.Nil(err)
}

func (suite *AddressServiceTestSuite) TestDeleteFail() {
	suite.mockRepo.On("Delete", mock.Anything, "a1", "u1").Return(errors.New("not found")).Times(1)

	err := suite.service.Delete(context.Background(), "a1", "u1")
	suite.NotNil(err)
}

// SetDefault
// =================================================================================================

func (suite *AddressServiceTestSuite) TestSetDefaultSuccess() {
	suite.mockRepo.On("SetDefault", mock.Anything, "a1", "u1").Return(nil).Times(1)

	err := suite.service.SetDefault(context.Background(), "a1", "u1")
	suite.Nil(err)
}

func (suite *AddressServiceTestSuite) TestSetDefaultFail() {
	suite.mockRepo.On("SetDefault", mock.Anything, "a1", "u1").Return(errors.New("not found")).Times(1)

	err := suite.service.SetDefault(context.Background(), "a1", "u1")
	suite.NotNil(err)
}
