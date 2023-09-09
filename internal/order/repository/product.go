package repository

import (
	"context"

	"goshop/internal/order/model"
	"goshop/pkg/dbs"
)

//go:generate mockery --name=IProductRepository
type IProductRepository interface {
	GetProductByID(ctx context.Context, id string) (*model.Product, error)
}

type ProductRepo struct {
	db dbs.IDatabase
}

func NewProductRepository(db dbs.IDatabase) *ProductRepo {
	return &ProductRepo{db: db}
}

func (r *ProductRepo) GetProductByID(ctx context.Context, id string) (*model.Product, error) {
	var product model.Product
	if err := r.db.FindById(ctx, id, &product); err != nil {
		return nil, err
	}

	return &product, nil
}
