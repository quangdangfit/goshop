package repositories

import (
	"go.uber.org/dig"
)

func Inject(container *dig.Container) error {
	if err := container.Provide(func() IProductRepository { return NewProductRepository() }); err != nil {
		return err
	}

	if err := container.Provide(func() IOrderRepository { return NewOrderRepository() }); err != nil {
		return err
	}

	if err := container.Provide(func() IUserRepository { return NewUserRepository() }); err != nil {
		return err
	}

	return nil
}
