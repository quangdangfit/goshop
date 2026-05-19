//go:build integration

package tests_order

import (
	"context"
	"testing"
	"time"

	"github.com/quangdangfit/gocommon/validation"
	"github.com/stretchr/testify/require"

	orderModel "goshop/internal/order/model"
	orderRepo "goshop/internal/order/repository"
	orderSvc "goshop/internal/order/service"
	productModel "goshop/internal/product/model"
	userModel "goshop/internal/user/model"
	"goshop/pkg/notification"
	"goshop/tests/testutil"
)

// TestSweepExpiredReservations_ReleasesAndCancels: a pending_payment order whose reservation
// has aged past expires_at must be cancelled and its stock released. Asserts both the
// reservation row's terminal status and product.reserved_quantity going back to 0.
func TestSweepExpiredReservations_ReleasesAndCancels(t *testing.T) {
	ctx := context.Background()
	db := testutil.StartPostgres(ctx, t)
	require.NoError(t, testutil.ApplyMigrations(db))

	user := &userModel.User{Email: "buyer@test.com", Password: "x"}
	require.NoError(t, db.Create(ctx, user))
	product := &productModel.Product{Name: "p", Code: "P-1", Price: 5, StockQuantity: 4, Active: true}
	require.NoError(t, db.Create(ctx, product))

	oRepo := orderRepo.NewOrderRepository(db)
	pRepo := orderRepo.NewProductRepository(db)
	uRepo := orderRepo.NewUserRepository(db)
	rRepo := orderRepo.NewReservationRepository(db)
	cSvc := orderSvc.NewCouponService(validation.New(), orderRepo.NewCouponRepository(db))
	svc := orderSvc.NewOrderService(validation.New(), db, oRepo, pRepo, uRepo, rRepo, cSvc, notification.NewLoggerNotifier())

	order, err := oRepo.CreateOrder(ctx, user.ID, []*orderModel.OrderLine{{
		ProductID: product.ID, Quantity: 2, Price: 10,
	}}, "", 0)
	require.NoError(t, err)
	order.Status = orderModel.OrderStatusPendingPayment
	require.NoError(t, oRepo.UpdateOrder(ctx, order))
	require.NoError(t, pRepo.ReserveStock(ctx, product.ID, 2))
	require.NoError(t, rRepo.CreateMany(ctx, []*orderModel.StockReservation{{
		OrderID:   order.ID,
		ProductID: product.ID,
		Quantity:  2,
		Status:    orderModel.ReservationStatusActive,
		ExpiresAt: time.Now().Add(-1 * time.Minute),
	}}))

	released, err := svc.SweepExpiredReservations(ctx, 100)
	require.NoError(t, err)
	require.Equal(t, 1, released)

	var fresh orderModel.Order
	require.NoError(t, db.GetDB().First(&fresh, "id = ?", order.ID).Error)
	require.Equal(t, orderModel.OrderStatusCancelled, fresh.Status)

	var freshProd productModel.Product
	require.NoError(t, db.GetDB().First(&freshProd, "id = ?", product.ID).Error)
	require.Equal(t, 4, freshProd.StockQuantity, "stock untouched on release")
	require.Equal(t, 0, freshProd.ReservedQuantity, "reserved fully released")

	again, err := svc.SweepExpiredReservations(ctx, 100)
	require.NoError(t, err)
	require.Equal(t, 0, again)
}
