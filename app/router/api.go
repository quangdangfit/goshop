package router

import (
	"github.com/gin-gonic/gin"
	"github.com/quangdangfit/gocommon/logger"
	"go.uber.org/dig"

	"goshop/app/api"
	"goshop/app/middleware"
)

func RegisterRoute(r *gin.Engine, container *dig.Container) error {
	err := container.Invoke(func(
		user *api.UserAPI,
		product *api.ProductAPI,
		order *api.OrderAPI,
	) error {
		authMiddleware := middleware.JWTAuth()
		refreshAuthMiddleware := middleware.JWTRefresh()
		authRoute := r.Group("/auth")
		{
			authRoute.POST("/register", user.Register)
			authRoute.POST("/login", user.Login)
			authRoute.POST("/refresh", refreshAuthMiddleware, user.RefreshToken)
			authRoute.GET("/me", authMiddleware, user.GetMe)
			authRoute.PUT("/change-password", authMiddleware, user.ChangePassword)
		}

		//--------------------------------API-----------------------------------
		api1 := r.Group("/api/v1")

		// Products
		productAPI := api1.Group("/products")
		{
			productAPI.GET("", product.ListProducts)
			productAPI.POST("", authMiddleware, product.CreateProduct)
			productAPI.PUT("/:id", authMiddleware, product.UpdateProduct)
			productAPI.GET("/:id", product.GetProductByID)
		}

		// Orders
		orderAPI := api1.Group("/orders", authMiddleware)
		{
			orderAPI.POST("", order.PlaceOrder)
			orderAPI.GET("/:id", order.GetOrderByID)
			orderAPI.GET("", order.GetOrders)
			orderAPI.PUT("/:id/cancel", order.CancelOrder)
		}

		return nil
	})

	if err != nil {
		logger.Error(err)
	}

	return err
}
