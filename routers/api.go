package routers

import (
	"github.com/gin-gonic/gin"
	"goshop/middlewares/jwt"
	"goshop/objects/category"
	"goshop/objects/product"
	"goshop/objects/user"
)

func API(e *gin.Engine) {
	apiV1 := e.Group("api/v1")
	apiV1.Use(jwt.JWT())
	{
		userService := user.NewService()
		apiV1.GET("/users/:uuid", userService.GetUserByID)

		proService := product.NewService()
		apiV1.GET("/products", proService.GetProducts)
		apiV1.POST("/products", proService.CreateProduct)
		apiV1.GET("/products/:uuid", proService.GetProductByID)
		apiV1.PUT("/products/:uuid", proService.UpdateProduct)

		categService := category.NewService()
		apiV1.GET("/categories", categService.GetCategories)
		apiV1.POST("/categories", categService.CreateCategory)
		apiV1.GET("/categories/:uuid", categService.GetCategoryByID)
		apiV1.GET("/categories/:uuid/products", proService.GetProductByCategory)
		apiV1.PUT("/categories/:uuid", categService.UpdateCategory)
	}
}
