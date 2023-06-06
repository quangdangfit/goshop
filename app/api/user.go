package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/quangdangfit/gocommon/logger"

	"goshop/app/serializers"
	"goshop/app/services"
	"goshop/pkg/response"
	"goshop/pkg/utils"
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

func (u *User) checkPermission(uuid string, data map[string]interface{}) bool {
	return data["uuid"] == uuid
}

func (u *User) Login(c *gin.Context) {
	var item serializers.Login
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	user, _, err := u.service.Login(ctx, &item)
	if err != nil {
		logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, utils.PrepareResponse(nil, err.Error(), ""))
		return
	}

	var res serializers.User
	copier.Copy(&res, &user)
	c.JSON(http.StatusOK, utils.PrepareResponse(res, "OK", ""))
}

func (u *User) Register(c *gin.Context) {
	var params serializers.RegisterReq
	if err := c.ShouldBindJSON(&params); c.Request.Body == nil || err != nil {
		logger.Error("Failed to get body", err)
		response.Error(c, http.StatusBadRequest, err, "Invalid Parameters")
		return
	}

	if err := u.validator.ValidateStruct(params); err != nil {
		response.Error(c, http.StatusBadRequest, err, "Invalid Parameters")
		return
	}

	ctx := c.Request.Context()
	user, token, err := u.service.Register(ctx, &params)
	if err != nil {
		logger.Error(err.Error())
		response.Error(c, http.StatusInternalServerError, err, "Internal Server Error")
		return
	}

	var res serializers.RegisterRes
	copier.Copy(&res.User, &user)
	res.AccessToken = token

	response.JSON(c, http.StatusOK, res)
}

func (u *User) GetUserByID(c *gin.Context) {
	userUUID := c.Param("uuid")
	ctx := c.Request.Context()
	user, err := u.service.GetUserByID(ctx, userUUID)
	if err != nil {
		logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, utils.PrepareResponse(nil, err.Error(), utils.ErrorNotExistUser))
		return
	}

	var res serializers.User
	copier.Copy(&res, &user)
	c.JSON(http.StatusOK, utils.PrepareResponse(res, "OK", ""))
}
