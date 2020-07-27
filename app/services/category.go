package services

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jinzhu/copier"
	"gitlab.com/quangdangfit/gocommon/utils/logger"

	"goshop/app/models"
	"goshop/app/repositories"
	"goshop/pkg/utils"
)

type ICategory interface {
	GetCategories(ctx context.Context, query models.CategoryQueryRequest) (*[]models.Category, error)
	GetCategoryByID(ctx context.Context, uuid string) (*models.Category, error)
	CreateCategory(c *gin.Context)
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

func (categ *category) CreateCategory(c *gin.Context) {
	var reqBody models.CategoryBodyRequest
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

	categories, err := categ.repo.CreateCategory(&reqBody)
	if err != nil {
		logger.Error("Failed to create category", err.Error())
		c.JSON(http.StatusBadRequest, utils.PrepareResponse(nil, err.Error(), ""))
		return
	}

	var res models.CategoryResponse
	copier.Copy(&res, &categories)
	c.JSON(http.StatusOK, utils.PrepareResponse(res, "OK", ""))
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
