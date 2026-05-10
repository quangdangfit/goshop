package http

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/quangdangfit/gocommon/logger"
	"github.com/stretchr/testify/require"

	"goshop/pkg/config"
	"goshop/pkg/payment"
)

type stubPayments struct {
	createFn func(ctx context.Context, orderID string) (*payment.Intent, error)
	hookFn   func(ctx context.Context, payload []byte, sig string) error
}

func (s *stubPayments) CreateIntentForOrder(ctx context.Context, o string) (*payment.Intent, error) {
	return s.createFn(ctx, o)
}
func (s *stubPayments) HandleWebhook(ctx context.Context, payload []byte, sig string) error {
	return s.hookFn(ctx, payload, sig)
}

func setupRouter(svc *stubPayments) *gin.Engine {
	logger.Initialize(config.ProductionEnv)
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewHandler(svc)
	r.POST("/orders/:id/payment-intent", h.CreatePaymentIntent)
	r.POST("/webhooks/stripe", h.StripeWebhook)
	return r
}

func TestCreatePaymentIntent_OK(t *testing.T) {
	svc := &stubPayments{createFn: func(_ context.Context, id string) (*payment.Intent, error) {
		require.Equal(t, "o1", id)
		return &payment.Intent{ID: "pi_1", ClientSecret: "cs", Amount: 1000, Currency: "usd", Status: "requires_payment_method"}, nil
	}}
	r := setupRouter(svc)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/orders/o1/payment-intent", nil))
	require.Equal(t, http.StatusOK, w.Code)
}

func TestCreatePaymentIntent_Error(t *testing.T) {
	svc := &stubPayments{createFn: func(_ context.Context, _ string) (*payment.Intent, error) {
		return nil, errors.New("nope")
	}}
	r := setupRouter(svc)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/orders/o1/payment-intent", nil))
	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestStripeWebhook_OK(t *testing.T) {
	called := false
	svc := &stubPayments{hookFn: func(_ context.Context, payload []byte, sig string) error {
		called = true
		require.Equal(t, "sig", sig)
		require.Equal(t, "{}", string(payload))
		return nil
	}}
	r := setupRouter(svc)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/webhooks/stripe", bytes.NewReader([]byte("{}")))
	req.Header.Set("Stripe-Signature", "sig")
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
	require.True(t, called)
}

func TestStripeWebhook_InvalidSignature(t *testing.T) {
	svc := &stubPayments{hookFn: func(_ context.Context, _ []byte, _ string) error { return payment.ErrInvalidSignature }}
	r := setupRouter(svc)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/webhooks/stripe", bytes.NewReader([]byte("{}")))
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestStripeWebhook_OtherError(t *testing.T) {
	svc := &stubPayments{hookFn: func(_ context.Context, _ []byte, _ string) error { return errors.New("processing failed") }}
	r := setupRouter(svc)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/webhooks/stripe", bytes.NewReader([]byte("{}")))
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusInternalServerError, w.Code)
}

type errReader struct{}

func (errReader) Read(_ []byte) (int, error) { return 0, errors.New("read failed") }
func (errReader) Close() error               { return nil }

func TestStripeWebhook_ReadBodyError(t *testing.T) {
	svc := &stubPayments{}
	r := setupRouter(svc)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/webhooks/stripe", io.NopCloser(errReader{}))
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusBadRequest, w.Code)
}
