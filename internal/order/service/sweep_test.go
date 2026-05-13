package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/quangdangfit/gocommon/logger"
	"github.com/quangdangfit/gocommon/validation"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"goshop/internal/order/model"
	orderRepo "goshop/internal/order/repository"
	orderMocks "goshop/internal/order/repository/mocks"
	serviceMocks "goshop/internal/order/service/mocks"
	"goshop/pkg/config"
	dbsMocks "goshop/pkg/dbs/mocks"
	notifMocks "goshop/pkg/notification/mocks"
)

type sweepFixture struct {
	svc         OrderService
	db          *dbsMocks.Database
	repo        *orderMocks.OrderRepository
	productRepo *orderMocks.ProductRepository
	userRepo    *orderMocks.UserRepository
	reservRepo  *orderMocks.ReservationRepository
	notifier    *notifMocks.Notifier
}

func newSweepFixture(t *testing.T) *sweepFixture {
	logger.Initialize(config.ProductionEnv)
	db := dbsMocks.NewDatabase(t)
	repo := orderMocks.NewOrderRepository(t)
	productRepo := orderMocks.NewProductRepository(t)
	userRepo := orderMocks.NewUserRepository(t)
	reservRepo := orderMocks.NewReservationRepository(t)
	couponSvc := serviceMocks.NewCouponService(t)
	notifier := notifMocks.NewNotifier(t)
	db.On("WithTransaction", mock.Anything).Return(func(fn func() error) error { return fn() }).Maybe()
	svc := NewOrderService(validation.New(), db, repo, productRepo, userRepo, reservRepo, couponSvc, notifier)
	return &sweepFixture{svc, db, repo, productRepo, userRepo, reservRepo, notifier}
}

func TestSweep_EmptyBatchReturnsZero(t *testing.T) {
	f := newSweepFixture(t)
	f.reservRepo.On("FindExpired", mock.Anything, mock.Anything, 100).Return(nil, nil).Once()
	n, err := f.svc.SweepExpiredReservations(context.Background(), 0)
	require.NoError(t, err)
	require.Zero(t, n)
}

func TestSweep_FindExpiredError(t *testing.T) {
	f := newSweepFixture(t)
	f.reservRepo.On("FindExpired", mock.Anything, mock.Anything, 50).Return(nil, errors.New("db")).Once()
	_, err := f.svc.SweepExpiredReservations(context.Background(), 50)
	require.Error(t, err)
}

func TestSweep_HappyPathReleasesAndCancels(t *testing.T) {
	f := newSweepFixture(t)
	expired := []*model.StockReservation{
		{ID: "r1", OrderID: "o1", ProductID: "p1", Quantity: 1},
	}
	f.reservRepo.On("FindExpired", mock.Anything, mock.Anything, 100).Return(expired, nil).Once()
	f.repo.On("GetOrderByID", mock.Anything, "o1", false).
		Return(&model.Order{ID: "o1", Status: model.OrderStatusPendingPayment, UserID: "u1"}, nil)
	f.productRepo.On("ReleaseReservation", mock.Anything, "p1", 1).Return(nil).Once()
	f.reservRepo.On("UpdateStatus", mock.Anything, []string{"r1"}, model.ReservationStatusReleased).Return(nil).Once()
	f.repo.On("UpdateOrder", mock.Anything, mock.Anything).Return(nil).Once()
	// Notify goroutine
	f.userRepo.On("GetUserByID", mock.Anything, "u1").Return(nil, nil).Maybe()

	n, err := f.svc.SweepExpiredReservations(context.Background(), 100)
	require.NoError(t, err)
	require.Equal(t, 1, n)
	time.Sleep(20 * time.Millisecond)
}

