package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"goshop/internal/order/model"
	"goshop/pkg/dbs"
)

//go:generate mockery --name=ProductRepository
type ProductRepository interface {
	GetProductByID(ctx context.Context, id string) (*model.Product, error)
	DecrementStock(ctx context.Context, id string, qty int) error
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

func (r *productRepo) DecrementStock(ctx context.Context, id string, qty int) error {
	result := r.db.GetDB().WithContext(ctx).Model(&model.Product{}).
		Where("id = ? AND stock_quantity >= ?", id, qty).
		UpdateColumn("stock_quantity", gorm.Expr("stock_quantity - ?", qty))
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("insufficient stock")
	}
	return nil
}
