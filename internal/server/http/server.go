package http

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/quangdangfit/gocommon/logger"
	"github.com/quangdangfit/gocommon/validation"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "goshop/docs"
	cartHttp "goshop/internal/cart/port/http"
	orderHttp "goshop/internal/order/port/http"
	productHttp "goshop/internal/product/port/http"
	userHttp "goshop/internal/user/port/http"
	"goshop/pkg/config"
	"goshop/pkg/dbs"
	"goshop/pkg/middleware"
	"goshop/pkg/redis"
	"goshop/pkg/response"
)

type Server struct {
	engine    *gin.Engine
	httpSvr   *http.Server
	cfg       *config.Schema
	validator validation.Validation
	db        dbs.Database
	cache     redis.Redis
}

func NewServer(validator validation.Validation, db dbs.Database, cache redis.Redis) *Server {
	return &Server{
		engine:    gin.Default(),
		cfg:       config.GetConfig(),
		validator: validator,
		db:        db,
		cache:     cache,
	}
}

func (s *Server) Run() error {
	_ = s.engine.SetTrustedProxies(nil)
	if s.cfg.Environment == config.ProductionEnv {
		gin.SetMode(gin.ReleaseMode)
	}

	s.engine.Use(middleware.CORS())
	s.engine.Use(middleware.RateLimit(s.cache))

	if err := s.MapRoutes(); err != nil {
		log.Fatalf("MapRoutes Error: %v", err)
	}
	s.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	s.engine.GET("/health", func(c *gin.Context) {
		response.JSON(c, http.StatusOK, nil)
	})

	s.httpSvr = &http.Server{
		Addr:              fmt.Sprintf(":%d", s.cfg.HttpPort),
		Handler:           s.engine,
		ReadHeaderTimeout: 10 * time.Second,
	}

	// Start http server
	logger.Info("HTTP server is listening on PORT: ", s.cfg.HttpPort)
	if err := s.httpSvr.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Running HTTP server: %v", err)
	}

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	logger.Info("Shutting down HTTP server...")
	return s.httpSvr.Shutdown(ctx)
}

func (s *Server) GetEngine() *gin.Engine {
	return s.engine
}

func (s *Server) MapRoutes() error {
	v1 := s.engine.Group("/api/v1")
	userHttp.Routes(v1, s.db, s.validator)
	productHttp.Routes(v1, s.db, s.validator, s.cache)
	orderHttp.Routes(v1, s.db, s.validator)
	cartHttp.Routes(v1, s.db, s.validator)
	return nil
}
