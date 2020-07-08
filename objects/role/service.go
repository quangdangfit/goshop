package role

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jinzhu/copier"
	"gitlab.com/quangdangfit/gocommon/utils/logger"
	"goshop/utils"
	"net/http"
)

type Service interface {
	CreateRole(c *gin.Context)
}

type service struct {
	repo Repository
}

func NewService() Service {
	return &service{repo: NewRepository()}
}

func (s *service) CreateRole(c *gin.Context) {
	var reqBody RoleRequest
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

	user, err := s.repo.CreateRole(&reqBody)
	if err != nil {
		logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, utils.PrepareResponse(nil, err.Error(), ""))
		return
	}

	var res RoleResponse
	copier.Copy(&res, &user)
	c.JSON(http.StatusOK, utils.PrepareResponse(res, "OK", ""))
}
