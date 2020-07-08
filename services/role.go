package services

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jinzhu/copier"
	"gitlab.com/quangdangfit/gocommon/utils/logger"
	"goshop/models"
	"goshop/repositories"
	"goshop/utils"
	"net/http"
)

type RoleService interface {
	CreateRole(c *gin.Context)
}

type roleService struct {
	repo repositories.RoleRepository
}

func NewService(repo repositories.RoleRepository) RoleService {
	return &roleService{repo: repo}
}

func (r *roleService) CreateRole(c *gin.Context) {
	var reqBody models.RoleRequest
	if err := c.ShouldBindJSON(&reqBody); err != nil {
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

	user, err := r.repo.CreateRole(&reqBody)
	if err != nil {
		logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, utils.PrepareResponse(nil, err.Error(), ""))
		return
	}

	var res models.RoleResponse
	copier.Copy(&res, &user)
	c.JSON(http.StatusOK, utils.PrepareResponse(res, "OK", ""))
}
