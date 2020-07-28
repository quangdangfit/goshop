package router

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/dig"
)

func InitGinEngine(container *dig.Container) *gin.Engine {
	app := gin.New()
	Docs(app)
	RegisterAPI(app, container)
	return app
}
