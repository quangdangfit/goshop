package app

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/dig"

	"goshop/app/api"
	"goshop/app/repositories"
	"goshop/app/services"
	"goshop/config"
)

func BuildContainer() *dig.Container {
	container := dig.New()

	// Inject repositories
	repositories.Inject(container)

	// Inject services
	services.Inject(container)

	// Inject APIs
	api.Inject(container)

	return container
}

func InitGinEngine(container *dig.Container) *gin.Engine {
	cfg := config.GetConfig()
	if cfg.Environment == config.ProductionEnv {
		gin.SetMode(gin.ReleaseMode)
	}
	app := gin.Default()
	api.RegisterAPI(app, container)
	return app
}
