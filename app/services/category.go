package services

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"gitlab.com/quangdangfit/gocommon/utils/logger"

	"goshop/app/models"
	"goshop/app/repositories"
	"goshop/app/schema"
	"goshop/pkg/utils"
)

type ICategory interface {
	GetCategories(ctx context.Context, query models.CategoryQueryRequest) (*[]models.Category, error)
	GetCategoryByID(ctx context.Context, uuid string) (*models.Category, error)
	CreateCategory(cxt context.Context, item schema.Category) (*models.Category, error)
	UpdateCategory(c *gin.Context)
}

type category struct {
	repo repositories.ICategoryRepository
}

func NewCategoryService(repo repositories.ICategoryRepository) ICategory {
	return &category{repo: repo}
}

func (categ *category) GetCategories(ctx context.Context, query models.CategoryQueryRequest) (*[]models.Category, error) {
	categories, err := categ.repo.GetCategories(query)
	if err != nil {
		logger.Error("Failed to get categories: ", err)
		return nil, err
	}

	return categories, nil
}

func (categ *category) GetCategoryByID(ctx context.Context, uuid string) (*models.Category, error) {
	category, err := categ.repo.GetCategoryByID(uuid)
	if err != nil {
		logger.Error("Failed to get category: ", err)
		return nil, err
	}

	return category, nil
}

func (categ *category) CreateCategory(cxt context.Context, item schema.Category) (*models.Category, error) {
	category, err := categ.repo.CreateCategory(&item)
	if err != nil {
		logger.Error("Failed to create category", err.Error())
		return nil, err
	}

	return category, nil
}

func (categ *category) UpdateCategory(c *gin.Context) {
	uuid := c.Param("uuid")
	var reqBody models.CategoryBodyRequest
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		logger.Error("Failed to parse request body: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	categories, err := categ.repo.UpdateCategory(uuid, &reqBody)
	if err != nil {
		logger.Error("Failed to update category: ", err)
		c.JSON(http.StatusBadRequest, utils.PrepareResponse(nil, err.Error(), ""))
		return
	}

	var res models.CategoryResponse
	copier.Copy(&res, &categories)
	c.JSON(http.StatusOK, utils.PrepareResponse(res, "OK", ""))
}
