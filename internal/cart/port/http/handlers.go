package http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/quangdangfit/gocommon/logger"

	"goshop/internal/cart/dto"
	"goshop/internal/cart/service"
	"goshop/pkg/response"
	"goshop/pkg/utils"
)

type CartHandler struct {
	service service.CartService
}

func NewCartHandler(service service.CartService) *CartHandler {
	return &CartHandler{service: service}
}

// GetCart godoc
//
//	@Summary	get my cart
//	@Tags		cart
//	@Produce	json
//	@Security	ApiKeyAuth
//	@Success	200	{object}	dto.Cart
//	@Router		/api/v1/cart [get]
func (h *CartHandler) GetCart(c *gin.Context) {
	userID := c.GetString("userId")
	if userID == "" {
		response.Error(c, http.StatusUnauthorized, errors.New("unauthorized"), "Unauthorized")
		return
	}

	cart, err := h.service.GetCartByUserID(c, userID)
	if err != nil {
		logger.Error("Failed to get cart: ", err)
		response.Error(c, http.StatusInternalServerError, err, "Something went wrong")
		return
	}

	var res dto.Cart
	utils.Copy(&res, &cart)
	response.JSON(c, http.StatusOK, res)
}

// AddProduct godoc
//
//	@Summary	add product to cart
//	@Tags		cart
//	@Produce	json
//	@Security	ApiKeyAuth
//	@Param		_	body		dto.CartLineReq	true	"Body"
//	@Success	200	{object}	dto.Cart
//	@Router		/api/v1/cart [post]
func (h *CartHandler) AddProduct(c *gin.Context) {
	userID := c.GetString("userId")
	if userID == "" {
		response.Error(c, http.StatusUnauthorized, errors.New("unauthorized"), "Unauthorized")
		return
	}

	var line dto.CartLineReq
	if err := c.ShouldBindJSON(&line); err != nil {
		logger.Error("Failed to get body", err)
		response.Error(c, http.StatusBadRequest, err, "Invalid parameters")
		return
	}

	cart, err := h.service.AddProduct(c, &dto.AddProductReq{
		UserID: userID,
		Line:   &line,
	})
	if err != nil {
		logger.Error("Failed to add product to cart: ", err)
		response.Error(c, http.StatusInternalServerError, err, "Something went wrong")
		return
	}

	var res dto.Cart
	utils.Copy(&res, &cart)
	response.JSON(c, http.StatusOK, res)
}

// RemoveProduct godoc
//
//	@Summary	remove product from cart
//	@Tags		cart
//	@Produce	json
//	@Security	ApiKeyAuth
//	@Param		productId	path	string	true	"Product ID"
//	@Router		/api/v1/cart/{productId} [delete]
func (h *CartHandler) RemoveProduct(c *gin.Context) {
	userID := c.GetString("userId")
	if userID == "" {
		response.Error(c, http.StatusUnauthorized, errors.New("unauthorized"), "Unauthorized")
		return
	}

	productID := c.Param("productId")
	if productID == "" {
		response.Error(c, http.StatusBadRequest, errors.New("bad request"), "Miss Product ID")
		return
	}

	cart, err := h.service.RemoveProduct(c, &dto.RemoveProductReq{
		UserID:    userID,
		ProductID: productID,
	})
	if err != nil {
		logger.Error("Failed to remove product from cart: ", err)
		response.Error(c, http.StatusInternalServerError, err, "Something went wrong")
		return
	}

	var res dto.Cart
	utils.Copy(&res, &cart)
	response.JSON(c, http.StatusOK, res)
}
