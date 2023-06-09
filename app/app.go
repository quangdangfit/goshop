package app

import (
	"github.com/gin-gonic/gin"
	"github.com/quangdangfit/gocommon/logger"
	"go.uber.org/dig"

	"goshop/app/api"
	"goshop/app/repositories"
	"goshop/app/services"
	"goshop/config"
)

func BuildContainer() *dig.Container {
	container := dig.New()

	// Inject repositories
	err := repositories.Inject(container)
	if err != nil {
		logger.Error("Failed to inject repositories", err)
	}

	// Inject services
	err = services.Inject(container)
	if err != nil {
		logger.Error("Failed to inject services", err)
	}

	// Inject APIs
	err = api.Inject(container)
	if err != nil {
		logger.Error("Failed to inject APIs", err)
	}

	return container
}

func InitGinEngine(container *dig.Container) *gin.Engine {
	cfg := config.GetConfig()
	if cfg.Environment == config.ProductionEnv {
		gin.SetMode(gin.ReleaseMode)
	}
	app := gin.Default()
	err := api.RegisterAPI(app, container)
	if err != nil {
		logger.Fatal("Failed to init GIN Engine", err)
	}
	return app
}
