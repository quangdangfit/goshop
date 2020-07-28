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

type Category struct {
	service services.ICategory
}

func NewCategoryAPI(service services.ICategory) *Category {
	return &Category{service: service}
}

func (categ *Category) GetCategories(c *gin.Context) {
	var query schema.CategoryQueryParam
	if err := c.ShouldBindQuery(&query); err != nil {
		logger.Error("Failed to parse request query: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	rs, err := categ.service.GetCategories(ctx, &query)
	if err != nil {
		logger.Error("Failed to get categories: ", err)
		c.JSON(http.StatusBadRequest, utils.PrepareResponse(nil, err.Error(), ""))
		return
	}

	var res []schema.Category
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

	var res schema.Category
	copier.Copy(&res, &category)
	c.JSON(http.StatusOK, utils.PrepareResponse(res, "OK", ""))
}

func (categ *Category) CreateCategory(c *gin.Context) {
	var item schema.Category
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
	categories, err := categ.service.CreateCategory(ctx, &item)
	if err != nil {
		logger.Error("Failed to create category", err.Error())
		c.JSON(http.StatusBadRequest, utils.PrepareResponse(nil, err.Error(), ""))
		return
	}

	var res schema.Category
	copier.Copy(&res, &categories)
	c.JSON(http.StatusOK, utils.PrepareResponse(res, "OK", ""))
}

func (categ *Category) UpdateCategory(c *gin.Context) {
	uuid := c.Param("uuid")
	var item schema.CategoryBodyParam
	if err := c.ShouldBind(&item); err != nil {
		logger.Error("Failed to parse request body: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	categories, err := categ.service.UpdateCategory(ctx, uuid, &item)
	if err != nil {
		logger.Error("Failed to update category: ", err)
		c.JSON(http.StatusBadRequest, utils.PrepareResponse(nil, err.Error(), ""))
		return
	}

	var res schema.Category
	copier.Copy(&res, &categories)
	c.JSON(http.StatusOK, utils.PrepareResponse(res, "OK", ""))
}
