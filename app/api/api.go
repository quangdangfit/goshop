package api

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"goshop/app/middleware"

	_ "goshop/docs"
)


func RegisterAPI(r *gin.Engine, userAPI *UserAPI, productAPI *ProductAPI, orderAPI *OrderAPI) {
	r.GET("swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	authMiddleware := middleware.JWTAuth()
	refreshAuthMiddleware := middleware.JWTRefresh()
	authRoute := r.Group("/auth")
	{
		authRoute.POST("/register", userAPI.Register)
		authRoute.POST("/login", userAPI.Login)
		authRoute.POST("/refresh", refreshAuthMiddleware, userAPI.RefreshToken)
		authRoute.GET("/me", authMiddleware, userAPI.GetMe)
		authRoute.PUT("/change-password", authMiddleware, userAPI.ChangePassword)
	}

	//--------------------------------API-----------------------------------
	api1 := r.Group("/api/v1")

	// Products
	productRoute := api1.Group("/products")
	{
		productRoute.GET("", productAPI.ListProducts)
		productRoute.POST("", authMiddleware, productAPI.CreateProduct)
		productRoute.PUT("/:id", authMiddleware, productAPI.UpdateProduct)
		productRoute.GET("/:id", productAPI.GetProductByID)
	}

	// Orders
	orderRoute := api1.Group("/orders", authMiddleware)
	{
		orderRoute.POST("", orderAPI.PlaceOrder)
		orderRoute.GET("/:id", orderAPI.GetOrderByID)
		orderRoute.GET("", orderAPI.GetOrders)
		orderRoute.PUT("/:id/cancel", orderAPI.CancelOrder)
	}
}
