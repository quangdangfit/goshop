package router

import (
	"github.com/gin-gonic/gin"

	"goshop/app/middleware/roles"
)

func Admin(e *gin.Engine) {
	admin := e.Group("admin")
	admin.Use(roles.CheckAdmin())
	{
		admin.POST("/roles", roleService.CreateRole)
	}
}
