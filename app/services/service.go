package services

import (
	"go.uber.org/dig"

	"goshop/app/repositories"
)

func Inject(container *dig.Container) {
	_ = container.Provide(func(
		repo repositories.IOrderRepository,
		productRepo repositories.IProductRepository,
	) IOrderService {
		return NewOrderService(repo, productRepo)
	})
	_ = container.Provide(func(
		repo repositories.IProductRepository,
	) IProductService {
		return NewProductService(repo)
	})
	_ = container.Provide(func(
		repo repositories.IUserRepository,
	) IUserService {
		return NewUserService(repo)
	})
}
