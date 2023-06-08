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

	if err := container.Provide(func() IOrderLineRepository { return NewOrderLineRepository() }); err != nil {
		return err
	}

	if err := container.Provide(func() IQuantityRepository { return NewQuantityRepository() }); err != nil {
		return err
	}

	if err := container.Provide(func() IUserRepository { return NewUserRepository() }); err != nil {
		return err
	}

	if err := container.Provide(func() IWarehouseRepository { return NewWarehouseRepository() }); err != nil {
		return err
	}
	return nil
}
