package services

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jinzhu/copier"
	"gitlab.com/quangdangfit/gocommon/utils/logger"

	"goshop/models"
	"goshop/repositories"
	"goshop/utils"
)

type ProductService interface {
	GetProducts(c *gin.Context)
	GetProductByID(c *gin.Context)
	CreateProduct(c *gin.Context)
	UpdateProduct(c *gin.Context)
	GetProductByCategory(c *gin.Context)
}

type product struct {
	repo repositories.ProductRepository
}

func NewProductService(repo repositories.ProductRepository) ProductService {
	return &product{repo: repo}
}

// GetProductByID godoc
// @Summary Get get product by uuid
// @Produce json
// @Param uuid path string true "Product UUID"
// @Security ApiKeyAuth
// @Success 200 {object} product.ProductResponse
// @Router /api/v1/products/{uuid} [get]
func (p *product) GetProductByID(c *gin.Context) {
	productId := c.Param("uuid")

	product, err := p.repo.GetProductByID(productId)
	if err != nil {
		logger.Error("Failed to get product: ", err)
		c.JSON(http.StatusBadRequest, utils.PrepareResponse(nil, err.Error(), ""))
		return
	}

	var res models.ProductResponse
	copier.Copy(&res, &product)
	c.JSON(http.StatusOK, utils.PrepareResponse(res, "OK", ""))
}

// GetProducts godoc
// @Summary Get list products
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} []product.ProductResponse
// @Router /api/v1/products [get]
func (p *product) GetProducts(c *gin.Context) {
	activeParam := c.Query("active")
	active := true
	if activeParam == "false" {
		active = false
	}

	products, err := p.repo.GetProducts(active)
	if err != nil {
		logger.Error("Failed to get products: ", err)
		c.JSON(http.StatusBadRequest, utils.PrepareResponse(nil, err.Error(), ""))
		return
	}

	var res []models.ProductResponse
	copier.Copy(&res, &products)
	c.JSON(http.StatusOK, utils.PrepareResponse(res, "OK", ""))
}

func (p *product) GetProductByCategory(c *gin.Context) {
	categUUID := c.Param("uuid")
	activeParam := c.Query("active")
	active := true
	if activeParam == "false" {
		active = false
	}

	products, err := p.repo.GetProductByCategory(categUUID, active)
	if err != nil {
		logger.Error("Failed to get products: ", err)
		c.JSON(http.StatusBadRequest, utils.PrepareResponse(nil, err.Error(), ""))
		return
	}

	var res []models.ProductResponse
	copier.Copy(&res, &products)
	c.JSON(http.StatusOK, utils.PrepareResponse(res, "OK", ""))
}

func (p *product) CreateProduct(c *gin.Context) {
	var reqBody models.ProductRequest
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

	products, err := p.repo.CreateProduct(&reqBody)
	if err != nil {
		logger.Error("Failed to create product", err.Error())
		c.JSON(http.StatusBadRequest, utils.PrepareResponse(nil, err.Error(), ""))
		return
	}

	var res []models.ProductResponse
	copier.Copy(&res, &products)
	c.JSON(http.StatusOK, utils.PrepareResponse(res, "OK", ""))
}

func (p *product) UpdateProduct(c *gin.Context) {
	uuid := c.Param("uuid")
	var reqBody models.ProductRequest
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		logger.Error("Failed to parse request body: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	products, err := p.repo.UpdateProduct(uuid, &reqBody)
	if err != nil {
		logger.Error("Failed to update product: ", err)
		c.JSON(http.StatusBadRequest, utils.PrepareResponse(nil, err.Error(), ""))
		return
	}

	var res models.ProductResponse
	copier.Copy(&res, &products)
	c.JSON(http.StatusOK, utils.PrepareResponse(res, "OK", ""))
}
