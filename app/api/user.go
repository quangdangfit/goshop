package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/quangdangfit/gocommon/logger"

	"goshop/app/serializers"
	"goshop/app/services"
	"goshop/pkg/response"
	"goshop/pkg/validation"
)

type User struct {
	validator validation.Validation
	service   services.IUserService
}

func NewUserAPI(service services.IUserService) *User {
	return &User{
		validator: validation.New(),
		service:   service,
	}
}

func (u *User) Login(c *gin.Context) {
	var params serializers.LoginReq
	if err := c.ShouldBindJSON(&params); c.Request.Body == nil || err != nil {
		logger.Error("Failed to get body", err)
		response.Error(c, http.StatusBadRequest, err, "Invalid parameters")
		return
	}

	user, accessToken, refreshToken, err := u.service.Login(c, &params)
	if err != nil {
		logger.Error("Failed to get body", err)
		response.Error(c, http.StatusBadRequest, err, "Something was wrong")
		return
	}

	var res serializers.LoginRes
	err = copier.Copy(&res.User, &user)
	if err != nil {
		logger.Error(err.Error())
		response.Error(c, http.StatusInternalServerError, err, "Something was wrong")
		return
	}
	res.AccessToken = accessToken
	res.RefreshToken = refreshToken
	response.JSON(c, http.StatusOK, res)
}

func (u *User) Register(c *gin.Context) {
	var params serializers.RegisterReq
	if err := c.ShouldBindJSON(&params); c.Request.Body == nil || err != nil {
		logger.Error("Failed to get body", err)
		response.Error(c, http.StatusBadRequest, err, "Invalid parameters")
		return
	}

	if err := u.validator.ValidateStruct(params); err != nil {
		response.Error(c, http.StatusBadRequest, err, "Invalid parameters")
		return
	}

	user, err := u.service.Register(c, &params)
	if err != nil {
		logger.Error(err.Error())
		response.Error(c, http.StatusInternalServerError, err, "Something was wrong")
		return
	}

	var res serializers.RegisterRes
	err = copier.Copy(&res.User, &user)
	if err != nil {
		logger.Error(err.Error())
		response.Error(c, http.StatusInternalServerError, err, "Something was wrong")
		return
	}
	response.JSON(c, http.StatusOK, res)
}

func (u *User) GetMe(c *gin.Context) {
	userID := c.GetString("userId")
	user, err := u.service.GetUserByID(c, userID)
	if err != nil {
		logger.Error(err.Error())
		response.Error(c, http.StatusInternalServerError, err, "Something was wrong")
		return
	}

	var res serializers.User
	err = copier.Copy(&res, &user)
	if err != nil {
		logger.Error(err.Error())
		response.Error(c, http.StatusInternalServerError, err, "Something was wrong")
		return
	}
	response.JSON(c, http.StatusOK, res)
}

func (u *User) RefreshToken(c *gin.Context) {
	userID := c.GetString("userId")
	accessToken, err := u.service.RefreshToken(c, userID)
	if err != nil {
		logger.Error("Failed to get body", err)
		response.Error(c, http.StatusBadRequest, err, "Something was wrong")
		return
	}

	res := serializers.RefreshTokenRes{
		AccessToken: accessToken,
	}
	response.JSON(c, http.StatusOK, res)
}
