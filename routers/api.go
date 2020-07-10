package routers

import (
	"github.com/gin-gonic/gin"

	"goshop/middleware/jwt"
)

func API(e *gin.Engine) {
	apiV1 := e.Group("api/v1")
	apiV1.Use(jwt.JWT())
	//if config.Config.Redis.Enable {
	//	apiV1.Use(cache.Cached())
	//}

	{
		apiV1.GET("/users/:uuid", userService.GetUserByID)

		apiV1.GET("/products", productService.GetProducts)
		apiV1.POST("/products", productService.CreateProduct)
		apiV1.GET("/products/:uuid", productService.GetProductByID)
		apiV1.PUT("/products/:uuid", productService.UpdateProduct)

		apiV1.GET("/categories", categoryService.GetCategories)
		apiV1.POST("/categories", categoryService.CreateCategory)
		apiV1.GET("/categories/:uuid", categoryService.GetCategoryByID)
		apiV1.GET("/categories/:uuid/products", productService.GetProductByCategory)
		apiV1.PUT("/categories/:uuid", categoryService.UpdateCategory)
	}
}
