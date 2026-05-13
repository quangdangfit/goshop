package service

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	orderModel "goshop/internal/order/model"
	orderService "goshop/internal/order/service"
	"goshop/internal/payment/model"
	"goshop/internal/payment/repository"
	"goshop/pkg/payment"
)

const stripeProvider = "stripe"

//go:generate mockery --name=PaymentService
type PaymentService interface {
	// CreateIntentForOrder creates (or reuses) a payment intent for the given order. Idempotent
	// per order: a second call returns the existing intent instead of charging twice.
	CreateIntentForOrder(ctx context.Context, orderID string) (*payment.Intent, error)
	// HandleWebhook verifies, deduplicates, and applies a provider webhook payload. Returns
	// nil for events that are valid but unrelated (ignored types, duplicates).
	HandleWebhook(ctx context.Context, payload []byte, signatureHeader string) error
}

// OrderQuery is the read-only slice of OrderService that PaymentService needs. Defined here
// to keep the dependency narrow and to avoid an import cycle with the order service.
type OrderQuery interface {
	GetOrderByID(ctx context.Context, id string) (*orderModel.Order, error)
}

type paymentService struct {
	provider     payment.Provider
	repo         repository.PaymentRepository
	orderQuery   OrderQuery
	orderService orderService.OrderService
	providerName string
}

func NewPaymentService(
	provider payment.Provider,
	repo repository.PaymentRepository,
	orderQuery OrderQuery,
	orderSvc orderService.OrderService,
) PaymentService {
	return &paymentService{
		provider:     provider,
		repo:         repo,
		orderQuery:   orderQuery,
		orderService: orderSvc,
		providerName: stripeProvider,
	}
}

func (s *paymentService) CreateIntentForOrder(ctx context.Context, orderID string) (*payment.Intent, error) {
	order, err := s.orderQuery.GetOrderByID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	if order.Status != orderModel.OrderStatusPendingPayment {
		return nil, fmt.Errorf("order %s is not pending payment (status=%s)", orderID, order.Status)
	}

	// Look up any existing payment row for this order. We deliberately do NOT short-circuit
	// on the stored row alone — the row doesn't (and shouldn't) persist client_secret, so the
	// FE needs us to re-issue the intent on every call. Stripe's idempotency key replays the
	// same intent (with its client_secret) for 24h, so a duplicate POST is safe.
	existing, err := s.repo.GetByOrderID(ctx, orderID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	amount := int64(order.FinalPrice * 100) // assume USD-style minor units
	currency := "usd"
	intent, err := s.provider.CreateIntent(ctx, payment.CreateIntentParams{
		Amount:         amount,
		Currency:       currency,
		OrderID:        order.ID,
		IdempotencyKey: "order_" + order.ID,
	})
	if err != nil {
		return nil, err
	}

	// First time for this order — persist a payment row so webhooks can look it up.
	// On subsequent calls the row already exists and Stripe returned the same intent
	// via idempotency replay, so there's nothing to write.
	if existing == nil {
		rec := &model.Payment{
			OrderID:          order.ID,
			Provider:         s.providerName,
			ProviderIntentID: intent.ID,
			Amount:           amount,
			Currency:         currency,
			Status:           model.PaymentStatusPending,
		}
		if err := s.repo.Create(ctx, rec); err != nil {
			return nil, err
		}
	}
	return intent, nil
}

func (s *paymentService) HandleWebhook(ctx context.Context, payload []byte, signatureHeader string) error {
	event, err := s.provider.VerifyWebhook(payload, signatureHeader)
	if err != nil {
		return err
	}

	// Dedup before doing any side-effects.
	if err := s.repo.RecordProviderEvent(ctx, s.providerName, event.ID); err != nil {
		if errors.Is(err, repository.ErrEventAlreadyProcessed) {
			return nil
		}
		return err
	}

	if event.OrderID == "" {
		return fmt.Errorf("webhook event %s missing order_id metadata", event.ID)
	}
	rec, err := s.repo.GetByOrderID(ctx, event.OrderID)
	if err != nil {
		return err
	}

	switch event.Type {
	case payment.EventPaymentSucceeded:
		rec.Status = model.PaymentStatusSucceeded
		if err := s.repo.Update(ctx, rec); err != nil {
			return err
		}
		if _, err := s.orderService.MarkOrderPaid(ctx, event.OrderID); err != nil {
			return err
		}
	case payment.EventPaymentFailed:
		rec.Status = model.PaymentStatusFailed
		if err := s.repo.Update(ctx, rec); err != nil {
			return err
		}
		if _, err := s.orderService.UpdateOrderStatus(ctx, event.OrderID, orderModel.OrderStatusPaymentFailed); err != nil {
			return err
		}
	case payment.EventPaymentCanceled:
		// Customer or system aborted the intent. Cancel the order and let the sweeper / order
		// flow release any remaining stock reservations.
		rec.Status = model.PaymentStatusCanceled
		if err := s.repo.Update(ctx, rec); err != nil {
			return err
		}
		if _, err := s.orderService.UpdateOrderStatus(ctx, event.OrderID, orderModel.OrderStatusCancelled); err != nil {
			return err
		}
	case payment.EventPaymentProcessing:
		// Async payment methods (e.g. bank debits) sit in processing. Reflect it on the payment
		// record so the FE can show the right state; order stays pending_payment.
		rec.Status = model.PaymentStatusProcessing
		if err := s.repo.Update(ctx, rec); err != nil {
			return err
		}
	case payment.EventPaymentRequiresAction:
		// 3DS / SCA challenge — the customer needs to come back to the payment page. No state
		// change to the order; just mark the payment so we don't auto-cancel.
		rec.Status = model.PaymentStatusRequiresAction
		if err := s.repo.Update(ctx, rec); err != nil {
			return err
		}
	default:
		// Unhandled event type — already deduped, nothing else to do.
	}
	return nil
}
