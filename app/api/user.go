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

type UserAPI struct {
	validator validation.Validation
	service   services.IUserService
}

func NewUserAPI(service services.IUserService) *UserAPI {
	return &UserAPI{
		validator: validation.New(),
		service:   service,
	}
}

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
		response.Error(c, http.StatusBadRequest, err, "Something went wrong")
		return
	}

	var res serializers.LoginRes
	err = copier.Copy(&res.User, &user)
	if err != nil {
		logger.Error(err.Error())
		response.Error(c, http.StatusInternalServerError, err, "Something went wrong")
		return
	}
	res.AccessToken = accessToken
	res.RefreshToken = refreshToken
	response.JSON(c, http.StatusOK, res)
}

// Register godoc
//
//	@Summary	Register new user
//	@Produce	json
//	@Param		b	body		serializers.RegisterReq	true	"Body"
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
	err = copier.Copy(&res.User, &user)
	if err != nil {
		logger.Error(err.Error())
		response.Error(c, http.StatusInternalServerError, err, "Something went wrong")
		return
	}
	response.JSON(c, http.StatusOK, res)
}

func (u *UserAPI) GetMe(c *gin.Context) {
	userID := c.GetString("userId")
	user, err := u.service.GetUserByID(c, userID)
	if err != nil {
		logger.Error(err.Error())
		response.Error(c, http.StatusInternalServerError, err, "Something went wrong")
		return
	}

	var res serializers.User
	err = copier.Copy(&res, &user)
	if err != nil {
		logger.Error(err.Error())
		response.Error(c, http.StatusInternalServerError, err, "Something went wrong")
		return
	}
	response.JSON(c, http.StatusOK, res)
}

func (u *UserAPI) RefreshToken(c *gin.Context) {
	userID := c.GetString("userId")
	accessToken, err := u.service.RefreshToken(c, userID)
	if err != nil {
		logger.Error("Failed to get body", err)
		response.Error(c, http.StatusBadRequest, err, "Something went wrong")
		return
	}

	res := serializers.RefreshTokenRes{
		AccessToken: accessToken,
	}
	response.JSON(c, http.StatusOK, res)
}
