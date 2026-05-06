package repository

import (
	"context"
	"time"

	"goshop/internal/order/model"
	"goshop/pkg/dbs"
)

//go:generate mockery --name=ReservationRepository
type ReservationRepository interface {
	CreateMany(ctx context.Context, items []*model.StockReservation) error
	FindActiveByOrderID(ctx context.Context, orderID string) ([]*model.StockReservation, error)
	FindExpired(ctx context.Context, now time.Time, limit int) ([]*model.StockReservation, error)
	UpdateStatus(ctx context.Context, ids []string, status model.ReservationStatus) error
}

type reservationRepo struct {
	db dbs.Database
}

func NewReservationRepository(db dbs.Database) ReservationRepository {
	return &reservationRepo{db: db}
}

func (r *reservationRepo) CreateMany(ctx context.Context, items []*model.StockReservation) error {
	if len(items) == 0 {
		return nil
	}
	return r.db.CreateInBatches(ctx, &items, len(items))
}

func (r *reservationRepo) FindActiveByOrderID(ctx context.Context, orderID string) ([]*model.StockReservation, error) {
	var rows []*model.StockReservation
	err := r.db.GetDB().WithContext(ctx).
		Where("order_id = ? AND status = ?", orderID, model.ReservationStatusActive).
		Find(&rows).Error
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (r *reservationRepo) FindExpired(ctx context.Context, now time.Time, limit int) ([]*model.StockReservation, error) {
	var rows []*model.StockReservation
	err := r.db.GetDB().WithContext(ctx).
		Where("status = ? AND expires_at < ?", model.ReservationStatusActive, now).
		Limit(limit).
		Find(&rows).Error
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (r *reservationRepo) UpdateStatus(ctx context.Context, ids []string, status model.ReservationStatus) error {
	if len(ids) == 0 {
		return nil
	}
	return r.db.GetDB().WithContext(ctx).
		Model(&model.StockReservation{}).
		Where("id IN ?", ids).
		Update("status", status).Error
}
