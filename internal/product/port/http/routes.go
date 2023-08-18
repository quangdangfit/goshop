package http

import (
	"github.com/gin-gonic/gin"
	"github.com/quangdangfit/gocommon/redis"
	"github.com/quangdangfit/gocommon/validation"

	"goshop/app/middleware"
	"goshop/internal/product/repository"
	"goshop/internal/product/service"
)

func Routes(r *gin.RouterGroup, validator validation.Validation, cache redis.IRedis) {
	productRepo := repository.NewProductRepository()
	productSvc := service.NewProductService(validator, productRepo)
	productHandler := NewProductHandler(cache, productSvc)

	authMiddleware := middleware.JWTAuth()

	productRoute := r.Group("/products")
	{
		productRoute.GET("", productHandler.ListProducts)
		productRoute.POST("", authMiddleware, productHandler.CreateProduct)
		productRoute.PUT("/:id", authMiddleware, productHandler.UpdateProduct)
		productRoute.GET("/:id", productHandler.GetProductByID)
	}
}
