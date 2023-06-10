package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/quangdangfit/gocommon/logger"

	"goshop/app/serializers"
	"goshop/app/services"
	"goshop/pkg/response"
	"goshop/pkg/utils"
	"goshop/pkg/validation"
)

type ProductAPI struct {
	validator validation.Validation
	service   services.IProductService
}

func NewProductAPI(service services.IProductService) *ProductAPI {
	return &ProductAPI{
		validator: validation.New(),
		service:   service,
	}
}

// GetProductByID godoc
//
//	@Summary	Get product by id
//	@Produce	json
//	@Param		id	path	string	true	"Product ID"
//	@Security	ApiKeyAuth
//	@Success	200	{object}	serializers.Product
//	@Router		/api/v1/products/{id} [get]
func (p *ProductAPI) GetProductByID(c *gin.Context) {
	productId := c.Param("id")
	product, err := p.service.GetProductByID(c, productId)
	if err != nil {
		logger.Error("Failed to get product detail: ", err)
		response.Error(c, http.StatusNotFound, err, "Not found")
		return
	}

	var res serializers.Product
	utils.Copy(&res, &product)
	response.JSON(c, http.StatusOK, res)
}

// ListProducts godoc
//
//	@Summary	Get list products
//	@Produce	json
//	@Security	ApiKeyAuth
//	@Success	200	{object}	serializers.ListProductRes
//	@Router		/api/v1/products [get]
func (p *ProductAPI) ListProducts(c *gin.Context) {
	var req serializers.ListProductReq
	if err := c.ShouldBindQuery(&req); err != nil {
		logger.Error("Failed to parse request query: ", err)
		response.Error(c, http.StatusBadRequest, err, "Invalid parameters")
		return
	}

	products, pagination, err := p.service.ListProducts(c, req)
	if err != nil {
		logger.Error("Failed to get list products: ", err)
		response.Error(c, http.StatusInternalServerError, err, "Something went wrong")
		return
	}

	var res serializers.ListProductRes
	utils.Copy(&res.Products, &products)
	res.Pagination = pagination
	response.JSON(c, http.StatusOK, res)
}

// CreateProduct godoc
//
//	@Summary	create product
//	@Produce	json
//	@Security	ApiKeyAuth
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

	ctx := c.Request.Context()
	product, err := p.service.Create(ctx, &req)
	if err != nil {
		logger.Error("Failed to create product", err.Error())
		response.Error(c, http.StatusInternalServerError, err, "Something went wrong")
		return
	}

	var res serializers.Product
	utils.Copy(&res, &product)
	response.JSON(c, http.StatusOK, res)
}

// UpdateProduct godoc
//
//	@Summary	update product
//	@Produce	json
//	@Security	ApiKeyAuth
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

	ctx := c.Request.Context()
	product, err := p.service.Update(ctx, productId, &req)
	if err != nil {
		logger.Error("Failed to update product", err.Error())
		response.Error(c, http.StatusInternalServerError, err, "Something went wrong")
		return
	}

	var res serializers.Product
	utils.Copy(&res, &product)
	response.JSON(c, http.StatusOK, res)
}
