package api

import (
	"go.uber.org/dig"
)

func Inject(container *dig.Container) error {
	_ = container.Provide(NewProductAPI)
	_ = container.Provide(NewWarehouseAPI)
	_ = container.Provide(NewQuantityAPI)
	_ = container.Provide(NewUserAPI)
	_ = container.Provide(NewOrderAPI)
	return nil
}
