package routers

import (
	"github.com/gin-gonic/gin"
	"goshop/objects/product"
)

func API(e *gin.Engine) {
	v1 := e.Group("api/v1")
	{
		proService := product.NewService()
		v1.GET("/products", proService.GetProducts)
		v1.GET("/products/:uuid", proService.GetProducts)
	}
}
