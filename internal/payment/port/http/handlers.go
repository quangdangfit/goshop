package http

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/quangdangfit/gocommon/logger"

	"goshop/internal/payment/service"
	"goshop/pkg/payment"
	"goshop/pkg/response"
)

// stripeSignatureHeader is the header Stripe uses for webhook HMAC delivery.
const stripeSignatureHeader = "Stripe-Signature"

type Handler struct {
	svc service.PaymentService
}

func NewHandler(svc service.PaymentService) *Handler {
	return &Handler{svc: svc}
}

type intentResponse struct {
	IntentID     string `json:"intent_id"`
	ClientSecret string `json:"client_secret"`
	Amount       int64  `json:"amount"`
	Currency     string `json:"currency"`
	Status       string `json:"status"`
}

// CreatePaymentIntent godoc
//
//	@Summary	Create a Stripe PaymentIntent for an order
//	@Tags		payments
//	@Produce	json
//	@Param		id	path		string	true	"Order ID"
//	@Success	200	{object}	intentResponse
//	@Failure	400	{object}	response.Response
//	@Router		/orders/{id}/payment-intent [post]
//	@Security	ApiKeyAuth
func (h *Handler) CreatePaymentIntent(c *gin.Context) {
	orderID := c.Param("id")
	intent, err := h.svc.CreateIntentForOrder(c.Request.Context(), orderID)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err, "create payment intent")
		return
	}
	response.JSON(c, http.StatusOK, intentResponse{
		IntentID:     intent.ID,
		ClientSecret: intent.ClientSecret,
		Amount:       intent.Amount,
		Currency:     intent.Currency,
		Status:       intent.Status,
	})
}

// StripeWebhook receives provider callbacks. Reads the raw body (signature is over the bytes
// as transmitted) and verifies before any side effect.
func (h *Handler) StripeWebhook(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err, "read webhook body")
		return
	}
	sig := c.GetHeader(stripeSignatureHeader)
	if err := h.svc.HandleWebhook(c.Request.Context(), body, sig); err != nil {
		if err == payment.ErrInvalidSignature {
			response.Error(c, http.StatusBadRequest, err, "invalid signature")
			return
		}
		logger.Errorf("stripe webhook: %s", err)
		response.Error(c, http.StatusInternalServerError, err, "webhook processing failed")
		return
	}
	response.JSON(c, http.StatusOK, gin.H{"received": true})
}
