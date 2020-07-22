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

type QuantitySerivce interface {
	GetQuantities(c *gin.Context)
	GetQuantityByID(c *gin.Context)
	CreateQuantity(c *gin.Context)
	UpdateQuantity(c *gin.Context)
}

type stockQuant struct {
	repo repositories.QuantityRepository
}

func NewQuantityService(repo repositories.QuantityRepository) QuantitySerivce {
	return &stockQuant{repo: repo}
}

func (categ *stockQuant) GetQuantities(c *gin.Context) {
	var reqQuery models.QuantityQueryRequest
	if err := c.ShouldBindQuery(&reqQuery); err != nil {
		logger.Error("Failed to parse request query: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var query map[string]interface{}
	data, _ := json.Marshal(reqQuery)
	json.Unmarshal(data, &query)
	quantities, err := categ.repo.GetQuantities(query)
	if err != nil {
		logger.Error("Failed to get quantities}: ", err)
		c.JSON(http.StatusBadRequest, utils.PrepareResponse(nil, err.Error(), ""))
		return
	}

	var res []models.QuantityResponse
	copier.Copy(&res, &quantities)
	c.JSON(http.StatusOK, utils.PrepareResponse(res, "OK", ""))
}

func (categ *stockQuant) GetQuantityByID(c *gin.Context) {
	stockQuantId := c.Param("uuid")

	stockQuant, err := categ.repo.GetQuantityByID(stockQuantId)
	if err != nil {
		logger.Error("Failed to get stockQuant: ", err)
		c.JSON(http.StatusBadRequest, utils.PrepareResponse(nil, err.Error(), ""))
		return
	}

	var res models.QuantityResponse
	copier.Copy(&res, &stockQuant)
	c.JSON(http.StatusOK, utils.PrepareResponse(res, "OK", ""))
}

func (categ *stockQuant) CreateQuantity(c *gin.Context) {
	var reqBody models.QuantityBodyRequest
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

	quantities, err := categ.repo.CreateQuantity(&reqBody)
	if err != nil {
		logger.Error("Failed to create stockQuant", err.Error())
		c.JSON(http.StatusBadRequest, utils.PrepareResponse(nil, err.Error(), ""))
		return
	}

	var res models.QuantityResponse
	copier.Copy(&res, &quantities)
	c.JSON(http.StatusOK, utils.PrepareResponse(res, "OK", ""))
}

func (categ *stockQuant) UpdateQuantity(c *gin.Context) {
	uuid := c.Param("uuid")
	var reqBody models.QuantityBodyRequest
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		logger.Error("Failed to parse request body: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	quantities, err := categ.repo.UpdateQuantity(uuid, &reqBody)
	if err != nil {
		logger.Error("Failed to update stockQuant: ", err)
		c.JSON(http.StatusBadRequest, utils.PrepareResponse(nil, err.Error(), ""))
		return
	}

	var res models.QuantityResponse
	copier.Copy(&res, &quantities)
	c.JSON(http.StatusOK, utils.PrepareResponse(res, "OK", ""))
}
