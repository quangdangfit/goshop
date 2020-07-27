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
	Service services.Category
}

func (categ *Category) init(service services.Category) {
	categ.Service = service
}

func (categ *Category) GetCategories(c *gin.Context) {
	var reqQuery models.CategoryQueryRequest
	if err := c.ShouldBindQuery(&reqQuery); err != nil {
		logger.Error("Failed to parse request query: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	rs, err := categ.Service.GetCategories(ctx, reqQuery)
	if err != nil {
		logger.Error("Failed to get categories: ", err)
		c.JSON(http.StatusBadRequest, utils.PrepareResponse(nil, err.Error(), ""))
		return
	}

	var res []models.CategoryResponse
	copier.Copy(&res, &rs)
	c.JSON(http.StatusOK, utils.PrepareResponse(res, "OK", ""))
}
