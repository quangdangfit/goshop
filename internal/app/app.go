package app

import (
	"github.com/gin-gonic/gin"

	"goshop/internal/app/api"
	"goshop/internal/config"
)

func InitGinEngine(
	userAPI *api.UserAPI,
	productAPI *api.ProductAPI,
	orderAPI *api.OrderAPI,
) *gin.Engine {
	cfg := config.GetConfig()
	if cfg.Environment == config.ProductionEnv {
		gin.SetMode(gin.ReleaseMode)
	}
	app := gin.Default()
	api.RegisterAPI(app, userAPI, productAPI, orderAPI)
	return app
}
