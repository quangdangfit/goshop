package http

import (
	"github.com/gin-gonic/gin"
	"github.com/quangdangfit/gocommon/validation"

	"goshop/internal/cart/repository"
	"goshop/internal/cart/service"
	"goshop/pkg/dbs"
	"goshop/pkg/middleware"
)

func Routes(r *gin.RouterGroup, db dbs.Database, validator validation.Validation) {
	cartRepo := repository.NewCartRepository(db)
	cartSvc := service.NewCartService(validator, cartRepo)
	cartHandler := NewCartHandler(cartSvc)

	authMiddleware := middleware.JWTAuth()

	cartRoute := r.Group("/cart", authMiddleware)
	{
		cartRoute.GET("", cartHandler.GetCart)
		cartRoute.POST("", cartHandler.AddProduct)
		cartRoute.DELETE("/:productId", cartHandler.RemoveProduct)
	}
}
