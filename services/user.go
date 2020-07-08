package services

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"gitlab.com/quangdangfit/gocommon/utils/logger"
	"goshop/models"
	"goshop/repositories"
	"goshop/utils"
	"net/http"
)

type User interface {
	Login(c *gin.Context)
	Register(c *gin.Context)
	GetUserByID(c *gin.Context)
}

type service struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) User {
	return &service{repo: repo}
}

func (s *service) validate(r models.RegisterRequest) bool {
	return utils.Validate(
		[]utils.Validation{
			{Value: r.Username, Valid: "username"},
			{Value: r.Email, Valid: "email"},
			{Value: r.Password, Valid: "password"},
		})
}

func (s *service) checkPermission(uuid string, data map[string]interface{}) bool {
	return data["uuid"] == uuid
}

func (s *service) Login(c *gin.Context) {
	var reqBody models.LoginRequest
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := s.repo.Login(&reqBody)
	if err != nil {
		logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, utils.PrepareResponse(nil, err.Error(), ""))
		return
	}

	var res models.UserResponse
	copier.Copy(&res, &user)
	res.Extra = map[string]interface{}{
		"token": utils.GenerateToken(user),
	}
	c.JSON(http.StatusOK, utils.PrepareResponse(res, "OK", ""))
}

func (s *service) Register(c *gin.Context) {
	var reqBody models.RegisterRequest
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	valid := s.validate(reqBody)
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request body is invalid"})
		return
	}

	user, err := s.repo.Register(&reqBody)
	if err != nil {
		logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, utils.PrepareResponse(nil, err.Error(), ""))
		return
	}

	var res models.UserResponse
	copier.Copy(&res, &user)
	res.Extra = map[string]interface{}{
		"token": utils.GenerateToken(user),
	}
	c.JSON(http.StatusOK, utils.PrepareResponse(res, "OK", ""))
}

func (s *service) GetUserByID(c *gin.Context) {
	userUUID := c.Param("uuid")
	user, err := s.repo.GetUserByID(userUUID)
	if err != nil {
		logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, utils.PrepareResponse(nil, err.Error(), utils.ErrorNotExistUser))
		return
	}

	var res models.UserResponse
	copier.Copy(&res, &user)
	c.JSON(http.StatusOK, utils.PrepareResponse(res, "OK", ""))
}
