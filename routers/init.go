package routers

import (
	"goshop/repositories"
	"goshop/services"
)

var userService services.User
var categoryService services.Category
var productService services.Product

func init() {
	userRepo := repositories.NewUserRepository()
	categoryRepo := repositories.NewCategoryRepository()
	productRepo := repositories.NewProductRepository()

	userService = services.NewUserService(userRepo)
	categoryService = services.NewCategoryService(categoryRepo)
	productService = services.NewProductService(productRepo)
}
