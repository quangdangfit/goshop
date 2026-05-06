package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/quangdangfit/gocommon/validation"

	orderRepository "goshop/internal/order/repository"
	orderService "goshop/internal/order/service"
	"goshop/internal/payment/repository"
	"goshop/internal/payment/service"
	"goshop/pkg/config"
	"goshop/pkg/dbs"
	"goshop/pkg/middleware"
	"goshop/pkg/notification"
	"goshop/pkg/payment"
	stripeProvider "goshop/pkg/payment/stripe"
	"goshop/pkg/response"
)

// Routes wires the payment domain. Uses the live config to construct a Stripe provider; the
// webhook route deliberately sits outside the JWT middleware (Stripe authenticates via the
// signature header instead).
func Routes(r *gin.RouterGroup, db dbs.Database, validator validation.Validation) {
	cfg := config.GetConfig()
	provider := stripeProvider.NewProvider(stripeProvider.Config{
		SecretKey:     cfg.StripeSecretKey,
		WebhookSecret: cfg.StripeWebhookSecret,
		APIBase:       cfg.StripeAPIBase,
	})

	paymentRepo := repository.NewPaymentRepository(db)

	// Build a minimal OrderService for MarkOrderPaid / UpdateOrderStatus on webhook events.
	orderSvc := orderService.NewOrderService(
		validator, db,
		orderRepository.NewOrderRepository(db),
		orderRepository.NewProductRepository(db),
		orderRepository.NewUserRepository(db),
		orderRepository.NewReservationRepository(db),
		orderService.NewCouponService(validator, orderRepository.NewCouponRepository(db)),
		notification.NewLoggerNotifier(),
	)

	paymentSvc := service.NewPaymentService(provider, paymentRepo, orderSvc, orderSvc)
	handler := NewHandler(paymentSvc)

	authMiddleware := middleware.JWTAuth()

	// /orders/:id/payment-intent — authenticated, used by the customer to start checkout.
	r.POST("/orders/:id/payment-intent", authMiddleware, handler.CreatePaymentIntent)

	// /webhooks/stripe — public, signature-verified.
	r.POST("/webhooks/stripe", handler.StripeWebhook)

	// /config/public — exposes the Stripe publishable key to the FE.
	r.GET("/config/public", func(c *gin.Context) {
		response.JSON(c, http.StatusOK, gin.H{
			"stripe_publishable_key": cfg.StripePublishableKey,
		})
	})

	// Touch payment to silence the unused import lint when the package is imported but no
	// type from it is referenced directly. (Used transitively via the provider.)
	_ = payment.ErrInvalidSignature
}
