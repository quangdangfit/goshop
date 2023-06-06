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

type Warehouse struct {
	service services.IWarehouseSerivce
}

func NewWarehouseAPI(service services.IWarehouseSerivce) *Warehouse {
	return &Warehouse{service: service}
}

func (w *Warehouse) GetWarehouses(c *gin.Context) {
	var query serializers.WarehouseQueryParam
	if err := c.ShouldBindQuery(&query); err != nil {
		logger.Error("Failed to parse request query: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	warehouses, err := w.service.GetWarehouses(ctx, &query)
	if err != nil {
		logger.Error("Failed to get warehouses}: ", err)
		c.JSON(http.StatusBadRequest, utils.PrepareResponse(nil, err.Error(), ""))
		return
	}

	var res []serializers.Warehouse
	copier.Copy(&res, &warehouses)
	c.JSON(http.StatusOK, utils.PrepareResponse(res, "OK", ""))
}

func (w *Warehouse) GetWarehouseByID(c *gin.Context) {
	warehouseId := c.Param("uuid")

	ctx := c.Request.Context()
	warehouse, err := w.service.GetWarehouseByID(ctx, warehouseId)
	if err != nil {
		logger.Error("Failed to get warehouse: ", err)
		c.JSON(http.StatusBadRequest, utils.PrepareResponse(nil, err.Error(), ""))
		return
	}

	var res serializers.Warehouse
	copier.Copy(&res, &warehouse)
	c.JSON(http.StatusOK, utils.PrepareResponse(res, "OK", ""))
}

func (w *Warehouse) CreateWarehouse(c *gin.Context) {
	var item serializers.WarehouseBodyParam
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
	warehouses, err := w.service.CreateWarehouse(ctx, &item)
	if err != nil {
		logger.Error("Failed to create warehouse", err.Error())
		c.JSON(http.StatusBadRequest, utils.PrepareResponse(nil, err.Error(), ""))
		return
	}

	var res serializers.Warehouse
	copier.Copy(&res, &warehouses)
	c.JSON(http.StatusOK, utils.PrepareResponse(res, "OK", ""))
}

func (w *Warehouse) UpdateWarehouse(c *gin.Context) {
	uuid := c.Param("uuid")
	var item serializers.WarehouseBodyParam
	if err := c.ShouldBindJSON(&item); err != nil {
		logger.Error("Failed to parse request body: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	warehouses, err := w.service.UpdateWarehouse(ctx, uuid, &item)
	if err != nil {
		logger.Error("Failed to update warehouse: ", err)
		c.JSON(http.StatusBadRequest, utils.PrepareResponse(nil, err.Error(), ""))
		return
	}

	var res serializers.Warehouse
	copier.Copy(&res, &warehouses)
	c.JSON(http.StatusOK, utils.PrepareResponse(res, "OK", ""))
}
