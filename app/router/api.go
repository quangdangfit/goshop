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
		warehouse *api.Warehouse,
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

		orderAPI := api1.Group("/orders", authMiddleware)
		{
			orderAPI.POST("", order.CreateOrder)
			orderAPI.GET("/:id", order.GetOrderByID)
			orderAPI.GET("", order.GetOrders)
			orderAPI.PUT("/:id", order.UpdateOrder)
		}

		{
			api1.GET("/warehouses", warehouse.GetWarehouses)
			api1.POST("/warehouses", warehouse.CreateWarehouse)
			api1.GET("/warehouses/:uuid", warehouse.GetWarehouseByID)
			api1.PUT("/warehouses/:uuid", warehouse.UpdateWarehouse)
		}

		return nil
	})

	if err != nil {
		logger.Error(err)
	}

	return err
}
