package service

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/quangdangfit/gocommon/logger"
	"github.com/quangdangfit/gocommon/validation"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"goshop/internal/order/model"
	orderMocks "goshop/internal/order/repository/mocks"
	serviceMocks "goshop/internal/order/service/mocks"
	"goshop/pkg/config"
	dbsMocks "goshop/pkg/dbs/mocks"
	"goshop/pkg/eventbus"
	notifMocks "goshop/pkg/notification/mocks"
)

type markPaidFixture struct {
	svc          *orderService
	db           *dbsMocks.Database
	repo         *orderMocks.OrderRepository
	productRepo  *orderMocks.ProductRepository
	reservRepo   *orderMocks.ReservationRepository
	bus          eventbus.Bus
	captured     []eventbus.LowStock
	capturedLock sync.Mutex
}

func newMarkPaidFixture(t *testing.T) *markPaidFixture {
	t.Helper()
	logger.Initialize(config.ProductionEnv)
	db := dbsMocks.NewDatabase(t)
	repo := orderMocks.NewOrderRepository(t)
	productRepo := orderMocks.NewProductRepository(t)
	reservRepo := orderMocks.NewReservationRepository(t)
	userRepo := orderMocks.NewUserRepository(t)
	couponSvc := serviceMocks.NewCouponService(t)
	notifier := notifMocks.NewNotifier(t)
	db.On("WithTransaction", mock.Anything).Return(func(fn func() error) error { return fn() }).Maybe()

	svc := NewOrderService(validation.New(), db, repo, productRepo, userRepo, reservRepo, couponSvc, notifier).(*orderService)
	bus := eventbus.New()
	svc.SetEventBus(bus)

	f := &markPaidFixture{
		svc: svc, db: db, repo: repo, productRepo: productRepo, reservRepo: reservRepo, bus: bus,
	}
	bus.Subscribe(eventbus.TopicLowStock, func(_ context.Context, ev eventbus.Event) {
		f.capturedLock.Lock()
		defer f.capturedLock.Unlock()
		f.captured = append(f.captured, ev.(eventbus.LowStock))
	})
	return f
}

func (f *markPaidFixture) waitForEvents(t *testing.T, n int) {
	t.Helper()
	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) {
		f.capturedLock.Lock()
		got := len(f.captured)
		f.capturedLock.Unlock()
		if got >= n {
			return
		}
		time.Sleep(5 * time.Millisecond)
	}
	t.Fatalf("expected %d LowStock events, got %d", n, len(f.captured))
}

// TestMarkOrderPaid_LowStockEvent table-drives the post-commit LowStock-event
// behavior: the happy commit path is identical across cases; only the final
// GetProductByID outcome (stock value or error) varies.
func TestMarkOrderPaid_LowStockEvent(t *testing.T) {
	tests := []struct {
		name             string
		productResult    *model.Product
		productErr       error
		wantEvents       int
		wantAvailable    int // only checked when wantEvents > 0
		wantOrderSuccess bool
	}{
		{
			name:             "publishes_event_when_below_threshold",
			productResult:    &model.Product{ID: "p1", StockQuantity: 4, ReservedQuantity: 1},
			wantEvents:       1,
			wantAvailable:    3,
			wantOrderSuccess: true,
		},
		{
			name:             "no_event_when_above_threshold",
			productResult:    &model.Product{ID: "p1", StockQuantity: 100, ReservedQuantity: 1},
			wantEvents:       0,
			wantOrderSuccess: true,
		},
		{
			name:             "lookup_error_is_non_fatal",
			productErr:       errors.New("db down"),
			wantEvents:       0,
			wantOrderSuccess: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := newMarkPaidFixture(t)
			order := &model.Order{ID: "o1", Status: model.OrderStatusPendingPayment}
			reservations := []*model.StockReservation{{ID: "r1", OrderID: "o1", ProductID: "p1", Quantity: 1}}

			f.repo.On("GetOrderByID", mock.Anything, "o1", true).Return(order, nil).Once()
			f.reservRepo.On("FindActiveByOrderID", mock.Anything, "o1").Return(reservations, nil).Once()
			f.productRepo.On("CommitReservation", mock.Anything, "p1", 1).Return(nil).Once()
			f.reservRepo.On("UpdateStatus", mock.Anything, []string{"r1"}, model.ReservationStatusCommitted).Return(nil).Once()
			f.repo.On("UpdateOrder", mock.Anything, mock.Anything).Return(nil).Once()
			f.productRepo.On("GetProductByID", mock.Anything, "p1").Return(tt.productResult, tt.productErr).Once()

			got, err := f.svc.MarkOrderPaid(context.Background(), "o1")
			if tt.wantOrderSuccess {
				require.NoError(t, err)
				require.Equal(t, model.OrderStatusPaid, got.Status)
			} else {
				require.Error(t, err)
			}

			if tt.wantEvents > 0 {
				f.waitForEvents(t, tt.wantEvents)
				require.Equal(t, "p1", f.captured[0].ProductID)
				require.Equal(t, tt.wantAvailable, f.captured[0].Available)
				require.Equal(t, LowStockThreshold, f.captured[0].Threshold)
			} else {
				time.Sleep(50 * time.Millisecond)
				f.capturedLock.Lock()
				defer f.capturedLock.Unlock()
				require.Empty(t, f.captured)
			}
		})
	}
}

