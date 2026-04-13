package http

import (
	"github.com/gin-gonic/gin"
	"github.com/quangdangfit/gocommon/validation"

	cartRepository "goshop/internal/cart/repository"
	"goshop/internal/order/repository"
	"goshop/internal/order/service"
	"goshop/pkg/dbs"
	"goshop/pkg/middleware"
	"goshop/pkg/notification"
)

func Routes(r *gin.RouterGroup, db dbs.Database, validator validation.Validation) {
	productRepo := repository.NewProductRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	couponRepo := repository.NewCouponRepository(db)
	userRepo := repository.NewUserRepository(db)
	cartRepo := cartRepository.NewCartRepository(db)

	couponSvc := service.NewCouponService(validator, couponRepo)
	notifier := notification.NewLoggerNotifier()

	orderSvc := service.NewOrderService(validator, orderRepo, productRepo, userRepo, cartRepo, couponSvc, notifier)
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
