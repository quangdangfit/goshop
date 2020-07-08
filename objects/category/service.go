package category

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jinzhu/copier"
	"gitlab.com/quangdangfit/gocommon/utils/logger"
	"goshop/utils"
	"net/http"
)

type Service interface {
	GetCategories(c *gin.Context)
	GetCategoryByID(c *gin.Context)
	CreateCategory(c *gin.Context)
	UpdateCategory(c *gin.Context)
}

type service struct {
	repo Repository
}

func NewService() Service {
	return &service{repo: NewRepository()}
}

func (s *service) GetCategories(c *gin.Context) {
	activeParam := c.Query("active")
	active := true
	if activeParam == "false" {
		active = false
	}

	categories, err := s.repo.GetCategories(active)
	if err != nil {
		logger.Error("Failed to get categories: ", err)
		c.JSON(http.StatusBadRequest, utils.PrepareResponse(nil, err.Error(), ""))
		return
	}

	var res []CategoryResponse
	copier.Copy(&res, &categories)
	c.JSON(http.StatusOK, utils.PrepareResponse(res, "OK", ""))
}

func (s *service) GetCategoryByID(c *gin.Context) {
	categoryId := c.Param("uuid")

	category, err := s.repo.GetCategoryByID(categoryId)
	if err != nil {
		logger.Error("Failed to get category: ", err)
		c.JSON(http.StatusBadRequest, utils.PrepareResponse(nil, err.Error(), ""))
		return
	}

	var res CategoryResponse
	copier.Copy(&res, &category)
	c.JSON(http.StatusOK, utils.PrepareResponse(res, "OK", ""))
}

func (s *service) CreateCategory(c *gin.Context) {
	var reqBody CategoryRequest
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

	categories, err := s.repo.CreateCategory(&reqBody)
	if err != nil {
		logger.Error("Failed to create category", err.Error())
		c.JSON(http.StatusBadRequest, utils.PrepareResponse(nil, err.Error(), ""))
		return
	}

	var res CategoryResponse
	copier.Copy(&res, &categories)
	c.JSON(http.StatusOK, utils.PrepareResponse(res, "OK", ""))
}

func (s *service) UpdateCategory(c *gin.Context) {
	uuid := c.Param("uuid")
	var reqBody CategoryRequest
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		logger.Error("Failed to parse request body: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	categories, err := s.repo.UpdateCategory(uuid, &reqBody)
	if err != nil {
		logger.Error("Failed to update category: ", err)
		c.JSON(http.StatusBadRequest, utils.PrepareResponse(nil, err.Error(), ""))
		return
	}

	var res CategoryResponse
	copier.Copy(&res, &categories)
	c.JSON(http.StatusOK, utils.PrepareResponse(res, "OK", ""))
}
