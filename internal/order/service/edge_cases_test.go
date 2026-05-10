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

	"goshop/internal/order/domain"
	"goshop/internal/order/model"
	orderMocks "goshop/internal/order/repository/mocks"
	serviceMocks "goshop/internal/order/service/mocks"
	"goshop/pkg/config"
	dbsMocks "goshop/pkg/dbs/mocks"
	notifMocks "goshop/pkg/notification/mocks"
)

func newEdgeFixture(t *testing.T) (OrderService, *dbsMocks.Database, *orderMocks.OrderRepository, *orderMocks.ProductRepository, *orderMocks.UserRepository, *orderMocks.ReservationRepository, *notifMocks.Notifier) {
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
	return svc, db, repo, productRepo, userRepo, reservRepo, notifier
}

func TestPlaceOrder_ReserveStockReturnsInsufficient(t *testing.T) {
	svc, _, repo, productRepo, _, _, _ := newEdgeFixture(t)
	productRepo.On("GetProductByID", mock.Anything, "p1").Return(&model.Product{Name: "p", Price: 1}, nil).Once()
	repo.On("CreateOrder", mock.Anything, "u1", mock.Anything, "", float64(0)).
		Return(&model.Order{ID: "o1", UserID: "u1", Lines: []*model.OrderLine{{ProductID: "p1", Quantity: 2}}}, nil).Once()
	repo.On("UpdateOrder", mock.Anything, mock.Anything).Return(nil).Once()
	productRepo.On("ReserveStock", mock.Anything, "p1", 2).
		Return(&InsufficientStockError{ProductID: "p1", Requested: 2}).Once()

	_, err := svc.PlaceOrder(context.Background(), &domain.PlaceOrderReq{
		UserID: "u1",
		Lines:  []domain.PlaceOrderLineReq{{ProductID: "p1", Quantity: 2}},
	})
	require.Error(t, err)
}

func TestPlaceOrder_CreateOrderError(t *testing.T) {
	svc, _, repo, productRepo, _, _, _ := newEdgeFixture(t)
	productRepo.On("GetProductByID", mock.Anything, "p1").Return(&model.Product{Name: "p", Price: 1}, nil).Once()
	repo.On("CreateOrder", mock.Anything, "u1", mock.Anything, "", float64(0)).Return(nil, errors.New("db")).Once()

	_, err := svc.PlaceOrder(context.Background(), &domain.PlaceOrderReq{
		UserID: "u1",
		Lines:  []domain.PlaceOrderLineReq{{ProductID: "p1", Quantity: 1}},
	})
	require.Error(t, err)
}

func TestPlaceOrder_CreateManyError(t *testing.T) {
	svc, _, repo, productRepo, _, reservRepo, _ := newEdgeFixture(t)
	productRepo.On("GetProductByID", mock.Anything, "p1").Return(&model.Product{Name: "p", Price: 1}, nil).Once()
	repo.On("CreateOrder", mock.Anything, "u1", mock.Anything, "", float64(0)).
		Return(&model.Order{ID: "o1", UserID: "u1", Lines: []*model.OrderLine{{ProductID: "p1", Quantity: 1}}}, nil).Once()
	repo.On("UpdateOrder", mock.Anything, mock.Anything).Return(nil).Once()
	productRepo.On("ReserveStock", mock.Anything, "p1", 1).Return(nil).Once()
	reservRepo.On("CreateMany", mock.Anything, mock.Anything).Return(errors.New("disk full")).Once()

	_, err := svc.PlaceOrder(context.Background(), &domain.PlaceOrderReq{
		UserID: "u1",
		Lines:  []domain.PlaceOrderLineReq{{ProductID: "p1", Quantity: 1}},
	})
	require.Error(t, err)
}

func TestPlaceOrder_NotifyHappy(t *testing.T) {
	svc, _, repo, productRepo, userRepo, reservRepo, notifier := newEdgeFixture(t)
	productRepo.On("GetProductByID", mock.Anything, "p1").Return(&model.Product{Name: "p", Price: 1}, nil).Once()
	repo.On("CreateOrder", mock.Anything, "u1", mock.Anything, "", float64(0)).
		Return(&model.Order{ID: "o1", UserID: "u1", Lines: []*model.OrderLine{{ProductID: "p1", Quantity: 1}}}, nil).Once()
	repo.On("UpdateOrder", mock.Anything, mock.Anything).Return(nil).Once()
	productRepo.On("ReserveStock", mock.Anything, "p1", 1).Return(nil).Once()
	reservRepo.On("CreateMany", mock.Anything, mock.Anything).Return(nil).Once()
	userRepo.On("GetUserByID", mock.Anything, "u1").Return(&model.User{ID: "u1", Email: "x@example.com"}, nil).Maybe()
	notifier.On("SendOrderPlaced", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("smtp")).Maybe()

	_, err := svc.PlaceOrder(context.Background(), &domain.PlaceOrderReq{
		UserID: "u1",
		Lines:  []domain.PlaceOrderLineReq{{ProductID: "p1", Quantity: 1}},
	})
	require.NoError(t, err)
	time.Sleep(20 * time.Millisecond)
}

func TestUpdateOrderStatus_GetOrderError(t *testing.T) {
	svc, _, repo, _, _, _, _ := newEdgeFixture(t)
	repo.On("GetOrderByID", mock.Anything, "o1", false).Return(nil, errors.New("not found")).Once()
	_, err := svc.UpdateOrderStatus(context.Background(), "o1", model.OrderStatusPaid)
	require.Error(t, err)
}

