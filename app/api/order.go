package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jinzhu/copier"
	"github.com/quangdangfit/gocommon/logger"

	"goshop/app/serializers"
	"goshop/app/services"
	"goshop/pkg/utils"
)

type Order struct {
	service services.IOrderSerivce
}

func NewOrderAPI(service services.IOrderSerivce) *Order {
	return &Order{service: service}
}

func (categ *Order) GetOrders(c *gin.Context) {
	var query serializers.OrderQueryParam
	if err := c.ShouldBindQuery(&query); err != nil {
		logger.Error("Failed to parse request query: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	orders, err := categ.service.GetOrders(ctx, &query)
	if err != nil {
		logger.Error("Failed to get orders: ", err)
		c.JSON(http.StatusBadRequest, utils.PrepareResponse(nil, err.Error(), ""))
		return
	}

	var res []serializers.Order
	copier.Copy(&res, &orders)
	c.JSON(http.StatusOK, utils.PrepareResponse(res, "OK", ""))
}

func (categ *Order) GetOrderByID(c *gin.Context) {
	orderId := c.Param("uuid")

	ctx := c.Request.Context()
	order, err := categ.service.GetOrderByID(ctx, orderId)
	if err != nil {
		logger.Error("Failed to get Order: ", err)
		c.JSON(http.StatusBadRequest, utils.PrepareResponse(nil, err.Error(), ""))
		return
	}

	var res serializers.Order
	copier.Copy(&res, &order)
	c.JSON(http.StatusOK, utils.PrepareResponse(res, "OK", ""))
}

func (categ *Order) CreateOrder(c *gin.Context) {
	var item serializers.OrderBodyParam
	if err := c.Bind(&item); err != nil {
		logger.Error("Failed to parse request body: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	validate := validator.New()
	err := validate.Struct(item)
	if err != nil {
		logger.Error("Request body is invalid: ", err.Error())
		c.JSON(http.StatusBadRequest, utils.PrepareResponse(nil, err.Error(), ""))
		return
	}

	ctx := c.Request.Context()
	orders, err := categ.service.CreateOrder(ctx, &item)
	if err != nil {
		logger.Error("Failed to create Order: ", err.Error())
		c.JSON(http.StatusBadRequest, utils.PrepareResponse(nil, err.Error(), ""))
		return
	}

	var res serializers.Order
	copier.Copy(&res, &orders)
	c.JSON(http.StatusOK, utils.PrepareResponse(res, "OK", ""))
}

func (categ *Order) UpdateOrder(c *gin.Context) {
	uuid := c.Param("uuid")
	var item serializers.OrderBodyParam
	if err := c.ShouldBindJSON(&item); err != nil {
		logger.Error("Failed to parse request body: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	orders, err := categ.service.UpdateOrder(ctx, uuid, &item)
	if err != nil {
		logger.Error("Failed to update Order: ", err)
		c.JSON(http.StatusBadRequest, utils.PrepareResponse(nil, err.Error(), ""))
		return
	}

	var res serializers.Order
	copier.Copy(&res, &orders)
	c.JSON(http.StatusOK, utils.PrepareResponse(res, "OK", ""))
}
