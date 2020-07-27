package services

import (
	"go.uber.org/dig"
)

func Inject(container *dig.Container) error {
	_ = container.Provide(NewCategoryService)
	_ = container.Provide(NewOrderService)
	_ = container.Provide(NewProductService)
	_ = container.Provide(NewQuantityService)
	_ = container.Provide(NewUserService)
	_ = container.Provide(NewWarehouseService)
	return nil
}
