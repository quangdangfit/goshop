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
	// ReserveStock atomically increments reserved_quantity by qty if the available stock
	// (stock_quantity - reserved_quantity) is at least qty. Returns ErrInsufficientStock
	// when no row matches.
	ReserveStock(ctx context.Context, id string, qty int) error
	// CommitReservation atomically decrements both stock_quantity and reserved_quantity by
	// qty, finalizing a previously reserved sale.
	CommitReservation(ctx context.Context, id string, qty int) error
	// ReleaseReservation atomically decrements reserved_quantity by qty, returning units to
	// available stock without touching the on-hand total.
	ReleaseReservation(ctx context.Context, id string, qty int) error
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

// DecrementStock atomically subtracts qty from stock_quantity, but only if the available
// stock (stock_quantity - reserved_quantity) is at least qty. Returns ErrInsufficientStock
// when no row matches.
func (r *productRepo) DecrementStock(ctx context.Context, id string, qty int) error {
	result := r.db.GetDB().WithContext(ctx).Model(&model.Product{}).
		Where("id = ? AND stock_quantity - reserved_quantity >= ?", id, qty).
		UpdateColumn("stock_quantity", gorm.Expr("stock_quantity - ?", qty))
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrInsufficientStock
	}
	return nil
}

// ErrInsufficientStock is returned when a stock-modifying operation cannot proceed because
// available stock is below the requested quantity.
var ErrInsufficientStock = errors.New("insufficient stock")

func (r *productRepo) ReserveStock(ctx context.Context, id string, qty int) error {
	result := r.db.GetDB().WithContext(ctx).Model(&model.Product{}).
		Where("id = ? AND stock_quantity - reserved_quantity >= ?", id, qty).
		UpdateColumn("reserved_quantity", gorm.Expr("reserved_quantity + ?", qty))
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrInsufficientStock
	}
	return nil
}

func (r *productRepo) CommitReservation(ctx context.Context, id string, qty int) error {
	result := r.db.GetDB().WithContext(ctx).Model(&model.Product{}).
		Where("id = ? AND reserved_quantity >= ? AND stock_quantity >= ?", id, qty, qty).
		Updates(map[string]any{
			"stock_quantity":    gorm.Expr("stock_quantity - ?", qty),
			"reserved_quantity": gorm.Expr("reserved_quantity - ?", qty),
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("commit reservation: row not in expected state")
	}
	return nil
}

func (r *productRepo) ReleaseReservation(ctx context.Context, id string, qty int) error {
	result := r.db.GetDB().WithContext(ctx).Model(&model.Product{}).
		Where("id = ? AND reserved_quantity >= ?", id, qty).
		UpdateColumn("reserved_quantity", gorm.Expr("reserved_quantity - ?", qty))
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("release reservation: row not in expected state")
	}
	return nil
}
