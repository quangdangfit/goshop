package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jinzhu/copier"
	"gitlab.com/quangdangfit/gocommon/utils/logger"

	"goshop/app/schema"
	"goshop/app/services"
	"goshop/pkg/utils"
)

type Product struct {
	service services.IProductService
}

func NewProductAPI(service services.IProductService) *Product {
	return &Product{service: service}
}

// GetProductByID godoc
// @Summary Get get product by uuid
// @Produce json
// @Param uuid path string true "Product UUID"
// @Security ApiKeyAuth
// @Success 200 {object} product.ProductResponse
// @Router /api/v1/products/{uuid} [get]
func (p *Product) GetProductByID(c *gin.Context) {
	productId := c.Param("uuid")

	ctx := c.Request.Context()
	product, err := p.service.GetProductByID(ctx, productId)
	if err != nil {
		logger.Error("Failed to get product: ", err)
		c.JSON(http.StatusBadRequest, utils.PrepareResponse(nil, err.Error(), ""))
		return
	}

	var res schema.Product
	copier.Copy(&res, &product)
	c.JSON(http.StatusOK, utils.PrepareResponse(res, "OK", ""))
}

// GetProducts godoc
// @Summary Get list products
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} []product.ProductResponse
// @Router /api/v1/products [get]
func (categ *Product) GetProducts(c *gin.Context) {
	var params schema.ProductQueryParam
	if err := c.ShouldBindQuery(&params); err != nil {
		logger.Error("Failed to parse request query: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	rs, err := categ.service.GetProducts(ctx, params)
	if err != nil {
		logger.Error("Failed to get categories: ", err)
		c.JSON(http.StatusBadRequest, utils.PrepareResponse(nil, err.Error(), ""))
		return
	}

	var res []schema.Product
	copier.Copy(&res, &rs)
	c.JSON(http.StatusOK, utils.PrepareResponse(res, "OK", ""))
}

func (p *Product) GetProductByCategoryID(c *gin.Context) {
	categUUID := c.Param("uuid")

	ctx := c.Request.Context()
	products, err := p.service.GetProductByCategoryID(ctx, categUUID)
	if err != nil {
		logger.Error("Failed to get products: ", err)
		c.JSON(http.StatusBadRequest, utils.PrepareResponse(nil, err.Error(), ""))
		return
	}

	var res []schema.Product
	copier.Copy(&res, &products)
	c.JSON(http.StatusOK, utils.PrepareResponse(res, "OK", ""))
}

func (p *Product) CreateProduct(c *gin.Context) {
	var item schema.ProductBodyParam
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
	products, err := p.service.CreateProduct(ctx, &item)
	if err != nil {
		logger.Error("Failed to create product", err.Error())
		c.JSON(http.StatusBadRequest, utils.PrepareResponse(nil, err.Error(), ""))
		return
	}

	var res []schema.Product
	copier.Copy(&res, &products)
	c.JSON(http.StatusOK, utils.PrepareResponse(res, "OK", ""))
}

func (p *Product) UpdateProduct(c *gin.Context) {
	uuid := c.Param("uuid")
	var item schema.ProductBodyParam
	if err := c.ShouldBindJSON(&item); err != nil {
		logger.Error("Failed to parse request body: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	products, err := p.service.UpdateProduct(ctx, uuid, &item)
	if err != nil {
		logger.Error("Failed to update product: ", err)
		c.JSON(http.StatusBadRequest, utils.PrepareResponse(nil, err.Error(), ""))
		return
	}

	var res schema.Product
	copier.Copy(&res, &products)
	c.JSON(http.StatusOK, utils.PrepareResponse(res, "OK", ""))
}
