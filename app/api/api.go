package api

import (
	"go.uber.org/dig"
)

func Inject(container *dig.Container) error {
	_ = container.Provide(NewCategoryAPI)
	_ = container.Provide(NewProductAPI)
	_ = container.Provide(NewWarehouseAPI)
	_ = container.Provide(NewQuantityAPI)
	_ = container.Provide(NewUserAPI)
	return nil
}
