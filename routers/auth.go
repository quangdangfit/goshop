package routers

import (
	"github.com/gin-gonic/gin"
	"goshop/objects/user"
)

func Auth(e *gin.Engine) {
	userService := user.NewService()
	e.POST("auth/register", userService.Register)
	e.POST("auth/login", userService.Login)
}
