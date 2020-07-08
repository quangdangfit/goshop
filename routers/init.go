package routers

import (
	"goshop/repositories"
	"goshop/services"
)

var roleService services.RoleService
var userService services.User
var categoryService services.Category
var productService services.Product

func init() {
	roleRepo := repositories.NewRoleRepository()
	userRepo := repositories.NewUserRepository()
	categoryRepo := repositories.NewCategoryRepository()
	productRepo := repositories.NewProductRepository()

	roleService = services.NewService(roleRepo)
	userService = services.NewUserService(userRepo)
	categoryService = services.NewCategoryService(categoryRepo)
	productService = services.NewProductService(productRepo)
}
