package services

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jinzhu/copier"
	"gitlab.com/quangdangfit/gocommon/utils/logger"

	"goshop/models"
	"goshop/repositories"
	"goshop/utils"
)

type WarehouseSerivce interface {
	GetWarehouses(c *gin.Context)
	GetWarehouseByID(c *gin.Context)
	CreateWarehouse(c *gin.Context)
	UpdateWarehouse(c *gin.Context)
}

type warehouse struct {
	repo repositories.WarehouseRepository
}

func NewWarehouseService(repo repositories.WarehouseRepository) WarehouseSerivce {
	return &warehouse{repo: repo}
}

func (categ *warehouse) GetWarehouses(c *gin.Context) {
	var reqQuery models.WarehouseQueryRequest
	if err := c.ShouldBindQuery(&reqQuery); err != nil {
		logger.Error("Failed to parse request query: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var query map[string]interface{}
	data, _ := json.Marshal(reqQuery)
	json.Unmarshal(data, &query)
	warehouses, err := categ.repo.GetWarehouses(query)
	if err != nil {
		logger.Error("Failed to get warehouses}: ", err)
		c.JSON(http.StatusBadRequest, utils.PrepareResponse(nil, err.Error(), ""))
		return
	}

	var res []models.WarehouseResponse
	copier.Copy(&res, &warehouses)
	c.JSON(http.StatusOK, utils.PrepareResponse(res, "OK", ""))
}

func (categ *warehouse) GetWarehouseByID(c *gin.Context) {
	warehouseId := c.Param("uuid")

	warehouse, err := categ.repo.GetWarehouseByID(warehouseId)
	if err != nil {
		logger.Error("Failed to get warehouse: ", err)
		c.JSON(http.StatusBadRequest, utils.PrepareResponse(nil, err.Error(), ""))
		return
	}

	var res models.WarehouseResponse
	copier.Copy(&res, &warehouse)
	c.JSON(http.StatusOK, utils.PrepareResponse(res, "OK", ""))
}

func (categ *warehouse) CreateWarehouse(c *gin.Context) {
	var reqBody models.WarehouseBodyRequest
	if err := c.Bind(&reqBody); err != nil {
		logger.Error("Failed to parse request body: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	validate := validator.New()
	err := validate.Struct(reqBody)
	if err != nil {
		logger.Error("Request body is invalid: ", err.Error())
		c.JSON(http.StatusBadRequest, utils.PrepareResponse(nil, err.Error(), ""))
		return
	}

	warehouses, err := categ.repo.CreateWarehouse(&reqBody)
	if err != nil {
		logger.Error("Failed to create warehouse", err.Error())
		c.JSON(http.StatusBadRequest, utils.PrepareResponse(nil, err.Error(), ""))
		return
	}

	var res models.WarehouseResponse
	copier.Copy(&res, &warehouses)
	c.JSON(http.StatusOK, utils.PrepareResponse(res, "OK", ""))
}

func (categ *warehouse) UpdateWarehouse(c *gin.Context) {
	uuid := c.Param("uuid")
	var reqBody models.WarehouseBodyRequest
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		logger.Error("Failed to parse request body: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	warehouses, err := categ.repo.UpdateWarehouse(uuid, &reqBody)
	if err != nil {
		logger.Error("Failed to update warehouse: ", err)
		c.JSON(http.StatusBadRequest, utils.PrepareResponse(nil, err.Error(), ""))
		return
	}

	var res models.WarehouseResponse
	copier.Copy(&res, &warehouses)
	c.JSON(http.StatusOK, utils.PrepareResponse(res, "OK", ""))
}
