package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"goshop/internal/payment/model"
	"goshop/pkg/dbs"
)

var ErrEventAlreadyProcessed = errors.New("provider event already processed")

//go:generate mockery --name=PaymentRepository
type PaymentRepository interface {
	GetByOrderID(ctx context.Context, orderID string) (*model.Payment, error)
	Create(ctx context.Context, p *model.Payment) error
	Update(ctx context.Context, p *model.Payment) error
	// RecordProviderEvent inserts a (provider, event_id) row and returns ErrEventAlreadyProcessed
	// if the event has been seen before. Used to make webhook handling idempotent.
	RecordProviderEvent(ctx context.Context, provider, eventID string) error
}

type paymentRepo struct {
	db dbs.Database
}

func NewPaymentRepository(db dbs.Database) PaymentRepository {
	return &paymentRepo{db: db}
}

func (r *paymentRepo) GetByOrderID(ctx context.Context, orderID string) (*model.Payment, error) {
	var p model.Payment
	err := r.db.GetDB().WithContext(ctx).Where("order_id = ?", orderID).First(&p).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *paymentRepo) Create(ctx context.Context, p *model.Payment) error {
	return r.db.Create(ctx, p)
}

func (r *paymentRepo) Update(ctx context.Context, p *model.Payment) error {
	return r.db.Update(ctx, p)
}

func (r *paymentRepo) RecordProviderEvent(ctx context.Context, provider, eventID string) error {
	err := r.db.GetDB().WithContext(ctx).Create(&model.ProviderEvent{Provider: provider, EventID: eventID}).Error
	if err == nil {
		return nil
	}
	// Postgres unique-violation surfaces as duplicate key error; treat as already-processed.
	if isDuplicateKey(err) {
		return ErrEventAlreadyProcessed
	}
	return err
}

func isDuplicateKey(err error) bool {
	if err == nil {
		return false
	}
	// gorm wraps the driver error; pq returns "duplicate key value violates unique constraint".
	msg := err.Error()
	return errors.Is(err, gorm.ErrDuplicatedKey) ||
		containsAny(msg, "duplicate key", "UNIQUE constraint", "Duplicate entry")
}

func containsAny(s string, subs ...string) bool {
	for _, sub := range subs {
		if len(sub) > 0 && len(s) >= len(sub) && indexOf(s, sub) >= 0 {
			return true
		}
	}
	return false
}

func indexOf(s, sub string) int {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}
