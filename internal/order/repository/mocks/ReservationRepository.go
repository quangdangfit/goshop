// Hand-written mock following the mockery testify template style; mockery v3 cannot
// re-run against the go1.26 source tree, so this file is maintained manually.

package mocks

import (
	"context"
	"time"

	mock "github.com/stretchr/testify/mock"

	"goshop/internal/order/model"
)

func NewReservationRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *ReservationRepository {
	m := &ReservationRepository{}
	m.Test(t)
	t.Cleanup(func() { m.AssertExpectations(t) })
	return m
}

type ReservationRepository struct {
	mock.Mock
}

func (_m *ReservationRepository) CreateMany(ctx context.Context, items []*model.StockReservation) error {
	ret := _m.Called(ctx, items)
	if fn, ok := ret.Get(0).(func(context.Context, []*model.StockReservation) error); ok {
		return fn(ctx, items)
	}
	return ret.Error(0)
}

func (_m *ReservationRepository) FindActiveByOrderID(ctx context.Context, orderID string) ([]*model.StockReservation, error) {
	ret := _m.Called(ctx, orderID)
	var r0 []*model.StockReservation
	if fn, ok := ret.Get(0).(func(context.Context, string) ([]*model.StockReservation, error)); ok {
		return fn(ctx, orderID)
	}
	if ret.Get(0) != nil {
		r0 = ret.Get(0).([]*model.StockReservation)
	}
	return r0, ret.Error(1)
}

func (_m *ReservationRepository) FindExpired(ctx context.Context, now time.Time, limit int) ([]*model.StockReservation, error) {
	ret := _m.Called(ctx, now, limit)
	var r0 []*model.StockReservation
	if fn, ok := ret.Get(0).(func(context.Context, time.Time, int) ([]*model.StockReservation, error)); ok {
		return fn(ctx, now, limit)
	}
	if ret.Get(0) != nil {
		r0 = ret.Get(0).([]*model.StockReservation)
	}
	return r0, ret.Error(1)
}

func (_m *ReservationRepository) UpdateStatus(ctx context.Context, ids []string, status model.ReservationStatus) error {
	ret := _m.Called(ctx, ids, status)
	if fn, ok := ret.Get(0).(func(context.Context, []string, model.ReservationStatus) error); ok {
		return fn(ctx, ids, status)
	}
	return ret.Error(0)
}
