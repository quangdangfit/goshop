package api

import (
	"go.uber.org/dig"
)

func Inject(container *dig.Container) error {
	_ = container.Provide(NewCategoryAPI)
	_ = container.Provide(NewProductAPI)
	return nil
}
