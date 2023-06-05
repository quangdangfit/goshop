package router

import (
	"github.com/gin-gonic/gin"
	"github.com/quangdangfit/gocommon/logger"
	"go.uber.org/dig"

	"goshop/config"
)

func InitGinEngine(container *dig.Container) *gin.Engine {
	cfg := config.GetConfig()
	if cfg.Environment == config.ProductionEnv {
		gin.SetMode(gin.ReleaseMode)
	}
	app := gin.Default()
	Docs(app)
	err := RegisterRoute(app, container)
	if err != nil {
		logger.Fatal("Failed to init GIN Engine", err)
	}
	return app
}
