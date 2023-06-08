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

type Quantity struct {
	service services.IQuantityService
}

func NewQuantityAPI(service services.IQuantityService) *Quantity {
	return &Quantity{service: service}
}

func (q *Quantity) GetQuantities(c *gin.Context) {
	var query serializers.QuantityQueryParam
	if err := c.ShouldBindQuery(&query); err != nil {
		logger.Error("Failed to parse request query: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	quantities, err := q.service.GetQuantities(ctx, &query)
	if err != nil {
		logger.Error("Failed to get quantities}: ", err)
		c.JSON(http.StatusBadRequest, utils.PrepareResponse(nil, err.Error(), ""))
		return
	}

	var res []serializers.Quantity
	copier.Copy(&res, &quantities)
	c.JSON(http.StatusOK, utils.PrepareResponse(res, "OK", ""))
}

func (q *Quantity) GetQuantityByID(c *gin.Context) {
	quantityId := c.Param("uuid")

	ctx := c.Request.Context()
	quantity, err := q.service.GetQuantityByID(ctx, quantityId)
	if err != nil {
		logger.Error("Failed to get quantity: ", err)
		c.JSON(http.StatusBadRequest, utils.PrepareResponse(nil, err.Error(), ""))
		return
	}

	var res serializers.Quantity
	copier.Copy(&res, &quantity)
	c.JSON(http.StatusOK, utils.PrepareResponse(res, "OK", ""))
}

func (q *Quantity) CreateQuantity(c *gin.Context) {
	var item serializers.QuantityBodyParam
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
	quantities, err := q.service.CreateQuantity(ctx, &item)
	if err != nil {
		logger.Error("Failed to create quantity", err.Error())
		c.JSON(http.StatusBadRequest, utils.PrepareResponse(nil, err.Error(), ""))
		return
	}

	var res serializers.Quantity
	copier.Copy(&res, &quantities)
	c.JSON(http.StatusOK, utils.PrepareResponse(res, "OK", ""))
}

func (q *Quantity) UpdateQuantity(c *gin.Context) {
	uuid := c.Param("uuid")
	var item serializers.QuantityBodyParam
	if err := c.ShouldBindJSON(&item); err != nil {
		logger.Error("Failed to parse request body: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	quantities, err := q.service.UpdateQuantity(ctx, uuid, &item)
	if err != nil {
		logger.Error("Failed to update quantity: ", err)
		c.JSON(http.StatusBadRequest, utils.PrepareResponse(nil, err.Error(), ""))
		return
	}

	var res serializers.Quantity
	copier.Copy(&res, &quantities)
	c.JSON(http.StatusOK, utils.PrepareResponse(res, "OK", ""))
}
