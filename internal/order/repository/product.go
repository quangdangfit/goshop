package repository

import (
	"context"

	"gorm.io/gorm"

	"goshop/config"
	"goshop/internal/order/model"
)

//go:generate mockery --name=IProductRepository
type IProductRepository interface {
	GetProductByID(ctx context.Context, id string) (*model.Product, error)
}

type ProductRepo struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepo {
	_ = db.AutoMigrate(&model.Product{})
	return &ProductRepo{db: db}
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
