package http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/quangdangfit/gocommon/logger"

	"goshop/internal/order/dto"
	"goshop/internal/order/service"
	"goshop/pkg/response"
	"goshop/pkg/utils"
)

type OrderHandler struct {
	service service.IOrderService
}

func NewOrderHandler(service service.IOrderService) *OrderHandler {
	return &OrderHandler{
		service: service,
	}
}

// PlaceOrder godoc
//
//	@Summary	place order
//	@Tags		orders
//	@Produce	json
//	@Security	ApiKeyAuth
//	@Param		_	body		dto.PlaceOrderReq	true	"Body"
//	@Success	200	{object}	dto.Order
//	@Router		/api/v1/orders [post]
func (a *OrderHandler) PlaceOrder(c *gin.Context) {
	var req dto.PlaceOrderReq
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("Failed to get body", err)
		response.Error(c, http.StatusBadRequest, err, "Invalid parameters")
		return
	}

	req.UserID = c.GetString("userId")
	if req.UserID == "" {
		response.Error(c, http.StatusUnauthorized, errors.New("unauthorized"), "Unauthorized")
		return
	}

	order, err := a.service.PlaceOrder(c, &req)
	if err != nil {
		logger.Error("Failed to create OrderHandler: ", err.Error())
		response.Error(c, http.StatusInternalServerError, err, "Something went wrong")
		return
	}

	var res dto.Order
	utils.Copy(&res, &order)
	response.JSON(c, http.StatusOK, res)
}

// GetOrders godoc
//
//	@Summary	get my orders
//	@Tags		orders
//	@Produce	json
//	@Security	ApiKeyAuth
//	@Param		_	query		dto.ListOrderReq	true	"Query"
//	@Success	200	{object}	dto.ListOrderRes
//	@Router		/api/v1/orders [get]
func (a *OrderHandler) GetOrders(c *gin.Context) {
	var req dto.ListOrderReq
	if err := c.ShouldBindQuery(&req); err != nil {
		logger.Error("Failed to parse request req: ", err)
		response.Error(c, http.StatusBadRequest, err, "Invalid parameters")
		return
	}

	req.UserID = c.GetString("userId")
	if req.UserID == "" {
		response.Error(c, http.StatusUnauthorized, errors.New("unauthorized"), "Unauthorized")
		return
	}

	orders, pagination, err := a.service.GetMyOrders(c, &req)
	if err != nil {
		logger.Error("Failed to get orders: ", err)
		response.Error(c, http.StatusInternalServerError, err, "Something went wrong")
		return
	}

	var res dto.ListOrderRes
	res.Pagination = pagination
	utils.Copy(&res.Orders, &orders)
	response.JSON(c, http.StatusOK, res)
}

// GetOrderByID godoc
//
//	@Summary	get order details
//	@Tags		orders
//	@Produce	json
//	@Security	ApiKeyAuth
//	@Param		id	path		string	true	"Order ID"
//	@Success	200	{object}	dto.Order
//	@Router		/api/v1/orders/{id} [get]
func (a *OrderHandler) GetOrderByID(c *gin.Context) {
	userId := c.GetString("userId")
	if userId == "" {
		response.Error(c, http.StatusUnauthorized, errors.New("unauthorized"), "Unauthorized")
		return
	}

	orderId := c.Param("id")
	if orderId == "" {
		response.Error(c, http.StatusBadRequest, errors.New("bad request"), "Miss Order ID")
		return
	}

	order, err := a.service.GetOrderByID(c, orderId)
	if err != nil {
		logger.Errorf("Failed to get order, id: %s, error: %s ", orderId, err)
		response.Error(c, http.StatusNotFound, err, "Not found")
		return
	}

	var res dto.Order
	utils.Copy(&res, &order)
	response.JSON(c, http.StatusOK, res)
}

// CancelOrder godoc
//
//	@Summary	cancel order
//	@Tags		orders
//	@Produce	json
//	@Security	ApiKeyAuth
//	@Param		id	path	string	true	"Order ID"
//	@Router		/api/v1/orders/{id}/cancel [put]
func (a *OrderHandler) CancelOrder(c *gin.Context) {
	userID := c.GetString("userId")
	if userID == "" {
		response.Error(c, http.StatusUnauthorized, errors.New("unauthorized"), "Unauthorized")
		return
	}

	orderID := c.Param("id")
	if orderID == "" {
		response.Error(c, http.StatusBadRequest, errors.New("bad request"), "Miss Order ID")
		return
	}

	order, err := a.service.CancelOrder(c, orderID, userID)
	if err != nil {
		logger.Errorf("Failed to cancel order, id: %s, error: %s", orderID, err)
		response.Error(c, http.StatusInternalServerError, err, "Something went wrong")
		return
	}

	var res dto.Order
	utils.Copy(&res, &order)
	response.JSON(c, http.StatusOK, res)
}