func TestMarkOrderPaid_IdempotentOnAlreadyPaid(t *testing.T) {
	f := newMarkPaidFixture(t)
	order := &model.Order{ID: "o1", Status: model.OrderStatusPaid}
	f.repo.On("GetOrderByID", mock.Anything, "o1", true).Return(order, nil).Once()

	got, err := f.svc.MarkOrderPaid(context.Background(), "o1")
	require.NoError(t, err)
	require.Equal(t, model.OrderStatusPaid, got.Status)
}

func TestMarkOrderPaid_RejectsCancelledOrFailed(t *testing.T) {
	tests := []struct {
		name   string
		status model.OrderStatus
	}{
		{"cancelled", model.OrderStatusCancelled},
		{"payment_failed", model.OrderStatusPaymentFailed},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := newMarkPaidFixture(t)
			order := &model.Order{ID: "o1", Status: tt.status}
			f.repo.On("GetOrderByID", mock.Anything, "o1", true).Return(order, nil).Once()

			_, err := f.svc.MarkOrderPaid(context.Background(), "o1")
			require.Error(t, err)
		})
	}
}

func TestMarkOrderPaid_GetOrderError(t *testing.T) {
	f := newMarkPaidFixture(t)
	f.repo.On("GetOrderByID", mock.Anything, "o1", true).Return(nil, errors.New("not found")).Once()

	_, err := f.svc.MarkOrderPaid(context.Background(), "o1")
	require.Error(t, err)
}

func TestMarkOrderPaid_CommitReservationFails(t *testing.T) {
	f := newMarkPaidFixture(t)
	order := &model.Order{ID: "o1", Status: model.OrderStatusPendingPayment}
	reservations := []*model.StockReservation{{ID: "r1", OrderID: "o1", ProductID: "p1", Quantity: 1}}

	f.repo.On("GetOrderByID", mock.Anything, "o1", true).Return(order, nil).Once()
	f.reservRepo.On("FindActiveByOrderID", mock.Anything, "o1").Return(reservations, nil).Once()
	f.productRepo.On("CommitReservation", mock.Anything, "p1", 1).Return(errors.New("constraint violation")).Once()

	_, err := f.svc.MarkOrderPaid(context.Background(), "o1")
	require.Error(t, err)
}

func TestMarkOrderPaid_FallsBackToDefaultBus(t *testing.T) {
	logger.Initialize(config.ProductionEnv)
	defaultBus := eventbus.New()
	eventbus.SetDefault(defaultBus)
	t.Cleanup(func() { eventbus.SetDefault(eventbus.New()) })

	var captured []eventbus.LowStock
	var mu sync.Mutex
	defaultBus.Subscribe(eventbus.TopicLowStock, func(_ context.Context, ev eventbus.Event) {
		mu.Lock()
		defer mu.Unlock()
		captured = append(captured, ev.(eventbus.LowStock))
	})

	db := dbsMocks.NewDatabase(t)
	repo := orderMocks.NewOrderRepository(t)
	productRepo := orderMocks.NewProductRepository(t)
	reservRepo := orderMocks.NewReservationRepository(t)
	userRepo := orderMocks.NewUserRepository(t)
	couponSvc := serviceMocks.NewCouponService(t)
	notifier := notifMocks.NewNotifier(t)
	db.On("WithTransaction", mock.Anything).Return(func(fn func() error) error { return fn() }).Maybe()
	svc := NewOrderService(validation.New(), db, repo, productRepo, userRepo, reservRepo, couponSvc, notifier)

	order := &model.Order{ID: "o1", Status: model.OrderStatusPendingPayment}
	reservations := []*model.StockReservation{{ID: "r1", OrderID: "o1", ProductID: "p1", Quantity: 1}}
	repo.On("GetOrderByID", mock.Anything, "o1", true).Return(order, nil).Once()
	reservRepo.On("FindActiveByOrderID", mock.Anything, "o1").Return(reservations, nil).Once()
	productRepo.On("CommitReservation", mock.Anything, "p1", 1).Return(nil).Once()
	reservRepo.On("UpdateStatus", mock.Anything, []string{"r1"}, model.ReservationStatusCommitted).Return(nil).Once()
	repo.On("UpdateOrder", mock.Anything, mock.Anything).Return(nil).Once()
	productRepo.On("GetProductByID", mock.Anything, "p1").
		Return(&model.Product{ID: "p1", StockQuantity: 2, ReservedQuantity: 1}, nil).Once()

	_, err := svc.MarkOrderPaid(context.Background(), "o1")
	require.NoError(t, err)

	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) {
		mu.Lock()
		if len(captured) > 0 {
			mu.Unlock()
			return
		}
		mu.Unlock()
		time.Sleep(5 * time.Millisecond)
	}
	t.Fatal("expected event on default bus")
}
