package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"gitlab.com/quangdangfit/gocommon/utils/logger"

	"goshop/app/models"
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

func (categ *Product) GetProducts(c *gin.Context) {
	var params schema.ProductQueryParams
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

	var res []models.ProductResponse
	copier.Copy(&res, &rs)
	c.JSON(http.StatusOK, utils.PrepareResponse(res, "OK", ""))
}
