package api

import (
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

func (a *OrderAPI) CreateOrder(c *gin.Context) {
	var req serializers.PlaceOrderReq
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("Failed to get body", err)
		response.Error(c, http.StatusBadRequest, err, "Invalid parameters")
		return
	}

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
	var query serializers.OrderQueryParam
	if err := c.ShouldBindQuery(&query); err != nil {
		logger.Error("Failed to parse request query: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	orders, err := a.service.GetOrders(ctx, &query)
	if err != nil {
		logger.Error("Failed to get orders: ", err)
		c.JSON(http.StatusBadRequest, utils.PrepareResponse(nil, err.Error(), ""))
		return
	}

	var res []serializers.Order
	copier.Copy(&res, &orders)
	c.JSON(http.StatusOK, utils.PrepareResponse(res, "OK", ""))
}

func (a *OrderAPI) GetOrderByID(c *gin.Context) {
	orderId := c.Param("uuid")

	ctx := c.Request.Context()
	order, err := a.service.GetOrderByID(ctx, orderId)
	if err != nil {
		logger.Error("Failed to get OrderAPI: ", err)
		c.JSON(http.StatusBadRequest, utils.PrepareResponse(nil, err.Error(), ""))
		return
	}

	var res serializers.Order
	copier.Copy(&res, &order)
	c.JSON(http.StatusOK, utils.PrepareResponse(res, "OK", ""))
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
