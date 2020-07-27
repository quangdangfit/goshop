package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"gitlab.com/quangdangfit/gocommon/utils/logger"

	"goshop/app/models"
	"goshop/app/services"
	"goshop/pkg/utils"
)

type Category struct {
	service services.ICategory
}

func NewCategoryAPI(service services.ICategory) *Category {
	return &Category{service: service}
}

func (categ *Category) GetCategories(c *gin.Context) {
	var reqQuery models.CategoryQueryRequest
	if err := c.ShouldBindQuery(&reqQuery); err != nil {
		logger.Error("Failed to parse request query: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	rs, err := categ.service.GetCategories(ctx, reqQuery)
	if err != nil {
		logger.Error("Failed to get categories: ", err)
		c.JSON(http.StatusBadRequest, utils.PrepareResponse(nil, err.Error(), ""))
		return
	}

	var res []models.CategoryResponse
	copier.Copy(&res, &rs)
	c.JSON(http.StatusOK, utils.PrepareResponse(res, "OK", ""))
}

func (categ *Category) GetCategoryByID(c *gin.Context) {
	categoryId := c.Param("uuid")

	ctx := c.Request.Context()
	category, err := categ.service.GetCategoryByID(ctx, categoryId)
	if err != nil {
		logger.Error("Failed to get category: ", err)
		c.JSON(http.StatusBadRequest, utils.PrepareResponse(nil, err.Error(), ""))
		return
	}

	var res models.CategoryResponse
	copier.Copy(&res, &category)
	c.JSON(http.StatusOK, utils.PrepareResponse(res, "OK", ""))
}
