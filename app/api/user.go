package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/quangdangfit/gocommon/logger"

	"goshop/app/schema"
	"goshop/app/services"
	"goshop/pkg/utils"
)

type User struct {
	service services.IUserService
}

func NewUserAPI(service services.IUserService) *User {
	return &User{service: service}
}

func (u *User) validate(r schema.Register) bool {
	return utils.Validate(
		[]utils.Validation{
			{Value: r.Username, Valid: "username"},
			{Value: r.Email, Valid: "email"},
			{Value: r.Password, Valid: "password"},
		})
}

func (u *User) checkPermission(uuid string, data map[string]interface{}) bool {
	return data["uuid"] == uuid
}

func (u *User) Login(c *gin.Context) {
	var item schema.Login
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	user, token, err := u.service.Login(ctx, &item)
	if err != nil {
		logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, utils.PrepareResponse(nil, err.Error(), ""))
		return
	}

	var res schema.User
	copier.Copy(&res, &user)
	res.Extra = map[string]interface{}{
		"token": token,
	}
	c.JSON(http.StatusOK, utils.PrepareResponse(res, "OK", ""))
}

func (u *User) Register(c *gin.Context) {
	var item schema.Register
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	valid := u.validate(item)
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request body is invalid"})
		return
	}

	ctx := c.Request.Context()
	user, token, err := u.service.Register(ctx, &item)
	if err != nil {
		logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, utils.PrepareResponse(nil, err.Error(), ""))
		return
	}

	var res schema.User
	copier.Copy(&res, &user)
	res.Extra = map[string]interface{}{
		"token": token,
	}
	c.JSON(http.StatusOK, utils.PrepareResponse(res, "OK", ""))
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

	var res schema.User
	copier.Copy(&res, &user)
	c.JSON(http.StatusOK, utils.PrepareResponse(res, "OK", ""))
}
