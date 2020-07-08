package routers

import (
	"goshop/repositories"
	"goshop/services"
)

var categoryService services.Category
var productService services.Product

func init() {
	categoryRepo := repositories.NewCategoryRepository()
	productRepo := repositories.NewProductRepository()

	categoryService = services.NewCategoryService(categoryRepo)
	productService = services.NewProductService(productRepo)
}
