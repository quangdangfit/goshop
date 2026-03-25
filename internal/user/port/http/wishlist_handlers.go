package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/quangdangfit/gocommon/logger"

	"goshop/internal/user/dto"
	"goshop/internal/user/service"
	"goshop/pkg/apperror"
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
		apperror.ErrUnauthorized.HTTPError(c)
		return
	}
	items, err := h.service.GetWishlist(c, userID)
	if err != nil {
		logger.Error("Failed to get wishlist: ", err)
		apperror.ToHTTPError(c, err, http.StatusInternalServerError, "Something went wrong")
		return
	}
	var res []*dto.WishlistItem
	utils.Copy(&res, &items)
	response.JSON(c, http.StatusOK, res)
}

func (h *WishlistHandler) AddProduct(c *gin.Context) {
	userID := c.GetString("userId")
	if userID == "" {
		apperror.ErrUnauthorized.HTTPError(c)
		return
	}
	var req dto.AddToWishlistReq
	if err := c.ShouldBindJSON(&req); c.Request.Body == nil || err != nil {
		apperror.Wrap(apperror.ErrBadRequest, err).HTTPError(c)
		return
	}
	if err := h.service.AddProduct(c, userID, &req); err != nil {
		apperror.ToHTTPError(c, err, http.StatusInternalServerError, "Something went wrong")
		return
	}
	response.JSON(c, http.StatusOK, nil)
}

func (h *WishlistHandler) RemoveProduct(c *gin.Context) {
	userID := c.GetString("userId")
	if userID == "" {
		apperror.ErrUnauthorized.HTTPError(c)
		return
	}
	productID := c.Param("productId")
	if productID == "" {
		apperror.WrapMessage(apperror.ErrBadRequest, nil, "Missing product ID").HTTPError(c)
		return
	}
	if err := h.service.RemoveProduct(c, userID, productID); err != nil {
		apperror.ToHTTPError(c, err, http.StatusInternalServerError, "Something went wrong")
		return
	}
	response.JSON(c, http.StatusOK, nil)
}
