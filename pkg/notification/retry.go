package notification

import (
	"context"
	"time"

	"github.com/quangdangfit/gocommon/logger"
)

// DeadLetterSink is the contract for sinks that store notifications which exhausted retries.
// Production wires this to a DB-backed sink; tests use an in-memory implementation.
type DeadLetterSink interface {
	Record(ctx context.Context, eventType, userEmail, payload string, lastErr error) error
}

// RetryConfig tunes the retrying notifier. Defaults are exponential backoff: 100ms, 200ms,
// 400ms — bounded so a slow downstream cannot stall callers.
type RetryConfig struct {
	MaxAttempts  int
	InitialDelay time.Duration
}

func (c RetryConfig) withDefaults() RetryConfig {
	if c.MaxAttempts <= 0 {
		c.MaxAttempts = 3
	}
	if c.InitialDelay <= 0 {
		c.InitialDelay = 100 * time.Millisecond
	}
	return c
}

// RetryingNotifier wraps an inner Notifier with bounded retry + DLQ-on-exhaustion. The
// dead-letter sink is invoked once per logical event after the final retry fails so
// operators can replay it later.
type RetryingNotifier struct {
	inner Notifier
	cfg   RetryConfig
	dlq   DeadLetterSink
}

func NewRetryingNotifier(inner Notifier, cfg RetryConfig, dlq DeadLetterSink) Notifier {
	return &RetryingNotifier{inner: inner, cfg: cfg.withDefaults(), dlq: dlq}
}

func (r *RetryingNotifier) SendOrderPlaced(ctx context.Context, orderID, userEmail string) error {
	return r.run(ctx, "order_placed", userEmail, orderID, func() error {
		return r.inner.SendOrderPlaced(ctx, orderID, userEmail)
	})
}

func (r *RetryingNotifier) SendOrderStatusChanged(ctx context.Context, orderID, userEmail, newStatus string) error {
	return r.run(ctx, "order_status_changed", userEmail, orderID+"|"+newStatus, func() error {
		return r.inner.SendOrderStatusChanged(ctx, orderID, userEmail, newStatus)
	})
}

func (r *RetryingNotifier) run(ctx context.Context, eventType, userEmail, payload string, op func() error) error {
	delay := r.cfg.InitialDelay
	var lastErr error
	for attempt := 1; attempt <= r.cfg.MaxAttempts; attempt++ {
		if err := op(); err == nil {
			return nil
		} else {
			lastErr = err
		}
		if attempt == r.cfg.MaxAttempts {
			break
		}
		select {
		case <-ctx.Done():
			lastErr = ctx.Err()
			delay = 0 // skip the next backoff; we're shutting down
		case <-time.After(delay):
		}
		if ctx.Err() != nil {
			break
		}
		delay *= 2
	}
	logger.Errorf("notifier exhausted retries for %s/%s: %s", eventType, userEmail, lastErr)
	if r.dlq != nil {
		if dlqErr := r.dlq.Record(ctx, eventType, userEmail, payload, lastErr); dlqErr != nil {
			logger.Errorf("dead-letter sink failed for %s/%s: %s", eventType, userEmail, dlqErr)
		}
	}
	return lastErr
}
