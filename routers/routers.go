package routers

import (
	"github.com/gin-gonic/gin"
	"goshop/objects/category"
	"goshop/objects/product"
	"goshop/objects/user"
)

func API(e *gin.Engine) {
	v1 := e.Group("api/v1")
	{
		userService := user.NewService()
		v1.POST("/register", userService.Register)
		v1.POST("/login", userService.Login)

		proService := product.NewService()
		v1.GET("/products", proService.GetProducts)
		v1.POST("/products", proService.CreateProduct)
		v1.GET("/products/:uuid", proService.GetProductByID)
		v1.PUT("/products/:uuid", proService.UpdateProduct)

		categService := category.NewService()
		v1.GET("/categories", categService.GetCategories)
		v1.POST("/categories", categService.CreateCategory)
		v1.GET("/categories/:uuid", categService.GetCategoryByID)
		v1.GET("/categories/:uuid/products", proService.GetProductByCategory)
		v1.PUT("/categories/:uuid", categService.UpdateCategory)
	}
}
