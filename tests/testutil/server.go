package testutil

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/quangdangfit/gocommon/logger"
	"github.com/quangdangfit/gocommon/validation"

	notificationModel "goshop/internal/notification/model"
	orderModel "goshop/internal/order/model"
	paymentModel "goshop/internal/payment/model"
	productModel "goshop/internal/product/model"
	httpServer "goshop/internal/server/http"
	userModel "goshop/internal/user/model"
	"goshop/pkg/config"
	"goshop/pkg/dbs"
	"goshop/pkg/redis"
)

// HTTPTestEnv bundles the live components an HTTP integration suite needs. Cleanup
// terminates the Postgres + Redis containers; call it from TestMain after m.Run.
type HTTPTestEnv struct {
	Engine  *gin.Engine
	DB      dbs.Database
	Cache   redis.Redis
	Cleanup func()
}

// NewHTTPEnv boots Postgres + Redis containers, auto-migrates all domain models, and wires
// the full Gin router (matching cmd/api startup). Intended for use from TestMain.
func NewHTTPEnv(ctx context.Context) (*HTTPTestEnv, error) {
	gin.SetMode(gin.TestMode)
	logger.Initialize(config.ProductionEnv)

	db, dbCleanup, err := StartPostgresM(ctx)
	if err != nil {
		return nil, fmt.Errorf("start postgres: %w", err)
	}

	cache, cacheCleanup, err := StartRedisM(ctx)
	if err != nil {
		dbCleanup()
		return nil, fmt.Errorf("start redis: %w", err)
	}

	if err := db.AutoMigrate(
		&userModel.User{}, &userModel.Address{}, &userModel.Wishlist{},
		&productModel.Category{}, &productModel.Product{}, &productModel.Review{},
		&orderModel.Coupon{}, &orderModel.Order{}, &orderModel.OrderLine{}, &orderModel.StockReservation{},
		&paymentModel.Payment{}, &paymentModel.ProviderEvent{},
		&notificationModel.Preference{}, &notificationModel.DeadLetterNotification{},
	); err != nil {
		cacheCleanup()
		dbCleanup()
		return nil, fmt.Errorf("auto-migrate: %w", err)
	}

	server := httpServer.NewServer(validation.New(), db, cache)
	if err := server.MapRoutes(); err != nil {
		cacheCleanup()
		dbCleanup()
		return nil, fmt.Errorf("map routes: %w", err)
	}

	return &HTTPTestEnv{
		Engine: server.GetEngine(),
		DB:     db,
		Cache:  cache,
		Cleanup: func() {
			cacheCleanup()
			dbCleanup()
		},
	}, nil
}
