package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/quangdangfit/gocommon/logger"

	"goshop/internal/product/domain"
	"goshop/internal/product/service"
	"goshop/pkg/apperror"
	"goshop/pkg/config"
	"goshop/pkg/redis"
	"goshop/pkg/response"
	"goshop/pkg/utils"
)

type ProductHandler struct {
	cache   redis.Redis
	service service.ProductService
}

func NewProductHandler(
	cache redis.Redis,
	service service.ProductService,
) *ProductHandler {
	return &ProductHandler{
		cache:   cache,
		service: service,
	}
}

// GetProductByID godoc
//
//	@Summary	Get product by id
//	@Tags		products
//	@Produce	json
//	@Param		id	path	string	true	"Product ID"
//	@Router		/api/v1/products/{id} [get]
func (p *ProductHandler) GetProductByID(c *gin.Context) {
	var res domain.Product
	cacheKey := c.Request.URL.RequestURI()
	err := p.cache.Get(cacheKey, &res)
	if err == nil {
		response.JSON(c, http.StatusOK, res)
		return
	}

	productId := c.Param("id")
	product, err := p.service.GetProductByID(c, productId)
	if err != nil {
		logger.Error("Failed to get product detail: ", err)
		apperror.Wrap(apperror.ErrNotFound, err).HTTPError(c)
		return
	}

	utils.Copy(&res, &product)
	response.JSON(c, http.StatusOK, res)
	_ = p.cache.SetWithExpiration(cacheKey, res, config.ProductCachingTime)
}

// ListProducts godoc
//
//	@Summary	Get list products
//	@Tags		products
//	@Produce	json
//	@Success	200	{object}	dto.ListProductRes
//	@Router		/api/v1/products [get]
func (p *ProductHandler) ListProducts(c *gin.Context) {
	var req domain.ListProductReq
	if err := c.ShouldBindQuery(&req); err != nil {
		logger.Error("Failed to parse request query: ", err)
		apperror.Wrap(apperror.ErrBadRequest, err).HTTPError(c)
		return
	}

	var res domain.ListProductRes
	cacheKey := c.Request.URL.RequestURI()
	err := p.cache.Get(cacheKey, &res)
	if err == nil {
		response.JSON(c, http.StatusOK, res)
		return
	}

	products, pagination, err := p.service.ListProducts(c, &req)
	if err != nil {
		logger.Error("Failed to get list products: ", err)
		apperror.ToHTTPError(c, err, http.StatusInternalServerError, "Something went wrong")
		return
	}

	utils.Copy(&res.Products, &products)
	res.Pagination = pagination
	response.JSON(c, http.StatusOK, res)
	_ = p.cache.SetWithExpiration(cacheKey, res, config.ProductCachingTime)
}

// CreateProduct godoc
//
//	@Summary	create product
//	@Tags		products
//	@Produce	json
//	@Security	ApiKeyAuth
//	@Param		_	body	dto.CreateProductReq	true	"Body"
//	@Router		/api/v1/products [post]
func (p *ProductHandler) CreateProduct(c *gin.Context) {
	var req domain.CreateProductReq
	if err := c.ShouldBindJSON(&req); c.Request.Body == nil || err != nil {
		logger.Error("Failed to get body", err)
		apperror.Wrap(apperror.ErrBadRequest, err).HTTPError(c)
		return
	}

	product, err := p.service.Create(c, &req)
	if err != nil {
		logger.Error("Failed to create product", err.Error())
		apperror.ToHTTPError(c, err, http.StatusInternalServerError, "Something went wrong")
		return
	}

	var res domain.Product
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
//	@Param		id	path	string					true	"Product ID"
//	@Param		_	body	dto.UpdateProductReq	true	"Body"
//	@Router		/api/v1/products/{id} [put]
func (p *ProductHandler) UpdateProduct(c *gin.Context) {
	productId := c.Param("id")
	var req domain.UpdateProductReq
	if err := c.ShouldBindJSON(&req); c.Request.Body == nil || err != nil {
		logger.Error("Failed to get body", err)
		apperror.Wrap(apperror.ErrBadRequest, err).HTTPError(c)
		return
	}

	product, err := p.service.Update(c, productId, &req)
	if err != nil {
		logger.Error("Failed to update product", err.Error())
		apperror.ToHTTPError(c, err, http.StatusInternalServerError, "Something went wrong")
		return
	}

	var res domain.Product
	utils.Copy(&res, &product)
	response.JSON(c, http.StatusOK, res)
	_ = p.cache.RemovePattern("*product*")
}
