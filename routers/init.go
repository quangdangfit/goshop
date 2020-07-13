package routers

import (
	"goshop/repositories"
	"goshop/services"
)

var roleService services.RoleService
var userService services.UserService
var categoryService services.CategorySerivce
var productService services.ProductService
var warehouseService services.WarehouseSerivce
var quantityService services.QuantitySerivce
var orderService services.OrderSerivce

func init() {
	roleRepo := repositories.NewRoleRepository()
	userRepo := repositories.NewUserRepository()
	categoryRepo := repositories.NewCategoryRepository()
	productRepo := repositories.NewProductRepository()
	warehouseRepo := repositories.NewWarehouseRepository()
	quantityRepo := repositories.NewQuantityRepository()
	orderRepo := repositories.NewOrderRepository()

	roleService = services.NewService(roleRepo)
	userService = services.NewUserService(userRepo)
	categoryService = services.NewCategoryService(categoryRepo)
	productService = services.NewProductService(productRepo)
	warehouseService = services.NewWarehouseService(warehouseRepo)
	quantityService = services.NewQuantityService(quantityRepo)
	orderService = services.NewOrderService(orderRepo)
}
