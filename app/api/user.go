package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/quangdangfit/gocommon/logger"
	"github.com/quangdangfit/gocommon/validation"

	"goshop/app/serializers"
	"goshop/app/services"
	"goshop/pkg/response"
	"goshop/pkg/utils"
)

type UserAPI struct {
	validator validation.Validation
	service   services.IUserService
}

func NewUserAPI(validator validation.Validation, service services.IUserService) *UserAPI {
	return &UserAPI{
		validator: validator,
		service:   service,
	}
}

// Login godoc
//
//	@Summary	Login
//	@Produce	json
//	@Param		_	body		serializers.LoginReq	true	"Body"
//	@Success	200	{object}	serializers.LoginRes
//	@Router		/auth/login [post]
func (u *UserAPI) Login(c *gin.Context) {
	var req serializers.LoginReq
	if err := c.ShouldBindJSON(&req); c.Request.Body == nil || err != nil {
		logger.Error("Failed to get body", err)
		response.Error(c, http.StatusBadRequest, err, "Invalid parameters")
		return
	}

	if err := u.validator.ValidateStruct(req); err != nil {
		response.Error(c, http.StatusBadRequest, err, "Invalid parameters")
		return
	}

	user, accessToken, refreshToken, err := u.service.Login(c, &req)
	if err != nil {
		logger.Error("Failed to get body", err)
		response.Error(c, http.StatusInternalServerError, err, "Something went wrong")
		return
	}

	var res serializers.LoginRes
	utils.Copy(&res.User, &user)
	res.AccessToken = accessToken
	res.RefreshToken = refreshToken
	response.JSON(c, http.StatusOK, res)
}

// Register godoc
//
//	@Summary	Register new user
//	@Produce	json
//	@Param		_	body		serializers.RegisterReq	true	"Body"
//	@Success	200	{object}	serializers.RegisterRes
//	@Router		/auth/register [post]
func (u *UserAPI) Register(c *gin.Context) {
	var req serializers.RegisterReq
	if err := c.ShouldBindJSON(&req); c.Request.Body == nil || err != nil {
		logger.Error("Failed to get body", err)
		response.Error(c, http.StatusBadRequest, err, "Invalid parameters")
		return
	}

	if err := u.validator.ValidateStruct(req); err != nil {
		response.Error(c, http.StatusBadRequest, err, "Invalid parameters")
		return
	}

	user, err := u.service.Register(c, &req)
	if err != nil {
		logger.Error(err.Error())
		response.Error(c, http.StatusInternalServerError, err, "Something went wrong")
		return
	}

	var res serializers.RegisterRes
	utils.Copy(&res.User, &user)
	response.JSON(c, http.StatusOK, res)
}

// GetMe godoc
//
//	@Summary	get my profile
//	@Security	ApiKeyAuth
//	@Produce	json
//	@Success	200	{object}	serializers.User
//	@Router		/auth/me [get]
func (u *UserAPI) GetMe(c *gin.Context) {
	userID := c.GetString("userId")
	user, err := u.service.GetUserByID(c, userID)
	if err != nil {
		logger.Error(err.Error())
		response.Error(c, http.StatusInternalServerError, err, "Something went wrong")
		return
	}

	var res serializers.User
	utils.Copy(&res, &user)
	response.JSON(c, http.StatusOK, res)
}

func (u *UserAPI) RefreshToken(c *gin.Context) {
	userID := c.GetString("userId")
	accessToken, err := u.service.RefreshToken(c, userID)
	if err != nil {
		logger.Error("Failed to refresh token", err)
		response.Error(c, http.StatusInternalServerError, err, "Something went wrong")
		return
	}

	res := serializers.RefreshTokenRes{
		AccessToken: accessToken,
	}
	response.JSON(c, http.StatusOK, res)
}

// ChangePassword godoc
//
//	@Summary	changes the password
//	@Security	ApiKeyAuth
//	@Produce	json
//	@Param		_	body	serializers.ChangePasswordReq	true	"Body"
//	@Router		/auth/change-password [put]
func (u *UserAPI) ChangePassword(c *gin.Context) {
	var req serializers.ChangePasswordReq
	if err := c.ShouldBindJSON(&req); c.Request.Body == nil || err != nil {
		logger.Error("Failed to get body", err)
		response.Error(c, http.StatusBadRequest, err, "Invalid parameters")
		return
	}

	if err := u.validator.ValidateStruct(req); err != nil {
		response.Error(c, http.StatusBadRequest, err, "Invalid parameters")
		return
	}

	userID := c.GetString("userId")
	err := u.service.ChangePassword(c, userID, &req)
	if err != nil {
		logger.Error(err.Error())
		response.Error(c, http.StatusInternalServerError, err, "Something went wrong")
		return
	}
	response.JSON(c, http.StatusOK, nil)
}
