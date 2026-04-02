package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/quangdangfit/gocommon/logger"

	"goshop/internal/user/domain"
	"goshop/internal/user/service"
	"goshop/pkg/apperror"
	"goshop/pkg/response"
	"goshop/pkg/utils"
)

type UserHandler struct {
	service service.UserService
}

func NewUserHandler(service service.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

// Login godoc
//
//	@Summary	Login
//	@Tags		users
//	@Produce	json
//	@Param		_	body		domain.LoginReq	true	"Body"
//	@Success	200	{object}	domain.LoginRes
//	@Router		/api/v1/auth/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req domain.LoginReq
	if err := c.ShouldBindJSON(&req); c.Request.Body == nil || err != nil {
		logger.Error("Failed to get body ", err)
		apperror.Wrap(apperror.ErrBadRequest, err).HTTPError(c)
		return
	}

	user, accessToken, refreshToken, err := h.service.Login(c, &req)
	if err != nil {
		logger.Error("Failed to login ", err)
		apperror.ToHTTPError(c, err, http.StatusInternalServerError, "Something went wrong")
		return
	}

	var res domain.LoginRes
	utils.Copy(&res.User, &user)
	res.AccessToken = accessToken
	res.RefreshToken = refreshToken
	response.JSON(c, http.StatusOK, res)
}

// Register godoc
//
//	@Summary	Register new user
//	@Tags		users
//	@Produce	json
//	@Param		_	body		domain.RegisterReq	true	"Body"
//	@Success	200	{object}	domain.RegisterRes
//	@Router		/api/v1/auth/register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var req domain.RegisterReq
	if err := c.ShouldBindJSON(&req); c.Request.Body == nil || err != nil {
		logger.Error("Failed to get body", err)
		apperror.Wrap(apperror.ErrBadRequest, err).HTTPError(c)
		return
	}

	user, err := h.service.Register(c, &req)
	if err != nil {
		logger.Error(err.Error())
		apperror.ToHTTPError(c, err, http.StatusInternalServerError, "Something went wrong")
		return
	}

	var res domain.RegisterRes
	utils.Copy(&res.User, &user)
	response.JSON(c, http.StatusOK, res)
}

// GetMe godoc
//
//	@Summary	get my profile
//	@Tags		users
//	@Security	ApiKeyAuth
//	@Produce	json
//	@Success	200	{object}	domain.User
//	@Router		/api/v1/auth/me [get]
func (h *UserHandler) GetMe(c *gin.Context) {
	userID := c.GetString("userId")
	if userID == "" {
		apperror.ErrUnauthorized.HTTPError(c)
		return
	}

	user, err := h.service.GetUserByID(c, userID)
	if err != nil {
		logger.Error(err.Error())
		apperror.ToHTTPError(c, err, http.StatusInternalServerError, "Something went wrong")
		return
	}

	var res domain.User
	utils.Copy(&res, &user)
	response.JSON(c, http.StatusOK, res)
}

func (h *UserHandler) RefreshToken(c *gin.Context) {
	userID := c.GetString("userId")
	if userID == "" {
		apperror.ErrUnauthorized.HTTPError(c)
		return
	}

	accessToken, err := h.service.RefreshToken(c, userID)
	if err != nil {
		logger.Error("Failed to refresh token", err)
		apperror.ToHTTPError(c, err, http.StatusInternalServerError, "Something went wrong")
		return
	}

	res := domain.RefreshTokenRes{
		AccessToken: accessToken,
	}
	response.JSON(c, http.StatusOK, res)
}

// ChangePassword godoc
//
//	@Summary	changes the password
//	@Tags		users
//	@Security	ApiKeyAuth
//	@Produce	json
//	@Param		_	body	domain.ChangePasswordReq	true	"Body"
//	@Router		/api/v1/auth/change-password [put]
func (h *UserHandler) ChangePassword(c *gin.Context) {
	var req domain.ChangePasswordReq
	if err := c.ShouldBindJSON(&req); c.Request.Body == nil || err != nil {
		logger.Error("Failed to get body", err)
		apperror.Wrap(apperror.ErrBadRequest, err).HTTPError(c)
		return
	}

	userID := c.GetString("userId")
	err := h.service.ChangePassword(c, userID, &req)
	if err != nil {
		logger.Error(err.Error())
		apperror.ToHTTPError(c, err, http.StatusInternalServerError, "Something went wrong")
		return
	}
	response.JSON(c, http.StatusOK, nil)
}
