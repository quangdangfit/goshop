package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/quangdangfit/gocommon/logger"

	"goshop/internal/order/dto"
	"goshop/internal/order/model"
	"goshop/internal/order/service"
	"goshop/pkg/apperror"
	"goshop/pkg/response"
	"goshop/pkg/utils"
)

type OrderHandler struct {
	service service.OrderService
}

func NewOrderHandler(service service.OrderService) *OrderHandler {
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
		apperror.Wrap(apperror.ErrBadRequest, err).HTTPError(c)
		return
	}

	req.UserID = c.GetString("userId")
	if req.UserID == "" {
		apperror.ErrUnauthorized.HTTPError(c)
		return
	}

	order, err := a.service.PlaceOrder(c, &req)
	if err != nil {
		logger.Error("Failed to create OrderHandler: ", err.Error())
		apperror.ToHTTPError(c, err, http.StatusInternalServerError, "Something went wrong")
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
		apperror.Wrap(apperror.ErrBadRequest, err).HTTPError(c)
		return
	}

	req.UserID = c.GetString("userId")
	if req.UserID == "" {
		apperror.ErrUnauthorized.HTTPError(c)
		return
	}

	orders, pagination, err := a.service.GetMyOrders(c, &req)
	if err != nil {
		logger.Error("Failed to get orders: ", err)
		apperror.ToHTTPError(c, err, http.StatusInternalServerError, "Something went wrong")
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
		apperror.ErrUnauthorized.HTTPError(c)
		return
	}

	orderId := c.Param("id")
	if orderId == "" {
		apperror.WrapMessage(apperror.ErrBadRequest, nil, "Missing order ID").HTTPError(c)
		return
	}

	order, err := a.service.GetOrderByID(c, orderId)
	if err != nil {
		logger.Errorf("Failed to get order, id: %s, error: %s ", orderId, err)
		apperror.Wrap(apperror.ErrNotFound, err).HTTPError(c)
		return
	}

	var res dto.Order
	utils.Copy(&res, &order)
	response.JSON(c, http.StatusOK, res)
}

// UpdateOrderStatus godoc
//
//	@Summary	update order status (admin)
//	@Tags		orders
//	@Produce	json
//	@Security	ApiKeyAuth
//	@Param		id		path	string	true	"Order ID"
//	@Param		status	query	string	true	"New status"
//	@Router		/api/v1/orders/{id}/status [put]
func (a *OrderHandler) UpdateOrderStatus(c *gin.Context) {
	orderID := c.Param("id")
	if orderID == "" {
		apperror.WrapMessage(apperror.ErrBadRequest, nil, "Missing order ID").HTTPError(c)
		return
	}

	status := c.Query("status")
	if status == "" {
		apperror.WrapMessage(apperror.ErrBadRequest, nil, "Missing status").HTTPError(c)
		return
	}

	order, err := a.service.UpdateOrderStatus(c, orderID, model.OrderStatus(status))
	if err != nil {
		logger.Errorf("Failed to update order status, id: %s, error: %s", orderID, err)
		apperror.ToHTTPError(c, err, http.StatusInternalServerError, "Something went wrong")
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
		apperror.ErrUnauthorized.HTTPError(c)
		return
	}

	orderID := c.Param("id")
	if orderID == "" {
		apperror.WrapMessage(apperror.ErrBadRequest, nil, "Missing order ID").HTTPError(c)
		return
	}

	order, err := a.service.CancelOrder(c, orderID, userID)
	if err != nil {
		logger.Errorf("Failed to cancel order, id: %s, error: %s", orderID, err)
		apperror.ToHTTPError(c, err, http.StatusInternalServerError, "Something went wrong")
		return
	}

	var res dto.Order
	utils.Copy(&res, &order)
	response.JSON(c, http.StatusOK, res)
}
