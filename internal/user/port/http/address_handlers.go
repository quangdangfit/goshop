package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/quangdangfit/gocommon/logger"

	"goshop/internal/user/domain"
	"goshop/internal/user/service"
	"goshop/pkg/apperror"
	"goshop/pkg/response"
	"goshop/pkg/utils"
)

type AddressHandler struct {
	service service.AddressService
}

func NewAddressHandler(svc service.AddressService) *AddressHandler {
	return &AddressHandler{service: svc}
}

func (h *AddressHandler) ListAddresses(c *gin.Context) {
	userID := c.GetString("userId")
	if userID == "" {
		apperror.ErrUnauthorized.HTTPError(c)
		return
	}
	addresses, err := h.service.ListAddresses(c, userID)
	if err != nil {
		logger.Error("Failed to list addresses: ", err)
		apperror.ToHTTPError(c, err, http.StatusInternalServerError, "Something went wrong")
		return
	}
	var res []*domain.Address
	if err := utils.Copy(&res, &addresses); err != nil {
		apperror.ToHTTPError(c, err, http.StatusInternalServerError, "Something went wrong")
		return
	}
	response.JSON(c, http.StatusOK, res)
}

func (h *AddressHandler) GetAddressByID(c *gin.Context) {
	userID := c.GetString("userId")
	if userID == "" {
		apperror.ErrUnauthorized.HTTPError(c)
		return
	}
	id := c.Param("id")
	if id == "" {
		apperror.WrapMessage(apperror.ErrBadRequest, nil, "Missing address ID").HTTPError(c)
		return
	}
	address, err := h.service.GetAddressByID(c, id, userID)
	if err != nil {
		apperror.Wrap(apperror.ErrNotFound, err).HTTPError(c)
		return
	}
	var res domain.Address
	if err := utils.Copy(&res, address); err != nil {
		apperror.ToHTTPError(c, err, http.StatusInternalServerError, "Something went wrong")
		return
	}
	response.JSON(c, http.StatusOK, res)
}

func (h *AddressHandler) CreateAddress(c *gin.Context) {
	userID := c.GetString("userId")
	if userID == "" {
		apperror.ErrUnauthorized.HTTPError(c)
		return
	}
	var req domain.CreateAddressReq
	if err := c.ShouldBindJSON(&req); c.Request.Body == nil || err != nil {
		apperror.Wrap(apperror.ErrBadRequest, err).HTTPError(c)
		return
	}
	address, err := h.service.Create(c, userID, &req)
	if err != nil {
		logger.Error("Failed to create address: ", err)
		apperror.ToHTTPError(c, err, http.StatusInternalServerError, "Something went wrong")
		return
	}
	var res domain.Address
	if err := utils.Copy(&res, address); err != nil {
		apperror.ToHTTPError(c, err, http.StatusInternalServerError, "Something went wrong")
		return
	}
	response.JSON(c, http.StatusOK, res)
}

func (h *AddressHandler) UpdateAddress(c *gin.Context) {
	userID := c.GetString("userId")
	if userID == "" {
		apperror.ErrUnauthorized.HTTPError(c)
		return
	}
	id := c.Param("id")
	var req domain.UpdateAddressReq
	if err := c.ShouldBindJSON(&req); c.Request.Body == nil || err != nil {
		apperror.Wrap(apperror.ErrBadRequest, err).HTTPError(c)
		return
	}
	address, err := h.service.Update(c, id, userID, &req)
	if err != nil {
		apperror.ToHTTPError(c, err, http.StatusInternalServerError, "Something went wrong")
		return
	}
	var res domain.Address
	if err := utils.Copy(&res, address); err != nil {
		apperror.ToHTTPError(c, err, http.StatusInternalServerError, "Something went wrong")
		return
	}
	response.JSON(c, http.StatusOK, res)
}

func (h *AddressHandler) DeleteAddress(c *gin.Context) {
	userID := c.GetString("userId")
	if userID == "" {
		apperror.ErrUnauthorized.HTTPError(c)
		return
	}
	id := c.Param("id")
	if err := h.service.Delete(c, id, userID); err != nil {
		apperror.ToHTTPError(c, err, http.StatusInternalServerError, "Something went wrong")
		return
	}
	response.JSON(c, http.StatusOK, nil)
}

func (h *AddressHandler) SetDefaultAddress(c *gin.Context) {
	userID := c.GetString("userId")
	if userID == "" {
		apperror.ErrUnauthorized.HTTPError(c)
		return
	}
	id := c.Param("id")
	if err := h.service.SetDefault(c, id, userID); err != nil {
		apperror.ToHTTPError(c, err, http.StatusInternalServerError, "Something went wrong")
		return
	}
	response.JSON(c, http.StatusOK, nil)
}
