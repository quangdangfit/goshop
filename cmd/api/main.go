package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/quangdangfit/gocommon/logger"
	"github.com/quangdangfit/gocommon/validation"

	notificationModel "goshop/internal/notification/model"
	orderModel "goshop/internal/order/model"
	orderRepository "goshop/internal/order/repository"
	orderService "goshop/internal/order/service"
	paymentModel "goshop/internal/payment/model"
	productModel "goshop/internal/product/model"
	grpcServer "goshop/internal/server/grpc"
	httpServer "goshop/internal/server/http"
	userModel "goshop/internal/user/model"
	"goshop/pkg/config"
	"goshop/pkg/dbs"
	"goshop/pkg/notification"
	"goshop/pkg/redis"
)

//	@title			GoShop Swagger API
//	@version		1.0
//	@description	Swagger API for GoShop.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	Quang Dang
//	@contact.email	quangdangfit@gmail.com

//	@license.name	MIT
//	@license.url	https://github.com/MartinHeinz/go-project-blueprint/blob/master/LICENSE

//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						Authorization

//	@BasePath	/api/v1

func main() {
	cfg := config.LoadConfig()
	logger.Initialize(cfg.Environment)

	db, err := dbs.NewDatabase(cfg.DatabaseURI)
	if err != nil {
		logger.Fatal("Cannot connect to database", err)
	}

	err = db.AutoMigrate(
		&userModel.User{}, &userModel.Address{}, &userModel.Wishlist{},
		&productModel.Category{}, &productModel.Product{}, &productModel.Review{},
		orderModel.Coupon{}, orderModel.Order{}, orderModel.OrderLine{}, orderModel.StockReservation{},
		&paymentModel.Payment{}, &paymentModel.ProviderEvent{},
		&notificationModel.Preference{},
	)
	if err != nil {
		logger.Fatal("Database migration fail", err)
	}

	validator := validation.New()

	cache := redis.New(redis.Config{
		Address:  cfg.RedisURI,
		Password: cfg.RedisPassword,
		Database: cfg.RedisDB,
	})

	httpSvr := httpServer.NewServer(validator, db, cache)
	grpcSvr := grpcServer.NewServer(validator, db, cache)

	go func() {
		if err := httpSvr.Run(); err != nil {
			logger.Fatal(err)
		}
	}()

	go func() {
		if err := grpcSvr.Run(); err != nil {
			logger.Fatal(err)
		}
	}()

	// Background sweeper: release expired stock reservations and cancel their unpaid orders.
	sweeperCtx, sweeperCancel := context.WithCancel(context.Background())
	defer sweeperCancel()
	go runReservationSweeper(sweeperCtx, validator, db)

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down servers...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	grpcSvr.Shutdown()

	if err := httpSvr.Shutdown(ctx); err != nil {
		logger.Error("HTTP server forced to shutdown: ", err)
	}

	sweeperCancel()
	logger.Info("Servers exited gracefully")
}

func runReservationSweeper(ctx context.Context, validator validation.Validation, db dbs.Database) {
	svc := orderService.NewOrderService(
		validator, db,
		orderRepository.NewOrderRepository(db),
		orderRepository.NewProductRepository(db),
		orderRepository.NewUserRepository(db),
		orderRepository.NewReservationRepository(db),
		orderService.NewCouponService(validator, orderRepository.NewCouponRepository(db)),
		notification.NewLoggerNotifier(),
	)
	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			released, err := svc.SweepExpiredReservations(ctx, 100)
			if err != nil {
				logger.Error("reservation sweeper: ", err)
				continue
			}
			if released > 0 {
				logger.Infof("reservation sweeper released %d expired reservations", released)
			}
		}
	}
}
