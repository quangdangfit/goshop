package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/quangdangfit/gocommon/logger"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"goshop/internal/user/dto"
	"goshop/internal/user/model"
	srvMocks "goshop/internal/user/service/mocks"
	"goshop/pkg/config"
	"goshop/pkg/response"
	"goshop/pkg/utils"
)

type AddressHandlerTestSuite struct {
	suite.Suite
	mockService *srvMocks.AddressService
	handler     *AddressHandler
}

func (suite *AddressHandlerTestSuite) SetupTest() {
	logger.Initialize(config.ProductionEnv)
	suite.mockService = srvMocks.NewAddressService(suite.T())
	suite.handler = NewAddressHandler(suite.mockService)
}

func TestAddressHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(AddressHandlerTestSuite))
}

func (suite *AddressHandlerTestSuite) prepareContext(method, path string, body any) (*gin.Context, *httptest.ResponseRecorder) {
	requestBody, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(method, path, bytes.NewBuffer(requestBody))
	c, _ := gin.CreateTestContext(w)
	c.Request = r
	return c, w
}

// ListAddresses
// =================================================================================================

func (suite *AddressHandlerTestSuite) TestListAddressesSuccess() {
	ctx, writer := suite.prepareContext("GET", "/api/v1/addresses", nil)
	ctx.Set("userId", "u1")

	suite.mockService.On("ListAddresses", mock.Anything, "u1").
		Return([]*model.Address{
			{ID: "a1", UserID: "u1", Street: "123 Main St"},
			{ID: "a2", UserID: "u1", Street: "456 Oak Ave"},
		}, nil).Times(1)

	suite.handler.ListAddresses(ctx)

	var res response.Response
	var addresses []*dto.Address
	_ = json.Unmarshal(writer.Body.Bytes(), &res)
	utils.Copy(&addresses, &res.Result)
	suite.Equal(http.StatusOK, writer.Code)
	suite.Equal(2, len(addresses))
}

func (suite *AddressHandlerTestSuite) TestListAddressesUnauthorized() {
	ctx, writer := suite.prepareContext("GET", "/api/v1/addresses", nil)

	suite.handler.ListAddresses(ctx)

	suite.Equal(http.StatusUnauthorized, writer.Code)
}

func (suite *AddressHandlerTestSuite) TestListAddressesFail() {
	ctx, writer := suite.prepareContext("GET", "/api/v1/addresses", nil)
	ctx.Set("userId", "u1")

	suite.mockService.On("ListAddresses", mock.Anything, "u1").
		Return(nil, errors.New("db error")).Times(1)

	suite.handler.ListAddresses(ctx)

	suite.Equal(http.StatusInternalServerError, writer.Code)
}

// GetAddressByID
// =================================================================================================

func (suite *AddressHandlerTestSuite) TestGetAddressByIDSuccess() {
	ctx, writer := suite.prepareContext("GET", "/api/v1/addresses/a1", nil)
	ctx.Set("userId", "u1")
	ctx.Params = gin.Params{{Key: "id", Value: "a1"}}

	suite.mockService.On("GetAddressByID", mock.Anything, "a1", "u1").
		Return(&model.Address{ID: "a1", UserID: "u1", Street: "123 Main St"}, nil).Times(1)

	suite.handler.GetAddressByID(ctx)

	var res response.Response
	var address dto.Address
	_ = json.Unmarshal(writer.Body.Bytes(), &res)
	utils.Copy(&address, &res.Result)
	suite.Equal(http.StatusOK, writer.Code)
	suite.Equal("a1", address.ID)
}

func (suite *AddressHandlerTestSuite) TestGetAddressByIDUnauthorized() {
	ctx, writer := suite.prepareContext("GET", "/api/v1/addresses/a1", nil)

	suite.handler.GetAddressByID(ctx)

	suite.Equal(http.StatusUnauthorized, writer.Code)
}

func (suite *AddressHandlerTestSuite) TestGetAddressByIDMissingID() {
	ctx, writer := suite.prepareContext("GET", "/api/v1/addresses/", nil)
	ctx.Set("userId", "u1")
	// id param is empty

	suite.handler.GetAddressByID(ctx)

	suite.Equal(http.StatusBadRequest, writer.Code)
}

