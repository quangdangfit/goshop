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
		product *api.ProductAPI,
		warehouse *api.Warehouse,
		quantity *api.Quantity,
		user *api.UserAPI,
		order *api.Order,
	) error {
		authMiddleware := middleware.JWTAuth()
		refreshAuthMiddleware := middleware.JWTRefresh()
		authRoute := r.Group("/auth")
		{
			authRoute.POST("/register", user.Register)
			authRoute.POST("/login", user.Login)
			authRoute.POST("/refresh", refreshAuthMiddleware, user.RefreshToken)
			authRoute.GET("/me", authMiddleware, user.GetMe)
		}

		//--------------------------------API-----------------------------------
		api1 := r.Group("/api/v1")

		// Products
		productAPI := api1.Group("/products")
		{
			productAPI.GET("", product.ListProducts)
			productAPI.POST("", product.CreateProduct)
			productAPI.PUT("/:id", product.UpdateProduct)
			productAPI.GET("/:id", product.GetProductByID)
		}

		{
			api1.GET("/warehouses", warehouse.GetWarehouses)
			api1.POST("/warehouses", warehouse.CreateWarehouse)
			api1.GET("/warehouses/:uuid", warehouse.GetWarehouseByID)
			api1.PUT("/warehouses/:uuid", warehouse.UpdateWarehouse)
		}
		{
			api1.GET("/quantities", quantity.GetQuantities)
			api1.POST("/quantities", quantity.CreateQuantity)
			api1.GET("/quantities/:uuid", quantity.GetQuantityByID)
			api1.PUT("/quantities/:uuid", quantity.UpdateQuantity)
		}
		{
			api1.GET("/orders", order.GetOrders)
			api1.POST("/orders", order.CreateOrder)
			api1.GET("/orders/:uuid", order.GetOrderByID)
			api1.PUT("/orders/:uuid", order.UpdateOrder)
		}

		return nil
	})

	if err != nil {
		logger.Error(err)
	}

	return err
}
