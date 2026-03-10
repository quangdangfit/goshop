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

type WishlistHandler struct {
	service service.WishlistService
}

func NewWishlistHandler(svc service.WishlistService) *WishlistHandler {
	return &WishlistHandler{service: svc}
}

func (h *WishlistHandler) GetWishlist(c *gin.Context) {
	userID := c.GetString("userId")
	if userID == "" {
		response.Error(c, http.StatusUnauthorized, errors.New("unauthorized"), "Unauthorized")
		return
	}
	items, err := h.service.GetWishlist(c, userID)
	if err != nil {
		logger.Error("Failed to get wishlist: ", err)
		response.Error(c, http.StatusInternalServerError, err, "Something went wrong")
		return
	}
	var res []*dto.WishlistItem
	utils.Copy(&res, &items)
	response.JSON(c, http.StatusOK, res)
}

func (h *WishlistHandler) AddProduct(c *gin.Context) {
	userID := c.GetString("userId")
	if userID == "" {
		response.Error(c, http.StatusUnauthorized, errors.New("unauthorized"), "Unauthorized")
		return
	}
	var req dto.AddToWishlistReq
	if err := c.ShouldBindJSON(&req); c.Request.Body == nil || err != nil {
		response.Error(c, http.StatusBadRequest, err, "Invalid parameters")
		return
	}
	if err := h.service.AddProduct(c, userID, &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err, "Something went wrong")
		return
	}
	response.JSON(c, http.StatusOK, nil)
}

func (h *WishlistHandler) RemoveProduct(c *gin.Context) {
	userID := c.GetString("userId")
	if userID == "" {
		response.Error(c, http.StatusUnauthorized, errors.New("unauthorized"), "Unauthorized")
		return
	}
	productID := c.Param("productId")
	if productID == "" {
		response.Error(c, http.StatusBadRequest, errors.New("missing productId"), "Missing product ID")
		return
	}
	if err := h.service.RemoveProduct(c, userID, productID); err != nil {
		response.Error(c, http.StatusInternalServerError, err, "Something went wrong")
		return
	}
	response.JSON(c, http.StatusOK, nil)
}