func TestSweep_SkipsNonPendingOrder(t *testing.T) {
	f := newSweepFixture(t)
	expired := []*model.StockReservation{
		{ID: "r1", OrderID: "o1", ProductID: "p1", Quantity: 1},
	}
	f.reservRepo.On("FindExpired", mock.Anything, mock.Anything, 100).Return(expired, nil).Once()
	f.repo.On("GetOrderByID", mock.Anything, "o1", false).
		Return(&model.Order{ID: "o1", Status: model.OrderStatusPaid, UserID: "u1"}, nil)
	f.userRepo.On("GetUserByID", mock.Anything, "u1").Return(nil, nil).Maybe()

	n, err := f.svc.SweepExpiredReservations(context.Background(), 100)
	require.NoError(t, err)
	require.Equal(t, 1, n)
	time.Sleep(20 * time.Millisecond)
}

func TestSweep_TxErrorIsLoggedNotReturned(t *testing.T) {
	f := newSweepFixture(t)
	expired := []*model.StockReservation{
		{ID: "r1", OrderID: "o1", ProductID: "p1", Quantity: 1},
	}
	f.reservRepo.On("FindExpired", mock.Anything, mock.Anything, 100).Return(expired, nil).Once()
	f.repo.On("GetOrderByID", mock.Anything, "o1", false).
		Return(&model.Order{ID: "o1", Status: model.OrderStatusPendingPayment}, nil)
	f.productRepo.On("ReleaseReservation", mock.Anything, "p1", 1).Return(errors.New("locked")).Once()

	n, err := f.svc.SweepExpiredReservations(context.Background(), 100)
	require.NoError(t, err)
	require.Zero(t, n)
}

func TestSweep_OrphanedReservation_StillReleasesAndDoesNotError(t *testing.T) {
	f := newSweepFixture(t)
	expired := []*model.StockReservation{
		{ID: "r1", OrderID: "ghost", ProductID: "p1", Quantity: 2},
	}
	f.reservRepo.On("FindExpired", mock.Anything, mock.Anything, 100).Return(expired, nil).Once()
	// Order row is missing — GetOrderByID returns ErrRecordNotFound.
	f.repo.On("GetOrderByID", mock.Anything, "ghost", false).Return(nil, gorm.ErrRecordNotFound).Once()
	// Sweeper should still release the held stock and mark the reservation released.
	f.productRepo.On("ReleaseReservation", mock.Anything, "p1", 2).Return(nil).Once()
	f.reservRepo.On("UpdateStatus", mock.Anything, []string{"r1"}, model.ReservationStatusReleased).Return(nil).Once()
	// And NOT call UpdateOrder, since there is no order row to update.

	n, err := f.svc.SweepExpiredReservations(context.Background(), 100)
	require.NoError(t, err)
	require.Equal(t, 1, n)
	time.Sleep(20 * time.Millisecond)
}

func TestSweep_ReservationAlreadyReleased_IsTolerated(t *testing.T) {
	f := newSweepFixture(t)
	expired := []*model.StockReservation{
		{ID: "r1", OrderID: "o1", ProductID: "p1", Quantity: 1},
	}
	f.reservRepo.On("FindExpired", mock.Anything, mock.Anything, 100).Return(expired, nil).Once()
	f.repo.On("GetOrderByID", mock.Anything, "o1", false).
		Return(&model.Order{ID: "o1", Status: model.OrderStatusPendingPayment, UserID: "u1"}, nil)
	// The product's counter is already drained (drift). Sweeper should still
	// mark the reservation released and cancel the order.
	f.productRepo.On("ReleaseReservation", mock.Anything, "p1", 1).
		Return(orderRepo.ErrReservationAlreadyReleased).Once()
	f.reservRepo.On("UpdateStatus", mock.Anything, []string{"r1"}, model.ReservationStatusReleased).Return(nil).Once()
	f.repo.On("UpdateOrder", mock.Anything, mock.Anything).Return(nil).Once()
	f.userRepo.On("GetUserByID", mock.Anything, "u1").Return(nil, nil).Maybe()

	n, err := f.svc.SweepExpiredReservations(context.Background(), 100)
	require.NoError(t, err)
	require.Equal(t, 1, n)
	time.Sleep(20 * time.Millisecond)
}

func TestInsufficientStockErrorMessage(t *testing.T) {
	e := &InsufficientStockError{ProductID: "p1", Requested: 5}
	require.Contains(t, e.Error(), "p1")
	require.Contains(t, e.Error(), "5")
}
