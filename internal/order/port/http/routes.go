package http

import (
	"github.com/gin-gonic/gin"
	"github.com/quangdangfit/gocommon/validation"
	"gorm.io/gorm"

	"goshop/internal/order/repository"
	"goshop/internal/order/service"
	"goshop/pkg/middleware"
)

func Routes(r *gin.RouterGroup, db *gorm.DB, validator validation.Validation) {
	productRepo := repository.NewProductRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	productSvc := service.NewOrderService(validator, orderRepo, productRepo)
	orderHandler := NewOrderHandler(productSvc)

	authMiddleware := middleware.JWTAuth()

	orderRoute := r.Group("/orders", authMiddleware)
	{
		orderRoute.POST("", orderHandler.PlaceOrder)
		orderRoute.GET("/:id", orderHandler.GetOrderByID)
		orderRoute.GET("", orderHandler.GetOrders)
		orderRoute.PUT("/:id/cancel", orderHandler.CancelOrder)
	}
}
