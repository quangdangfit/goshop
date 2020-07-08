package services

import "goshop/repositories"

var CategoryService Category
var ProductService Product

func init() {
	categoryRepo := repositories.NewCategoryRepository()
	productRepo := repositories.NewProductRepository()

	CategoryService = NewCategoryService(categoryRepo)
	ProductService = NewProductService(productRepo)
}
