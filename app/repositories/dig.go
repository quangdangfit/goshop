package repositories

import (
	"go.uber.org/dig"
)

func Inject(container *dig.Container) {
	_ = container.Provide(func() IProductRepository { return NewProductRepository() })
	_ = container.Provide(func() IOrderRepository { return NewOrderRepository() })
	_ = container.Provide(func() IUserRepository { return NewUserRepository() })
}
