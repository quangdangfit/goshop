package services

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"gitlab.com/quangdangfit/gocommon/utils/logger"

	jwtMiddle "goshop/middleware/jwt"
	"goshop/models"
	"goshop/repositories"
	"goshop/utils"
)

type UserService interface {
	Login(c *gin.Context)
	Register(c *gin.Context)
	GetUserByID(c *gin.Context)
}

type user struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &user{repo: repo}
}

func (u *user) validate(r models.RegisterRequest) bool {
	return utils.Validate(
		[]utils.Validation{
			{Value: r.Username, Valid: "username"},
			{Value: r.Email, Valid: "email"},
			{Value: r.Password, Valid: "password"},
		})
}

func (u *user) checkPermission(uuid string, data map[string]interface{}) bool {
	return data["uuid"] == uuid
}

func (u *user) Login(c *gin.Context) {
	var reqBody models.LoginRequest
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := u.repo.Login(&reqBody)
	if err != nil {
		logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, utils.PrepareResponse(nil, err.Error(), ""))
		return
	}

	var res models.UserResponse
	copier.Copy(&res, &user)
	res.Extra = map[string]interface{}{
		"token": jwtMiddle.GenerateToken(user),
	}
	c.JSON(http.StatusOK, utils.PrepareResponse(res, "OK", ""))
}

func (u *user) Register(c *gin.Context) {
	var reqBody models.RegisterRequest
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	valid := u.validate(reqBody)
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request body is invalid"})
		return
	}

	user, err := u.repo.Register(&reqBody)
	if err != nil {
		logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, utils.PrepareResponse(nil, err.Error(), ""))
		return
	}

	var res models.UserResponse
	copier.Copy(&res, &user)
	res.Extra = map[string]interface{}{
		"token": jwtMiddle.GenerateToken(user),
	}
	c.JSON(http.StatusOK, utils.PrepareResponse(res, "OK", ""))
}

func (u *user) GetUserByID(c *gin.Context) {
	userUUID := c.Param("uuid")
	user, err := u.repo.GetUserByID(userUUID)
	if err != nil {
		logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, utils.PrepareResponse(nil, err.Error(), utils.ErrorNotExistUser))
		return
	}

	var res models.UserResponse
	copier.Copy(&res, &user)
	c.JSON(http.StatusOK, utils.PrepareResponse(res, "OK", ""))
}
