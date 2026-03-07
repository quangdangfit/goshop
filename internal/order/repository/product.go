package repository

import (
	"context"

	"goshop/internal/order/model"
	"goshop/pkg/dbs"
)

//go:generate mockery --name=ProductRepository
type ProductRepository interface {
	GetProductByID(ctx context.Context, id string) (*model.Product, error)
}

type productRepo struct {
	db dbs.Database
}

func NewProductRepository(db dbs.Database) ProductRepository {
	return &productRepo{db: db}
}

func (r *productRepo) GetProductByID(ctx context.Context, id string) (*model.Product, error) {
	var product model.Product
	if err := r.db.FindById(ctx, id, &product); err != nil {
		return nil, err
	}

	return &product, nil
}