func (suite *AddressHandlerTestSuite) TestGetAddressByIDNotFound() {
	ctx, writer := suite.prepareContext("GET", "/api/v1/addresses/notfound", nil)
	ctx.Set("userId", "u1")
	ctx.Params = gin.Params{{Key: "id", Value: "notfound"}}

	suite.mockService.On("GetAddressByID", mock.Anything, "notfound", "u1").
		Return(nil, errors.New("not found")).Times(1)

	suite.handler.GetAddressByID(ctx)

	suite.Equal(http.StatusNotFound, writer.Code)
}

// CreateAddress
// =================================================================================================

func (suite *AddressHandlerTestSuite) TestCreateAddressSuccess() {
	req := &dto.CreateAddressReq{
		Street: "123 Main St",
		City:   "Springfield",
	}
	ctx, writer := suite.prepareContext("POST", "/api/v1/addresses", req)
	ctx.Set("userId", "u1")

	suite.mockService.On("Create", mock.Anything, "u1", req).
		Return(&model.Address{ID: "a1", UserID: "u1", Street: "123 Main St"}, nil).Times(1)

	suite.handler.CreateAddress(ctx)

	var res response.Response
	var address dto.Address
	_ = json.Unmarshal(writer.Body.Bytes(), &res)
	utils.Copy(&address, &res.Result)
	suite.Equal(http.StatusOK, writer.Code)
	suite.Equal("a1", address.ID)
}

func (suite *AddressHandlerTestSuite) TestCreateAddressUnauthorized() {
	req := &dto.CreateAddressReq{Street: "123 Main St"}
	ctx, writer := suite.prepareContext("POST", "/api/v1/addresses", req)

	suite.handler.CreateAddress(ctx)

	suite.Equal(http.StatusUnauthorized, writer.Code)
}

func (suite *AddressHandlerTestSuite) TestCreateAddressInvalidBody() {
	req := map[string]any{"street": 123}
	ctx, writer := suite.prepareContext("POST", "/api/v1/addresses", req)
	ctx.Set("userId", "u1")

	suite.handler.CreateAddress(ctx)

	suite.Equal(http.StatusBadRequest, writer.Code)
}

func (suite *AddressHandlerTestSuite) TestCreateAddressFail() {
	req := &dto.CreateAddressReq{Street: "123 Main St", City: "Springfield"}
	ctx, writer := suite.prepareContext("POST", "/api/v1/addresses", req)
	ctx.Set("userId", "u1")

	suite.mockService.On("Create", mock.Anything, "u1", req).
		Return(nil, errors.New("db error")).Times(1)

	suite.handler.CreateAddress(ctx)

	suite.Equal(http.StatusInternalServerError, writer.Code)
}

// UpdateAddress
// =================================================================================================

func (suite *AddressHandlerTestSuite) TestUpdateAddressSuccess() {
	req := &dto.UpdateAddressReq{Street: "789 Elm St"}
	ctx, writer := suite.prepareContext("PUT", "/api/v1/addresses/a1", req)
	ctx.Set("userId", "u1")
	ctx.Params = gin.Params{{Key: "id", Value: "a1"}}

	suite.mockService.On("Update", mock.Anything, "a1", "u1", req).
		Return(&model.Address{ID: "a1", UserID: "u1", Street: "789 Elm St"}, nil).Times(1)

	suite.handler.UpdateAddress(ctx)

	var res response.Response
	var address dto.Address
	_ = json.Unmarshal(writer.Body.Bytes(), &res)
	utils.Copy(&address, &res.Result)
	suite.Equal(http.StatusOK, writer.Code)
	suite.Equal("a1", address.ID)
}

func (suite *AddressHandlerTestSuite) TestUpdateAddressUnauthorized() {
	req := &dto.UpdateAddressReq{Street: "789 Elm St"}
	ctx, writer := suite.prepareContext("PUT", "/api/v1/addresses/a1", req)

	suite.handler.UpdateAddress(ctx)

	suite.Equal(http.StatusUnauthorized, writer.Code)
}

