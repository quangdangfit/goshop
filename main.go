package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/quangdangfit/gocommon/logger"

	"goshop/app"
	"goshop/app/migrations"
	"goshop/app/router"
	"goshop/config"
	"goshop/dbs"
)

func main() {
	cfg := config.GetConfig()
	logger.Initialize(cfg.Environment)

	dbs.Init()

	migrations.Migrate()

	container := app.BuildContainer()
	engine := router.InitGinEngine(container)

	server := &http.Server{
		Addr:    ":8888",
		Handler: engine,
	}

	go func() {
		// service connections
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Listen: %s\n", err)
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
