// Package stripe implements payment.Provider against Stripe's REST API. Hand-rolled HTTP
// rather than the stripe-go SDK to keep the dependency surface small and to make the
// behavior easy to fake with stripe-mock in integration tests.
package stripe

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"goshop/pkg/payment"
)

const (
	defaultAPIBase  = "https://api.stripe.com"
	signatureMaxAge = 5 * time.Minute
)

// Provider is a payment.Provider backed by Stripe's REST API.
type Provider struct {
	secretKey     string
	webhookSecret string
	apiBase       string
	httpClient    *http.Client
	now           func() time.Time
}

// Config holds the inputs to NewProvider. APIBase is optional and lets tests point at
// stripe-mock; if empty, the real Stripe endpoint is used.
type Config struct {
	SecretKey     string
	WebhookSecret string
	APIBase       string
	HTTPClient    *http.Client
}

func NewProvider(cfg Config) *Provider {
	base := cfg.APIBase
	if base == "" {
		base = defaultAPIBase
	}
	client := cfg.HTTPClient
	if client == nil {
		client = &http.Client{Timeout: 10 * time.Second}
	}
	return &Provider{
		secretKey:     cfg.SecretKey,
		webhookSecret: cfg.WebhookSecret,
		apiBase:       strings.TrimRight(base, "/"),
		httpClient:    client,
		now:           time.Now,
	}
}

type stripePaymentIntent struct {
	ID           string `json:"id"`
	ClientSecret string `json:"client_secret"`
	Status       string `json:"status"`
	Amount       int64  `json:"amount"`
	Currency     string `json:"currency"`
}

type stripeError struct {
	Error struct {
		Type    string `json:"type"`
		Message string `json:"message"`
		Code    string `json:"code"`
	} `json:"error"`
}

// CreateIntent calls POST /v1/payment_intents with the order_id stored in metadata so the
// webhook handler can match the resulting event back to an order.
func (p *Provider) CreateIntent(ctx context.Context, params payment.CreateIntentParams) (*payment.Intent, error) {
	form := url.Values{}
	form.Set("amount", strconv.FormatInt(params.Amount, 10))
	form.Set("currency", strings.ToLower(params.Currency))
	form.Set("metadata[order_id]", params.OrderID)
	// Disable redirect-based methods so the integration works with PaymentElement out of the box.
	form.Set("automatic_payment_methods[enabled]", "true")
	form.Set("automatic_payment_methods[allow_redirects]", "never")

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.apiBase+"/v1/payment_intents", strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+p.secretKey)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if params.IdempotencyKey != "" {
		req.Header.Set("Idempotency-Key", params.IdempotencyKey)
	}

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("stripe create intent: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		var se stripeError
		_ = json.Unmarshal(body, &se)
		return nil, fmt.Errorf("stripe create intent: status=%d type=%s code=%s message=%s",
			resp.StatusCode, se.Error.Type, se.Error.Code, se.Error.Message)
	}

	var pi stripePaymentIntent
	if err := json.Unmarshal(body, &pi); err != nil {
		return nil, fmt.Errorf("stripe create intent: decode body: %w", err)
	}
	return &payment.Intent{
		ID:           pi.ID,
		ClientSecret: pi.ClientSecret,
		Status:       pi.Status,
		Amount:       pi.Amount,
		Currency:     pi.Currency,
	}, nil
}

// stripeEvent mirrors the relevant subset of the Stripe webhook event envelope.
type stripeEvent struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	Data struct {
		Object stripePaymentIntent `json:"object"`
	} `json:"data"`
}

// VerifyWebhook parses the Stripe-Signature header (t=...,v1=...) and validates its HMAC
// against the configured webhook secret. Rejects payloads older than signatureMaxAge to
// blunt replay attacks.
func (p *Provider) VerifyWebhook(payloadBytes []byte, signatureHeader string) (*payment.Event, error) {
	timestamp, sigs, err := parseStripeSignature(signatureHeader)
	if err != nil {
		return nil, err
	}
	if p.now().Sub(time.Unix(timestamp, 0)).Abs() > signatureMaxAge {
		return nil, payment.ErrInvalidSignature
	}

	mac := hmac.New(sha256.New, []byte(p.webhookSecret))
	_, _ = fmt.Fprintf(mac, "%d.", timestamp)
	_, _ = mac.Write(payloadBytes)
	expected := hex.EncodeToString(mac.Sum(nil))

	matched := false
	for _, s := range sigs {
		if hmac.Equal([]byte(s), []byte(expected)) {
			matched = true
			break
		}
	}
	if !matched {
		return nil, payment.ErrInvalidSignature
	}

	var ev stripeEvent
	if err := json.Unmarshal(payloadBytes, &ev); err != nil {
		return nil, fmt.Errorf("stripe webhook: decode body: %w", err)
	}

	// Re-decode metadata since the inner object's metadata field is provider-specific.
	var inner struct {
		Data struct {
			Object struct {
				Metadata map[string]string `json:"metadata"`
			} `json:"object"`
		} `json:"data"`
	}
	_ = json.Unmarshal(payloadBytes, &inner)

	return &payment.Event{
		ID:              ev.ID,
		Type:            payment.EventType(ev.Type),
		PaymentIntentID: ev.Data.Object.ID,
		OrderID:         inner.Data.Object.Metadata["order_id"],
		Raw:             payloadBytes,
	}, nil
}

// parseStripeSignature extracts t=<unix> and one or more v1=<hex> fields from the header.
func parseStripeSignature(header string) (int64, []string, error) {
	if header == "" {
		return 0, nil, payment.ErrInvalidSignature
	}
	var ts int64
	var v1 []string
	for _, part := range strings.Split(header, ",") {
		kv := strings.SplitN(strings.TrimSpace(part), "=", 2)
		if len(kv) != 2 {
			continue
		}
		switch kv[0] {
		case "t":
			n, err := strconv.ParseInt(kv[1], 10, 64)
			if err != nil {
				return 0, nil, payment.ErrInvalidSignature
			}
			ts = n
		case "v1":
			v1 = append(v1, kv[1])
		}
	}
	if ts == 0 || len(v1) == 0 {
		return 0, nil, payment.ErrInvalidSignature
	}
	return ts, v1, nil
}
