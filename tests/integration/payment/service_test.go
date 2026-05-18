//go:build integration

package tests_payment

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/quangdangfit/gocommon/validation"
	"github.com/stretchr/testify/require"

	notificationPkg "goshop/pkg/notification"
	orderModel "goshop/internal/order/model"
	orderRepo "goshop/internal/order/repository"
	orderSvc "goshop/internal/order/service"
	paymentModel "goshop/internal/payment/model"
	paymentRepo "goshop/internal/payment/repository"
	paymentSvc "goshop/internal/payment/service"
	productModel "goshop/internal/product/model"
	"goshop/tests/testutil"
	userModel "goshop/internal/user/model"
	"goshop/pkg/payment"
	"goshop/pkg/payment/stripe"
)

// TestStripeWebhook_PaidFlow exercises: place order (pending_payment + reservation) →
// create intent (against an httptest server impersonating Stripe) → simulate
// payment_intent.succeeded webhook → assert order=paid, stock committed, payment=succeeded.
//
// Runs against a real Postgres (testcontainers) and a fake Stripe (httptest.Server) so the
// HMAC signing path and DB transactions are both real.
func TestStripeWebhook_PaidFlow(t *testing.T) {
	ctx := context.Background()
	db := testutil.StartPostgres(ctx, t)
	require.NoError(t, db.AutoMigrate(
		&userModel.User{},
		&productModel.Product{}, &productModel.Category{},
		orderModel.Order{}, orderModel.OrderLine{}, orderModel.Coupon{}, orderModel.StockReservation{},
		&paymentModel.Payment{}, &paymentModel.ProviderEvent{},
	))

	// Seed a user + product.
	user := &userModel.User{Email: "buyer@test.com", Password: "x"}
	require.NoError(t, db.Create(ctx, user))
	product := &productModel.Product{Name: "p1", Code: "P-1", Price: 10, StockQuantity: 5, Active: true}
	require.NoError(t, db.Create(ctx, product))

	// Fake Stripe API: the only call is POST /v1/payment_intents during CreateIntentForOrder.
	stripeAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/v1/payment_intents", r.URL.Path)
		_ = r.Body.Close()
		_ = json.NewEncoder(w).Encode(map[string]any{
			"id": "pi_test_1", "client_secret": "pi_test_1_secret",
			"status": "requires_payment_method", "amount": 2000, "currency": "usd",
		})
	}))
	defer stripeAPI.Close()

	provider := stripe.NewProvider(stripe.Config{
		SecretKey:     "sk_test",
		WebhookSecret: "whsec_test",
		APIBase:       stripeAPI.URL,
	})

	// Build the order service stack and place an order.
	validator := validation.New()
	oRepo := orderRepo.NewOrderRepository(db)
	pRepo := orderRepo.NewProductRepository(db)
	uRepo := orderRepo.NewUserRepository(db)
	rRepo := orderRepo.NewReservationRepository(db)
	cSvc := orderSvc.NewCouponService(validator, orderRepo.NewCouponRepository(db))
	orderService := orderSvc.NewOrderService(validator, db, oRepo, pRepo, uRepo, rRepo, cSvc, notificationPkg.NewLoggerNotifier())

	order, err := oRepo.CreateOrder(ctx, user.ID, []*orderModel.OrderLine{{
		ProductID: product.ID, Quantity: 2, Price: 20,
	}}, "", 0)
	require.NoError(t, err)
	order.Status = orderModel.OrderStatusPendingPayment
	order.FinalPrice = 20
	require.NoError(t, oRepo.UpdateOrder(ctx, order))
	// Reserve manually since we bypassed the full service for setup brevity.
	require.NoError(t, pRepo.ReserveStock(ctx, product.ID, 2))
	require.NoError(t, rRepo.CreateMany(ctx, []*orderModel.StockReservation{{
		OrderID: order.ID, ProductID: product.ID, Quantity: 2,
		Status: orderModel.ReservationStatusActive, ExpiresAt: time.Now().Add(15 * time.Minute),
	}}))

	pSvc := paymentSvc.NewPaymentService(provider, paymentRepo.NewPaymentRepository(db), orderService, orderService)

	intent, err := pSvc.CreateIntentForOrder(ctx, order.ID)
	require.NoError(t, err)
	require.Equal(t, "pi_test_1", intent.ID)

	// Build a signed webhook for payment_intent.succeeded.
	body := []byte(fmt.Sprintf(`{"id":"evt_1","type":"payment_intent.succeeded","data":{"object":{"id":"pi_test_1","metadata":{"order_id":"%s"}}}}`, order.ID))
	ts := time.Now().Unix()
	mac := hmac.New(sha256.New, []byte("whsec_test"))
	_, _ = fmt.Fprintf(mac, "%d.", ts)
	_, _ = mac.Write(body)
	header := fmt.Sprintf("t=%d,v1=%s", ts, hex.EncodeToString(mac.Sum(nil)))

	require.NoError(t, pSvc.HandleWebhook(ctx, body, header))

	// Replaying the same event must be a no-op (idempotency).
	require.NoError(t, pSvc.HandleWebhook(ctx, body, header))

	// Assertions: order paid, stock decremented, reserved cleared, payment succeeded.
	var fresh orderModel.Order
	require.NoError(t, db.GetDB().First(&fresh, "id = ?", order.ID).Error)
	require.Equal(t, orderModel.OrderStatusPaid, fresh.Status)

	var freshProd productModel.Product
	require.NoError(t, db.GetDB().First(&freshProd, "id = ?", product.ID).Error)
	require.Equal(t, 3, freshProd.StockQuantity, "stock decreased by committed qty")
	require.Equal(t, 0, freshProd.ReservedQuantity, "reservation released into commit")

	var pay paymentModel.Payment
	require.NoError(t, db.GetDB().First(&pay, "order_id = ?", order.ID).Error)
	require.Equal(t, paymentModel.PaymentStatusSucceeded, pay.Status)

	// drain a body to silence unused-io warning
	_, _ = io.Discard.Write(nil)
	_ = payment.ErrInvalidSignature
}
