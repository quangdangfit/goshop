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

	"goshop/app"
	"goshop/app/api"
	"goshop/app/dbs"
	"goshop/app/repositories"
	"goshop/app/services"
	"goshop/config"
)

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

	engine := app.InitGinEngine(userAPI, productAPI, orderAPI)

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