func (suite *AddressHandlerTestSuite) TestUpdateAddressInvalidBody() {
	req := map[string]any{"street": 123}
	ctx, writer := suite.prepareContext("PUT", "/api/v1/addresses/a1", req)
	ctx.Set("userId", "u1")
	ctx.Params = gin.Params{{Key: "id", Value: "a1"}}

	suite.handler.UpdateAddress(ctx)

	suite.Equal(http.StatusBadRequest, writer.Code)
}

func (suite *AddressHandlerTestSuite) TestUpdateAddressFail() {
	req := &dto.UpdateAddressReq{Street: "789 Elm St"}
	ctx, writer := suite.prepareContext("PUT", "/api/v1/addresses/a1", req)
	ctx.Set("userId", "u1")
	ctx.Params = gin.Params{{Key: "id", Value: "a1"}}

	suite.mockService.On("Update", mock.Anything, "a1", "u1", req).
		Return(nil, errors.New("not found")).Times(1)

	suite.handler.UpdateAddress(ctx)

	suite.Equal(http.StatusInternalServerError, writer.Code)
}

// DeleteAddress
// =================================================================================================

func (suite *AddressHandlerTestSuite) TestDeleteAddressSuccess() {
	ctx, writer := suite.prepareContext("DELETE", "/api/v1/addresses/a1", nil)
	ctx.Set("userId", "u1")
	ctx.Params = gin.Params{{Key: "id", Value: "a1"}}

	suite.mockService.On("Delete", mock.Anything, "a1", "u1").Return(nil).Times(1)

	suite.handler.DeleteAddress(ctx)

	suite.Equal(http.StatusOK, writer.Code)
}

func (suite *AddressHandlerTestSuite) TestDeleteAddressUnauthorized() {
	ctx, writer := suite.prepareContext("DELETE", "/api/v1/addresses/a1", nil)

	suite.handler.DeleteAddress(ctx)

	suite.Equal(http.StatusUnauthorized, writer.Code)
}

func (suite *AddressHandlerTestSuite) TestDeleteAddressFail() {
	ctx, writer := suite.prepareContext("DELETE", "/api/v1/addresses/a1", nil)
	ctx.Set("userId", "u1")
	ctx.Params = gin.Params{{Key: "id", Value: "a1"}}

	suite.mockService.On("Delete", mock.Anything, "a1", "u1").Return(errors.New("not found")).Times(1)

	suite.handler.DeleteAddress(ctx)

	suite.Equal(http.StatusInternalServerError, writer.Code)
}

// SetDefaultAddress
// =================================================================================================

func (suite *AddressHandlerTestSuite) TestSetDefaultAddressSuccess() {
	ctx, writer := suite.prepareContext("PUT", "/api/v1/addresses/a1/default", nil)
	ctx.Set("userId", "u1")
	ctx.Params = gin.Params{{Key: "id", Value: "a1"}}

	suite.mockService.On("SetDefault", mock.Anything, "a1", "u1").Return(nil).Times(1)

	suite.handler.SetDefaultAddress(ctx)

	suite.Equal(http.StatusOK, writer.Code)
}

func (suite *AddressHandlerTestSuite) TestSetDefaultAddressUnauthorized() {
	ctx, writer := suite.prepareContext("PUT", "/api/v1/addresses/a1/default", nil)

	suite.handler.SetDefaultAddress(ctx)

	suite.Equal(http.StatusUnauthorized, writer.Code)
}

func (suite *AddressHandlerTestSuite) TestSetDefaultAddressFail() {
	ctx, writer := suite.prepareContext("PUT", "/api/v1/addresses/a1/default", nil)
	ctx.Set("userId", "u1")
	ctx.Params = gin.Params{{Key: "id", Value: "a1"}}

	suite.mockService.On("SetDefault", mock.Anything, "a1", "u1").Return(errors.New("not found")).Times(1)

	suite.handler.SetDefaultAddress(ctx)

	suite.Equal(http.StatusInternalServerError, writer.Code)
}
