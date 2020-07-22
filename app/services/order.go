package services

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jinzhu/copier"
	"gitlab.com/quangdangfit/gocommon/utils/logger"

	"goshop/app/models"
	"goshop/app/repositories"
	"goshop/pkg/utils"
)

type OrderSerivce interface {
	GetOrders(c *gin.Context)
	GetOrderByID(c *gin.Context)
	CreateOrder(c *gin.Context)
	UpdateOrder(c *gin.Context)
}

type order struct {
	repo repositories.OrderRepository
}

func NewOrderService(repo repositories.OrderRepository) OrderSerivce {
	return &order{repo: repo}
}

func (categ *order) GetOrders(c *gin.Context) {
	var reqQuery models.OrderQueryRequest
	if err := c.ShouldBindQuery(&reqQuery); err != nil {
		logger.Error("Failed to parse request query: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var query map[string]interface{}
	data, _ := json.Marshal(reqQuery)
	json.Unmarshal(data, &query)
	orders, err := categ.repo.GetOrders(query)
	if err != nil {
		logger.Error("Failed to get orders: ", err)
		c.JSON(http.StatusBadRequest, utils.PrepareResponse(nil, err.Error(), ""))
		return
	}

	var res []models.OrderResponse
	copier.Copy(&res, &orders)
	c.JSON(http.StatusOK, utils.PrepareResponse(res, "OK", ""))
}

func (categ *order) GetOrderByID(c *gin.Context) {
	orderId := c.Param("uuid")

	order, err := categ.repo.GetOrderByID(orderId)
	if err != nil {
		logger.Error("Failed to get order: ", err)
		c.JSON(http.StatusBadRequest, utils.PrepareResponse(nil, err.Error(), ""))
		return
	}

	var res models.OrderResponse
	copier.Copy(&res, &order)
	c.JSON(http.StatusOK, utils.PrepareResponse(res, "OK", ""))
}

func (categ *order) CreateOrder(c *gin.Context) {
	var reqBody models.OrderBodyRequest
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

	orders, err := categ.repo.CreateOrder(&reqBody)
	if err != nil {
		logger.Error("Failed to create order: ", err.Error())
		c.JSON(http.StatusBadRequest, utils.PrepareResponse(nil, err.Error(), ""))
		return
	}

	var res models.OrderResponse
	copier.Copy(&res, &orders)
	c.JSON(http.StatusOK, utils.PrepareResponse(res, "OK", ""))
}

func (categ *order) UpdateOrder(c *gin.Context) {
	uuid := c.Param("uuid")
	var reqBody models.OrderBodyRequest
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		logger.Error("Failed to parse request body: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	orders, err := categ.repo.UpdateOrder(uuid, &reqBody)
	if err != nil {
		logger.Error("Failed to update order: ", err)
		c.JSON(http.StatusBadRequest, utils.PrepareResponse(nil, err.Error(), ""))
		return
	}

	var res models.OrderResponse
	copier.Copy(&res, &orders)
	c.JSON(http.StatusOK, utils.PrepareResponse(res, "OK", ""))
}
