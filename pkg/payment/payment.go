// Package payment defines a transport-neutral interface for charging customers via a payment
// provider. The interface keeps domain code unaware of which provider is plugged in; today
// the only implementation is Stripe (see pkg/payment/stripe).
package payment

import (
	"context"
	"errors"
)

// Intent is a provider-agnostic view of a created payment intent. Only the fields the order
// flow actually needs are surfaced; Stripe's full PaymentIntent has many more.
type Intent struct {
	ID           string
	ClientSecret string
	Status       string
	Amount       int64
	Currency     string
}

// EventType enumerates the webhook events we care about. Anything else is ignored.
type EventType string

const (
	EventPaymentSucceeded EventType = "payment_intent.succeeded"
	EventPaymentFailed    EventType = "payment_intent.payment_failed"
)

// Event is the verified, normalized form of a provider webhook callback.
type Event struct {
	ID              string
	Type            EventType
	PaymentIntentID string
	OrderID         string // pulled from the intent's metadata
	Raw             []byte
}

// CreateIntentParams collects the inputs required to start a payment.
type CreateIntentParams struct {
	Amount         int64 // in minor currency units (cents for USD)
	Currency       string
	OrderID        string
	IdempotencyKey string
}

// Provider abstracts a payment processor. Webhook verification must validate the signature
// using a shared secret; callers should reject unverified events.
type Provider interface {
	CreateIntent(ctx context.Context, params CreateIntentParams) (*Intent, error)
	VerifyWebhook(payload []byte, signatureHeader string) (*Event, error)
}

// ErrInvalidSignature is returned by VerifyWebhook when the provided signature header does not
// match the expected HMAC of the body.
var ErrInvalidSignature = errors.New("invalid webhook signature")
