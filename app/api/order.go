package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/quangdangfit/gocommon/logger"
	"github.com/quangdangfit/gocommon/validation"

	"goshop/app/serializers"
	"goshop/app/services"
	"goshop/pkg/response"
	"goshop/pkg/utils"
)

type OrderAPI struct {
	validator validation.Validation
	service   services.IOrderService
}

func NewOrderAPI(validator validation.Validation, service services.IOrderService) *OrderAPI {
	return &OrderAPI{
		validator: validator,
		service:   service,
	}
}

// PlaceOrder godoc
//
//	@Summary	place order
//	@Produce	json
//	@Security	ApiKeyAuth
//	@Param		_	body		serializers.PlaceOrderReq	true	"Body"
//	@Success	200	{object}	serializers.Order
//	@Router		/api/v1/orders [post]
func (a *OrderAPI) PlaceOrder(c *gin.Context) {
	var req serializers.PlaceOrderReq
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("Failed to get body", err)
		response.Error(c, http.StatusBadRequest, err, "Invalid parameters")
		return
	}

	req.UserID = c.GetString("userId")
	if err := a.validator.ValidateStruct(req); err != nil {
		response.Error(c, http.StatusBadRequest, err, "Invalid parameters")
		return
	}

	order, err := a.service.PlaceOrder(c, &req)
	if err != nil {
		logger.Error("Failed to create OrderAPI: ", err.Error())
		response.Error(c, http.StatusInternalServerError, err, "Something went wrong")
		return
	}

	var res serializers.Order
	utils.Copy(&res, &order)
	response.JSON(c, http.StatusOK, res)
}

// GetOrders godoc
//
//	@Summary	get my orders
//	@Produce	json
//	@Security	ApiKeyAuth
//	@Param		_	query		serializers.ListOrderReq	true	"Query"
//	@Success	200	{object}	serializers.ListOrderRes
//	@Router		/api/v1/orders [get]
func (a *OrderAPI) GetOrders(c *gin.Context) {
	var req serializers.ListOrderReq
	if err := c.ShouldBindQuery(&req); err != nil {
		logger.Error("Failed to parse request req: ", err)
		response.Error(c, http.StatusBadRequest, err, "Invalid parameters")
		return
	}

	req.UserID = c.GetString("userId")
	orders, pagination, err := a.service.GetMyOrders(c, &req)
	if err != nil {
		logger.Error("Failed to get orders: ", err)
		response.Error(c, http.StatusInternalServerError, err, "Something went wrong")
		return
	}

	var res serializers.ListOrderRes
	res.Pagination = pagination
	utils.Copy(&res.Orders, &orders)
	response.JSON(c, http.StatusOK, res)
}

// GetOrderByID godoc
//
//	@Summary	get order details
//	@Produce	json
//	@Security	ApiKeyAuth
//	@Param		id	path		string	true	"Order ID"
//	@Success	200	{object}	serializers.Order
//	@Router		/api/v1/orders/{id} [get]
func (a *OrderAPI) GetOrderByID(c *gin.Context) {
	orderId := c.Param("id")
	order, err := a.service.GetOrderByID(c, orderId)
	if err != nil {
		logger.Errorf("Failed to get order, id: %s, error: %s ", orderId, err)
		response.Error(c, http.StatusNotFound, err, "Not found")
		return
	}

	var res serializers.Order
	utils.Copy(&res, &order)
	response.JSON(c, http.StatusOK, res)
}

// CancelOrder godoc
//
//	@Summary	cancel order
//	@Produce	json
//	@Security	ApiKeyAuth
//	@Param		id	path	string	true	"Order ID"
//	@Router		/api/v1/orders/{id}/cancel [put]
func (a *OrderAPI) CancelOrder(c *gin.Context) {
	orderID := c.Param("id")
	userID := c.GetString("userId")
	order, err := a.service.CancelOrder(c, orderID, userID)
	if err != nil {
		logger.Errorf("Failed to cancel order, id: %s, error: %s", orderID, err)
		response.Error(c, http.StatusInternalServerError, err, "Something went wrong")
		return
	}

	var res serializers.Order
	utils.Copy(&res, &order)
	response.JSON(c, http.StatusOK, nil)
}
