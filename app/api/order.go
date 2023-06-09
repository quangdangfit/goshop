package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/quangdangfit/gocommon/logger"

	"goshop/app/serializers"
	"goshop/app/services"
	"goshop/pkg/response"
	"goshop/pkg/utils"
	"goshop/pkg/validation"
)

type OrderAPI struct {
	validator validation.Validation
	service   services.IOrderService
}

func NewOrderAPI(service services.IOrderService) *OrderAPI {
	return &OrderAPI{
		validator: validation.New(),
		service:   service,
	}
}

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
	err = copier.Copy(&res, &order)
	if err != nil {
		logger.Error(err.Error())
		response.Error(c, http.StatusInternalServerError, err, "Something went wrong")
		return
	}
	response.JSON(c, http.StatusOK, res)
}

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
	err = copier.Copy(&res.Orders, &orders)
	if err != nil {
		logger.Error(err.Error())
		response.Error(c, http.StatusInternalServerError, err, "Something went wrong")
		return
	}

	response.JSON(c, http.StatusOK, res)
}

func (a *OrderAPI) GetOrderByID(c *gin.Context) {
	orderId := c.Param("id")
	if orderId == "" {
		response.Error(c, http.StatusBadRequest, errors.New("missing id"), "Invalid Parameters")
		return
	}

	ctx := c.Request.Context()
	order, err := a.service.GetOrderByID(ctx, orderId)
	if err != nil {
		logger.Errorf("Failed to get order, id: %s, error: %s ", orderId, err)
		response.Error(c, http.StatusNotFound, err, "Not found")
		return
	}

	var res serializers.Order
	err = copier.Copy(&res, &order)
	if err != nil {
		logger.Error(err.Error())
		response.Error(c, http.StatusInternalServerError, err, "Something went wrong")
		return
	}
	response.JSON(c, http.StatusOK, res)
}

func (a *OrderAPI) UpdateOrder(c *gin.Context) {
	uuid := c.Param("uuid")
	var req serializers.PlaceOrderReq
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("Failed to parse request body: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	orders, err := a.service.UpdateOrder(ctx, uuid, &req)
	if err != nil {
		logger.Error("Failed to update OrderAPI: ", err)
		c.JSON(http.StatusBadRequest, utils.PrepareResponse(nil, err.Error(), ""))
		return
	}

	var res serializers.Order
	copier.Copy(&res, &orders)
	c.JSON(http.StatusOK, utils.PrepareResponse(res, "OK", ""))
}
