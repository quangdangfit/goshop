package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/quangdangfit/gocommon/logger"
	"github.com/quangdangfit/gocommon/redis"
	"github.com/quangdangfit/gocommon/validation"

	"goshop/app/serializers"
	"goshop/app/services"
	"goshop/config"
	"goshop/pkg/response"
	"goshop/pkg/utils"
)

type ProductAPI struct {
	validator validation.Validation
	cache     redis.IRedis
	service   services.IProductService
}

func NewProductAPI(
	validator validation.Validation,
	cache redis.IRedis,
	service services.IProductService,
) *ProductAPI {
	return &ProductAPI{
		validator: validator,
		cache:     cache,
		service:   service,
	}
}

// GetProductByID godoc
//
//	@Summary	Get product by id
//	@Tags		products
//	@Produce	json
//	@Param		id	path		string	true	"Product ID"
//	@Success	200	{object}	serializers.Product
//	@Router		/api/v1/products/{id} [get]
func (p *ProductAPI) GetProductByID(c *gin.Context) {
	var res serializers.Product
	err := p.cache.Get(c.Request.URL.Path, &res)
	if err == nil {
		response.JSON(c, http.StatusOK, res)
		return
	}

	productId := c.Param("id")
	product, err := p.service.GetProductByID(c, productId)
	if err != nil {
		logger.Error("Failed to get product detail: ", err)
		response.Error(c, http.StatusNotFound, err, "Not found")
		return
	}

	utils.Copy(&res, &product)
	response.JSON(c, http.StatusOK, res)
	_ = p.cache.SetWithExpiration(c.Request.URL.Path, res, config.ProductCachingTime)
}

// ListProducts godoc
//
//	@Summary	Get list products
//	@Tags		products
//	@Produce	json
//	@Success	200	{object}	serializers.ListProductRes
//	@Router		/api/v1/products [get]
func (p *ProductAPI) ListProducts(c *gin.Context) {
	var req serializers.ListProductReq
	if err := c.ShouldBindQuery(&req); err != nil {
		logger.Error("Failed to parse request query: ", err)
		response.Error(c, http.StatusBadRequest, err, "Invalid parameters")
		return
	}

	var res serializers.ListProductRes
	err := p.cache.Get(c.Request.URL.Path, &res)
	if err == nil {
		response.JSON(c, http.StatusOK, res)
		return
	}

	products, pagination, err := p.service.ListProducts(c, req)
	if err != nil {
		logger.Error("Failed to get list products: ", err)
		response.Error(c, http.StatusInternalServerError, err, "Something went wrong")
		return
	}

	utils.Copy(&res.Products, &products)
	res.Pagination = pagination
	response.JSON(c, http.StatusOK, res)
	_ = p.cache.SetWithExpiration(c.Request.URL.Path, res, config.ProductCachingTime)
}

// CreateProduct godoc
//
//	@Summary	create product
//	@Tags		products
//	@Produce	json
//	@Security	ApiKeyAuth
//	@Param		_	body		serializers.CreateProductReq	true	"Body"
//	@Success	200	{object}	serializers.Product
//	@Router		/api/v1/products [post]
func (p *ProductAPI) CreateProduct(c *gin.Context) {
	var req serializers.CreateProductReq
	if err := c.ShouldBindJSON(&req); c.Request.Body == nil || err != nil {
		logger.Error("Failed to get body", err)
		response.Error(c, http.StatusBadRequest, err, "Invalid parameters")
		return
	}

	if err := p.validator.ValidateStruct(req); err != nil {
		response.Error(c, http.StatusBadRequest, err, "Invalid parameters")
		return
	}

	product, err := p.service.Create(c, &req)
	if err != nil {
		logger.Error("Failed to create product", err.Error())
		response.Error(c, http.StatusInternalServerError, err, "Something went wrong")
		return
	}

	var res serializers.Product
	utils.Copy(&res, &product)
	response.JSON(c, http.StatusOK, res)
	_ = p.cache.RemovePattern("*product*")
}

// UpdateProduct godoc
//
//	@Summary	update product
//	@Tags		products
//	@Produce	json
//	@Security	ApiKeyAuth
//	@Param		id	path		string							true	"Product ID"
//	@Param		_	body		serializers.UpdateProductReq	true	"Body"
//	@Success	200	{object}	serializers.Product
//	@Router		/api/v1/products/{id} [put]
func (p *ProductAPI) UpdateProduct(c *gin.Context) {
	productId := c.Param("id")
	var req serializers.UpdateProductReq
	if err := c.ShouldBindJSON(&req); c.Request.Body == nil || err != nil {
		logger.Error("Failed to get body", err)
		response.Error(c, http.StatusBadRequest, err, "Invalid parameters")
		return
	}

	if err := p.validator.ValidateStruct(req); err != nil {
		response.Error(c, http.StatusBadRequest, err, "Invalid parameters")
		return
	}

	product, err := p.service.Update(c, productId, &req)
	if err != nil {
		logger.Error("Failed to update product", err.Error())
		response.Error(c, http.StatusInternalServerError, err, "Something went wrong")
		return
	}

	var res serializers.Product
	utils.Copy(&res, &product)
	response.JSON(c, http.StatusOK, res)
	_ = p.cache.RemovePattern("*product*")
}
