package router

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/dig"
)

func InitGinEngine(container *dig.Container) *gin.Engine {
	app := gin.New()
	RegisterAPI(app, container)
	return app
}
