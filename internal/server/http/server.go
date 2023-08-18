package http

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/quangdangfit/gocommon/logger"
	"github.com/quangdangfit/gocommon/redis"
	"github.com/quangdangfit/gocommon/validation"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "goshop/docs"

	"goshop/config"
	"goshop/pkg/response"
)

type Server struct {
	engine *gin.Engine
	cfg    *config.Schema
}

func NewServer() *Server {
	return &Server{
		engine: gin.Default(),
		cfg:    config.GetConfig(),
	}
}

func (s Server) Run(validator validation.Validation, cache redis.IRedis) error {
	_ = s.engine.SetTrustedProxies(nil)

	if err := s.MapRoutes(s.engine, validator, cache); err != nil {
		log.Fatalf("MapRoutes Error: %v", err)
	}
	s.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	s.engine.GET("/health", func(c *gin.Context) {
		response.JSON(c, http.StatusOK, nil)
		return
	})

	// Start server
	logger.Info("Server is listening on PORT: ", s.cfg.Port)
	if err := s.engine.Run(fmt.Sprintf(":%d", s.cfg.Port)); err != nil {
		log.Fatalf("Running HTTP server: %v", err)
	}

	return nil
}
