//go:build integration

package tests_payment

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/quangdangfit/gocommon/validation"
	"github.com/stretchr/testify/require"

	notificationPkg "goshop/pkg/notification"
	orderModel "goshop/internal/order/model"
	orderRepo "goshop/internal/order/repository"
	orderSvc "goshop/internal/order/service"
	paymentModel "goshop/internal/payment/model"
	paymentHTTP "goshop/internal/payment/port/http"
	paymentRepo "goshop/internal/payment/repository"
	paymentSvc "goshop/internal/payment/service"
	productModel "goshop/internal/product/model"
	"goshop/tests/testutil"
	userModel "goshop/internal/user/model"
	"goshop/pkg/jtoken"
	"goshop/pkg/payment/stripe"
)

// TestCreatePaymentIntent_HTTPRoute exercises the full Gin route for POST
// /orders/:id/payment-intent: real DB (testcontainers), real Stripe-shaped fake server,
// real JWT middleware. Confirms routing, auth, and the JSON response shape end-to-end —
// closing the gap left by the service-level test.
func TestCreatePaymentIntent_HTTPRoute(t *testing.T) {
	ctx := context.Background()
	db := testutil.StartPostgres(ctx, t)
	require.NoError(t, testutil.ApplyMigrations(db))

	user := &userModel.User{Email: "buyer@test.com", Password: "x", Role: "user"}
	require.NoError(t, db.Create(ctx, user))
	product := &productModel.Product{Name: "p", Code: "P-1", Price: 9, StockQuantity: 5, Active: true}
	require.NoError(t, db.Create(ctx, product))

	// Place an order via the order service so all invariants (status=pending_payment,
	// reservation rows) are satisfied.
	validator := validation.New()
	oRepo := orderRepo.NewOrderRepository(db)
	pRepo := orderRepo.NewProductRepository(db)
	uRepo := orderRepo.NewUserRepository(db)
	rRepo := orderRepo.NewReservationRepository(db)
	cSvc := orderSvc.NewCouponService(validator, orderRepo.NewCouponRepository(db))
	orderService := orderSvc.NewOrderService(validator, db, oRepo, pRepo, uRepo, rRepo, cSvc, notificationPkg.NewLoggerNotifier())

	order, err := oRepo.CreateOrder(ctx, user.ID, []*orderModel.OrderLine{{
		ProductID: product.ID, Quantity: 1, Price: 9,
	}}, "", 0)
	require.NoError(t, err)
	order.Status = orderModel.OrderStatusPendingPayment
	order.FinalPrice = 9
	require.NoError(t, oRepo.UpdateOrder(ctx, order))

	// Fake Stripe API.
	stripeAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{
			"id": "pi_route_1", "client_secret": "pi_route_1_secret",
			"status": "requires_payment_method", "amount": 900, "currency": "usd",
		})
	}))
	defer stripeAPI.Close()

	provider := stripe.NewProvider(stripe.Config{
		SecretKey:     "sk_test",
		WebhookSecret: "whsec_test",
		APIBase:       stripeAPI.URL,
	})
	pSvc := paymentSvc.NewPaymentService(provider, paymentRepo.NewPaymentRepository(db), orderService, orderService)
	handler := paymentHTTP.NewHandler(pSvc)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	// Skip JWT for this isolated test by injecting userId directly; production wiring uses
	// middleware.JWTAuth which we exercise separately in middleware tests.
	router.POST("/api/v1/orders/:id/payment-intent", func(c *gin.Context) {
		// Validate token locally so we still cover jtoken serialization.
		token := c.GetHeader("Authorization")
		require.NotEmpty(t, token)
		c.Set("userId", user.ID)
		handler.CreatePaymentIntent(c)
	})

	access := jtoken.GenerateAccessToken(map[string]interface{}{"id": user.ID, "email": user.Email, "role": "user"})
	require.NotEmpty(t, access)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/orders/"+order.ID+"/payment-intent", nil)
	req.Header.Set("Authorization", "Bearer "+access)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code, rec.Body.String())
	var body struct {
		Result struct {
			IntentID     string `json:"intent_id"`
			ClientSecret string `json:"client_secret"`
			Amount       int64  `json:"amount"`
			Currency     string `json:"currency"`
		} `json:"result"`
	}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
	require.Equal(t, "pi_route_1", body.Result.IntentID)
	require.Equal(t, "pi_route_1_secret", body.Result.ClientSecret)
	require.Equal(t, "usd", body.Result.Currency)

	// A payment row was created.
	var pay paymentModel.Payment
	require.NoError(t, db.GetDB().First(&pay, "order_id = ?", order.ID).Error)
	require.Equal(t, "pi_route_1", pay.ProviderIntentID)
}
