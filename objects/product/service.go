package product

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jinzhu/copier"
	"gitlab.com/quangdangfit/gocommon/utils/logger"
	"goshop/utils"
	"net/http"
)

type Service interface {
	GetProducts(c *gin.Context)
	GetProductByID(c *gin.Context)
	CreateProduct(c *gin.Context)
	UpdateProduct(c *gin.Context)
	GetProductByCategory(c *gin.Context)
}

type service struct {
	repo Repository
}

func NewService() Service {
	return &service{repo: NewRepository()}
}

func (s *service) GetProductByID(c *gin.Context) {
	productId := c.Param("uuid")

	product, err := s.repo.GetProductByID(productId)
	if err != nil {
		logger.Error("Failed to get product: ", err)
		c.JSON(http.StatusBadRequest, utils.PrepareResponse(nil, err.Error(), ""))
		return
	}

	var res ProductResponse
	copier.Copy(&res, &product)
	c.JSON(http.StatusOK, utils.PrepareResponse(res, "OK", ""))
}

func (s *service) GetProducts(c *gin.Context) {
	activeParam := c.Query("active")
	active := true
	if activeParam == "false" {
		active = false
	}

	products, err := s.repo.GetProducts(active)
	if err != nil {
		logger.Error("Failed to get products: ", err)
		c.JSON(http.StatusBadRequest, utils.PrepareResponse(nil, err.Error(), ""))
		return
	}

	var res []ProductResponse
	copier.Copy(&res, &products)
	c.JSON(http.StatusOK, utils.PrepareResponse(res, "OK", ""))
}

func (s *service) GetProductByCategory(c *gin.Context) {
	categUUID := c.Param("uuid")
	activeParam := c.Query("active")
	active := true
	if activeParam == "false" {
		active = false
	}

	products, err := s.repo.GetProductByCategory(categUUID, active)
	if err != nil {
		logger.Error("Failed to get products: ", err)
		c.JSON(http.StatusBadRequest, utils.PrepareResponse(nil, err.Error(), ""))
		return
	}

	var res []ProductResponse
	copier.Copy(&res, &products)
	c.JSON(http.StatusOK, utils.PrepareResponse(res, "OK", ""))
}

func (s *service) CreateProduct(c *gin.Context) {
	var reqBody ProductRequest
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

	products, err := s.repo.CreateProduct(&reqBody)
	if err != nil {
		logger.Error("Failed to create product", err.Error())
		c.JSON(http.StatusBadRequest, utils.PrepareResponse(nil, err.Error(), ""))
		return
	}

	var res []ProductResponse
	copier.Copy(&res, &products)
	c.JSON(http.StatusOK, utils.PrepareResponse(res, "OK", ""))
}

func (s *service) UpdateProduct(c *gin.Context) {
	uuid := c.Param("uuid")
	var reqBody ProductRequest
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		logger.Error("Failed to parse request body: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	products, err := s.repo.UpdateProduct(uuid, &reqBody)
	if err != nil {
		logger.Error("Failed to update product: ", err)
		c.JSON(http.StatusBadRequest, utils.PrepareResponse(nil, err.Error(), ""))
		return
	}

	var res ProductResponse
	copier.Copy(&res, &products)
	c.JSON(http.StatusOK, utils.PrepareResponse(res, "OK", ""))
}
