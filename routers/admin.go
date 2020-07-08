package routers

import (
	"github.com/gin-gonic/gin"
	"goshop/middlewares/roles"
	"goshop/objects/role"
)

func Admin(e *gin.Engine) {
	admin := e.Group("admin")
	admin.Use(roles.CheckAdmin())
	{
		roleService := role.NewService()
		admin.POST("/roles", roleService.CreateRole)
	}
}
