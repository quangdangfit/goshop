package repositories

import (
	"context"

	"gorm.io/gorm"

	"goshop/app/models"
	"goshop/config"
	"goshop/dbs"
)

type IOrderLineRepository interface {
	CreateOrderLines(ctx context.Context, req []*models.OrderLine) error
}

type OrderLineRepo struct {
	db *gorm.DB
}

func NewOrderLineRepository() *OrderLineRepo {
	return &OrderLineRepo{db: dbs.Database}
}

func (r *OrderLineRepo) CreateOrderLines(ctx context.Context, req []*models.OrderLine) error {
	ctx, cancel := context.WithTimeout(ctx, config.DatabaseTimeout)
	defer cancel()

	if err := r.db.CreateInBatches(&req, len(req)).Error; err != nil {
		return err
	}

	return nil
}
