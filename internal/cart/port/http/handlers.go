package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/quangdangfit/gocommon/logger"

	"goshop/internal/cart/domain"
	"goshop/internal/cart/service"
	"goshop/pkg/apperror"
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
//	@Success	200	{object}	domain.Cart
//	@Router		/api/v1/cart [get]
func (h *CartHandler) GetCart(c *gin.Context) {
	userID := c.GetString("userId")
	if userID == "" {
		apperror.ErrUnauthorized.HTTPError(c)
		return
	}

	cart, err := h.service.GetCartByUserID(c, userID)
	if err != nil {
		logger.Error("Failed to get cart: ", err)
		apperror.ToHTTPError(c, err, http.StatusInternalServerError, "Something went wrong")
		return
	}

	var res domain.Cart
	utils.Copy(&res, &cart)
	response.JSON(c, http.StatusOK, res)
}

// AddProduct godoc
//
//	@Summary	add product to cart
//	@Tags		cart
//	@Produce	json
//	@Security	ApiKeyAuth
//	@Param		_	body		domain.CartLineReq	true	"Body"
//	@Success	200	{object}	domain.Cart
//	@Router		/api/v1/cart [post]
func (h *CartHandler) AddProduct(c *gin.Context) {
	userID := c.GetString("userId")
	if userID == "" {
		apperror.ErrUnauthorized.HTTPError(c)
		return
	}

	var line domain.CartLineReq
	if err := c.ShouldBindJSON(&line); err != nil {
		logger.Error("Failed to get body", err)
		apperror.Wrap(apperror.ErrBadRequest, err).HTTPError(c)
		return
	}

	cart, err := h.service.AddProduct(c, &domain.AddProductReq{
		UserID: userID,
		Line:   &line,
	})
	if err != nil {
		logger.Error("Failed to add product to cart: ", err)
		apperror.ToHTTPError(c, err, http.StatusInternalServerError, "Something went wrong")
		return
	}

	var res domain.Cart
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
		apperror.ErrUnauthorized.HTTPError(c)
		return
	}

	productID := c.Param("productId")
	if productID == "" {
		apperror.WrapMessage(apperror.ErrBadRequest, nil, "Missing product ID").HTTPError(c)
		return
	}

	cart, err := h.service.RemoveProduct(c, &domain.RemoveProductReq{
		UserID:    userID,
		ProductID: productID,
	})
	if err != nil {
		logger.Error("Failed to remove product from cart: ", err)
		apperror.ToHTTPError(c, err, http.StatusInternalServerError, "Something went wrong")
		return
	}

	var res domain.Cart
	utils.Copy(&res, &cart)
	response.JSON(c, http.StatusOK, res)
}
