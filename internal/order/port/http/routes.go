package http

import (
	"github.com/gin-gonic/gin"
	"github.com/quangdangfit/gocommon/validation"

	notificationRepo "goshop/internal/notification/repository"
	notificationSvc "goshop/internal/notification/service"
	"goshop/internal/order/repository"
	"goshop/internal/order/service"
	userRepository "goshop/internal/user/repository"
	"goshop/pkg/config"
	"goshop/pkg/dbs"
	"goshop/pkg/middleware"
	"goshop/pkg/notification"
)

func Routes(r *gin.RouterGroup, db dbs.Database, validator validation.Validation) {
	productRepo := repository.NewProductRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	couponRepo := repository.NewCouponRepository(db)
	userRepo := repository.NewUserRepository(db)
	reservationRepo := repository.NewReservationRepository(db)

	cfg := config.GetConfig()
	couponSvc := service.NewCouponService(validator, couponRepo)
	prefChecker := notificationSvc.NewDBPreferenceChecker(
		notificationSvc.NewUserRepoLookup(userRepository.NewUserRepository(db)),
		notificationRepo.NewPreferenceRepository(db),
	)
	notifier := notification.BuildDefault(notification.Settings{
		SMTPHost:     cfg.SMTPHost,
		SMTPPort:     cfg.SMTPPort,
		SMTPUser:     cfg.SMTPUser,
		SMTPPassword: cfg.SMTPPassword,
		EmailFrom:    cfg.EmailFrom,
		Prefs:        prefChecker,
	})

	orderSvc := service.NewOrderService(validator, db, orderRepo, productRepo, userRepo, reservationRepo, couponSvc, notifier)
	orderHandler := NewOrderHandler(orderSvc)
	couponHandler := NewCouponHandler(couponSvc)

	authMiddleware := middleware.JWTAuth()
	adminMiddleware := middleware.AdminOnly()

	orderRoute := r.Group("/orders", authMiddleware)
	{
		orderRoute.POST("", orderHandler.PlaceOrder)
		orderRoute.GET("/:id", orderHandler.GetOrderByID)
		orderRoute.GET("", orderHandler.GetOrders)
		orderRoute.PUT("/:id/cancel", orderHandler.CancelOrder)
		orderRoute.PUT("/:id/status", adminMiddleware, orderHandler.UpdateOrderStatus)
	}

	couponRoute := r.Group("/coupons", authMiddleware)
	{
		couponRoute.POST("", adminMiddleware, couponHandler.CreateCoupon)
		couponRoute.GET("/:code", couponHandler.GetCouponByCode)
	}
}
