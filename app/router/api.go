package router

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/quangdangfit/gocommon/utils/logger"
	"go.uber.org/dig"

	"goshop/app/api"
)

func RegisterAPI(r *gin.Engine, container *dig.Container) error {
	err := container.Invoke(func(
		category *api.Category,
		product *api.Product,
		warehouse *api.Warehouse,
		quantity *api.Quantity,
		user *api.User,
	) error {
		auth := r.Group("/auth")
		{
			auth.POST("auth/register", user.Register)
			auth.POST("auth/login", user.Login)
		}

		apiV1 := r.Group("api/v1")
		{
			apiV1.GET("/users/:uuid", user.GetUserByID)
		}
		{
			apiV1.GET("/products", product.GetProducts)
			apiV1.POST("/products", product.CreateProduct)
			apiV1.GET("/products/:uuid", product.GetProductByID)
			apiV1.PUT("/products/:uuid", product.UpdateProduct)
		}
		{
			apiV1.GET("/categories", category.GetCategories)
			apiV1.POST("/categories", category.CreateCategory)
			apiV1.GET("/categories/:uuid", category.GetCategoryByID)
			apiV1.GET("/categories/:uuid/products", product.GetProductByCategoryID)
			apiV1.PUT("/categories/:uuid", category.UpdateCategory)
		}
		{
			apiV1.GET("/warehouses", warehouse.GetWarehouses)
			apiV1.POST("/warehouses", warehouse.CreateWarehouse)
			apiV1.GET("/warehouses/:uuid", warehouse.GetWarehouseByID)
			apiV1.PUT("/warehouses/:uuid", warehouse.UpdateWarehouse)
		}
		{
			apiV1.GET("/quantities", quantity.GetQuantities)
			apiV1.POST("/quantities", quantity.CreateQuantity)
			apiV1.GET("/quantities/:uuid", quantity.GetQuantityByID)
			apiV1.PUT("/quantities/:uuid", quantity.UpdateQuantity)
		}
		{
			apiV1.GET("/orders", orderService.GetOrders)
			apiV1.POST("/orders", orderService.CreateOrder)
			apiV1.GET("/orders/:uuid", orderService.GetOrderByID)
			apiV1.PUT("/orders/:uuid", orderService.UpdateOrder)
		}

		return nil
	})

	if err != nil {
		logger.Error(err)
	}

	return err
}
