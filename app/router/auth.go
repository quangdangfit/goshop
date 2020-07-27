package router

import (
	"github.com/gin-gonic/gin"
)

func Auth(e *gin.Engine) {
	e.POST("auth/register", userService.Register)
	e.POST("auth/login", userService.Login)
}
