package api

import (
	"go.uber.org/dig"
)

func Inject(container *dig.Container) {
	_ = container.Provide(NewProductAPI)
	_ = container.Provide(NewUserAPI)
	_ = container.Provide(NewOrderAPI)
}
