package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/quangdangfit/gocommon/logger"
	"github.com/quangdangfit/gocommon/redis"
	"github.com/quangdangfit/gocommon/validation"

	"goshop/internal/app"
	"goshop/internal/app/api"
	"goshop/internal/app/dbs"
	"goshop/internal/app/repositories"
	"goshop/internal/app/services"
	"goshop/internal/config"

	"github.com/gin-gonic/gin"
	"github.com/penglongli/gin-metrics/ginmetrics"
)


//	@Title			Blueprint Swagger API
//	@Version		1.0
//	@Description	Swagger API for Golang Project Blueprint.
//	@termsOfService	http://swagger.io/terms/
//	@contact.name	API Support
//	@contact.email	quangdangfit@gmail.com
//	@license.name	MIT
//	@license.url	https://github.com/MartinHeinz/go-project-blueprint/blob/master/LICENSE

//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						Authorization
//	@BasePath					/api/v1


func main() {
	cfg := config.GetConfig()
	logger.Initialize(cfg.Environment)

	dbs.Init()

	validator := validation.New()
	cache := redis.New(redis.Config{
		Address:  cfg.RedisURI,
		Password: cfg.RedisPassword,
		Database: cfg.RedisDB,
	})

	userRepo := repositories.NewUserRepository()
	productRepo := repositories.NewProductRepository()
	orderRepo := repositories.NewOrderRepository()

	userSvc := services.NewUserService(userRepo)
	productSvc := services.NewProductService(productRepo)
	orderSvc := services.NewOrderService(orderRepo, productRepo)

	userAPI := api.NewUserAPI(validator, userSvc)
	productAPI := api.NewProductAPI(validator, cache, productSvc)
	orderAPI := api.NewOrderAPI(validator, orderSvc)



        api := gin.Default()

        metricRouter := gin.Default()
        m := ginmetrics.GetMonitor()
        //m.UseWithoutExposingEndpoint(engine)
        m.Expose(metricRouter)
        // +optional set metric path, default /debug/metrics
        m.SetMetricPath("/metrics")
        // +optional set slow time, default 5s
        m.SetSlowTime(10)
        // +optional set request duration, default {0.1, 0.3, 1.2, 5, 10}
        // used to p95, p99
        //m.SetDuration([]float64{0.1, 0.3, 1.2, 5, 10})

        // set middleware for gin
        m.Use(api)

	engine := app.InitGinEngine(api, userAPI, productAPI, orderAPI)


	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: engine,
	}

	go func() {
		logger.Infof("Listen at: %d\n", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Failed to start server: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal("Server Shutdown: ", err)
	}
	// catching ctx.Done(). timeout of 5 seconds.
	select {
	case <-ctx.Done():
		logger.Info("Timeout of 5 seconds.")
	}
	logger.Info("Server exiting")
}
