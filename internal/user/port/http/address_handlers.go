package http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/quangdangfit/gocommon/logger"

	"goshop/internal/user/dto"
	"goshop/internal/user/service"
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
		response.Error(c, http.StatusUnauthorized, errors.New("unauthorized"), "Unauthorized")
		return
	}
	addresses, err := h.service.ListAddresses(c, userID)
	if err != nil {
		logger.Error("Failed to list addresses: ", err)
		response.Error(c, http.StatusInternalServerError, err, "Something went wrong")
		return
	}
	var res []*dto.Address
	utils.Copy(&res, &addresses)
	response.JSON(c, http.StatusOK, res)
}

func (h *AddressHandler) GetAddressByID(c *gin.Context) {
	userID := c.GetString("userId")
	if userID == "" {
		response.Error(c, http.StatusUnauthorized, errors.New("unauthorized"), "Unauthorized")
		return
	}
	id := c.Param("id")
	if id == "" {
		response.Error(c, http.StatusBadRequest, errors.New("missing id"), "Missing address ID")
		return
	}
	address, err := h.service.GetAddressByID(c, id, userID)
	if err != nil {
		response.Error(c, http.StatusNotFound, err, "Not found")
		return
	}
	var res dto.Address
	utils.Copy(&res, address)
	response.JSON(c, http.StatusOK, res)
}

func (h *AddressHandler) CreateAddress(c *gin.Context) {
	userID := c.GetString("userId")
	if userID == "" {
		response.Error(c, http.StatusUnauthorized, errors.New("unauthorized"), "Unauthorized")
		return
	}
	var req dto.CreateAddressReq
	if err := c.ShouldBindJSON(&req); c.Request.Body == nil || err != nil {
		response.Error(c, http.StatusBadRequest, err, "Invalid parameters")
		return
	}
	address, err := h.service.Create(c, userID, &req)
	if err != nil {
		logger.Error("Failed to create address: ", err)
		response.Error(c, http.StatusInternalServerError, err, "Something went wrong")
		return
	}
	var res dto.Address
	utils.Copy(&res, address)
	response.JSON(c, http.StatusOK, res)
}

func (h *AddressHandler) UpdateAddress(c *gin.Context) {
	userID := c.GetString("userId")
	if userID == "" {
		response.Error(c, http.StatusUnauthorized, errors.New("unauthorized"), "Unauthorized")
		return
	}
	id := c.Param("id")
	var req dto.UpdateAddressReq
	if err := c.ShouldBindJSON(&req); c.Request.Body == nil || err != nil {
		response.Error(c, http.StatusBadRequest, err, "Invalid parameters")
		return
	}
	address, err := h.service.Update(c, id, userID, &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err, "Something went wrong")
		return
	}
	var res dto.Address
	utils.Copy(&res, address)
	response.JSON(c, http.StatusOK, res)
}

func (h *AddressHandler) DeleteAddress(c *gin.Context) {
	userID := c.GetString("userId")
	if userID == "" {
		response.Error(c, http.StatusUnauthorized, errors.New("unauthorized"), "Unauthorized")
		return
	}
	id := c.Param("id")
	if err := h.service.Delete(c, id, userID); err != nil {
		response.Error(c, http.StatusInternalServerError, err, "Something went wrong")
		return
	}
	response.JSON(c, http.StatusOK, nil)
}

func (h *AddressHandler) SetDefaultAddress(c *gin.Context) {
	userID := c.GetString("userId")
	if userID == "" {
		response.Error(c, http.StatusUnauthorized, errors.New("unauthorized"), "Unauthorized")
		return
	}
	id := c.Param("id")
	if err := h.service.SetDefault(c, id, userID); err != nil {
		response.Error(c, http.StatusInternalServerError, err, "Something went wrong")
		return
	}
	response.JSON(c, http.StatusOK, nil)
}
