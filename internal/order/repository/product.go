package repository

import (
	"context"

	"gorm.io/gorm"

	"goshop/app/dbs"
	"goshop/config"
	"goshop/internal/order/model"
)

type IProductRepository interface {
	GetProductByID(ctx context.Context, id string) (*model.Product, error)
}

type ProductRepo struct {
	db *gorm.DB
}

func NewProductRepository() *ProductRepo {
	return &ProductRepo{db: dbs.Database}
}

func (r *ProductRepo) GetProductByID(ctx context.Context, id string) (*model.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, config.DatabaseTimeout)
	defer cancel()

	var product model.Product
	if err := r.db.Where("id = ?", id).First(&product).Error; err != nil {
		return nil, err
	}

	return &product, nil
}