func TestUpdateOrderStatus_InvalidStatus(t *testing.T) {
	svc, _, repo, _, _, _, _ := newEdgeFixture(t)
	repo.On("GetOrderByID", mock.Anything, "o1", false).
		Return(&model.Order{ID: "o1", Status: model.OrderStatusNew}, nil).Once()
	_, err := svc.UpdateOrderStatus(context.Background(), "o1", model.OrderStatus("bogus"))
	require.Error(t, err)
}

func TestUpdateOrderStatus_DisallowedTransition(t *testing.T) {
	svc, _, repo, _, _, _, _ := newEdgeFixture(t)
	repo.On("GetOrderByID", mock.Anything, "o1", false).
		Return(&model.Order{ID: "o1", Status: model.OrderStatusDone}, nil).Once()
	_, err := svc.UpdateOrderStatus(context.Background(), "o1", model.OrderStatusCancelled)
	require.Error(t, err)
}

func TestUpdateOrderStatus_UpdateError(t *testing.T) {
	svc, _, repo, _, _, _, _ := newEdgeFixture(t)
	repo.On("GetOrderByID", mock.Anything, "o1", false).
		Return(&model.Order{ID: "o1", Status: model.OrderStatusPaid}, nil).Once()
	repo.On("UpdateOrder", mock.Anything, mock.Anything).Return(errors.New("db")).Once()
	_, err := svc.UpdateOrderStatus(context.Background(), "o1", model.OrderStatusInProgress)
	require.Error(t, err)
}

func TestCancelOrder_GetError(t *testing.T) {
	svc, _, repo, _, _, _, _ := newEdgeFixture(t)
	repo.On("GetOrderByID", mock.Anything, "o1", false).Return(nil, errors.New("not found")).Once()
	_, err := svc.CancelOrder(context.Background(), "o1", "u1")
	require.Error(t, err)
}

func TestCancelOrder_Forbidden(t *testing.T) {
	svc, _, repo, _, _, _, _ := newEdgeFixture(t)
	repo.On("GetOrderByID", mock.Anything, "o1", false).
		Return(&model.Order{ID: "o1", UserID: "owner"}, nil).Once()
	_, err := svc.CancelOrder(context.Background(), "o1", "intruder")
	require.Error(t, err)
}

func TestCancelOrder_TerminalStatus(t *testing.T) {
	for _, st := range []model.OrderStatus{model.OrderStatusDone, model.OrderStatusCancelled, model.OrderStatusPaid} {
		svc, _, repo, _, _, _, _ := newEdgeFixture(t)
		repo.On("GetOrderByID", mock.Anything, "o1", false).
			Return(&model.Order{ID: "o1", UserID: "u1", Status: st}, nil).Once()
		_, err := svc.CancelOrder(context.Background(), "o1", "u1")
		require.Error(t, err)
	}
}

func TestCancelOrder_HappyPath(t *testing.T) {
	svc, _, repo, productRepo, _, reservRepo, _ := newEdgeFixture(t)
	repo.On("GetOrderByID", mock.Anything, "o1", false).
		Return(&model.Order{ID: "o1", UserID: "u1", Status: model.OrderStatusPendingPayment}, nil).Once()
	reservRepo.On("FindActiveByOrderID", mock.Anything, "o1").
		Return([]*model.StockReservation{{ID: "r1", ProductID: "p1", Quantity: 1}}, nil).Once()
	productRepo.On("ReleaseReservation", mock.Anything, "p1", 1).Return(nil).Once()
	reservRepo.On("UpdateStatus", mock.Anything, []string{"r1"}, model.ReservationStatusReleased).Return(nil).Once()
	repo.On("UpdateOrder", mock.Anything, mock.Anything).Return(nil).Once()

	_, err := svc.CancelOrder(context.Background(), "o1", "u1")
	require.NoError(t, err)
}

func TestCancelOrder_ReleaseError(t *testing.T) {
	svc, _, repo, productRepo, _, reservRepo, _ := newEdgeFixture(t)
	repo.On("GetOrderByID", mock.Anything, "o1", false).
		Return(&model.Order{ID: "o1", UserID: "u1", Status: model.OrderStatusPendingPayment}, nil).Once()
	reservRepo.On("FindActiveByOrderID", mock.Anything, "o1").
		Return([]*model.StockReservation{{ID: "r1", ProductID: "p1", Quantity: 1}}, nil).Once()
	productRepo.On("ReleaseReservation", mock.Anything, "p1", 1).Return(errors.New("locked")).Once()

	_, err := svc.CancelOrder(context.Background(), "o1", "u1")
	require.Error(t, err)
}

func TestCancelOrder_FindReservationsError(t *testing.T) {
	svc, _, repo, _, _, reservRepo, _ := newEdgeFixture(t)
	repo.On("GetOrderByID", mock.Anything, "o1", false).
		Return(&model.Order{ID: "o1", UserID: "u1", Status: model.OrderStatusPendingPayment}, nil).Once()
	reservRepo.On("FindActiveByOrderID", mock.Anything, "o1").Return(nil, errors.New("db")).Once()

	_, err := svc.CancelOrder(context.Background(), "o1", "u1")
	require.Error(t, err)
}
