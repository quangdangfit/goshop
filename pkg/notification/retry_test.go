package notification

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"
)

type recordingDLQ struct {
	mu      sync.Mutex
	records []string
}

func (d *recordingDLQ) Record(_ context.Context, eventType, userEmail, payload string, lastErr error) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.records = append(d.records, eventType+"|"+userEmail+"|"+payload+"|"+lastErr.Error())
	return nil
}

type stubNotifier struct {
	calls   int
	failFor int
}

func (s *stubNotifier) SendOrderPlaced(_ context.Context, _, _ string) error {
	s.calls++
	if s.calls <= s.failFor {
		return errors.New("transient")
	}
	return nil
}

func (s *stubNotifier) SendOrderStatusChanged(_ context.Context, _, _, _ string) error {
	s.calls++
	if s.calls <= s.failFor {
		return errors.New("transient")
	}
	return nil
}

func TestRetryingNotifier_SucceedsAfterRetries(t *testing.T) {
	inner := &stubNotifier{failFor: 1} // fail once, then succeed
	dlq := &recordingDLQ{}
	n := NewRetryingNotifier(inner, RetryConfig{MaxAttempts: 3, InitialDelay: time.Millisecond}, dlq)

	if err := n.SendOrderPlaced(context.Background(), "o1", "a@x.com"); err != nil {
		t.Fatalf("expected success, got %v", err)
	}
	if inner.calls != 2 {
		t.Fatalf("expected 2 attempts, got %d", inner.calls)
	}
	if len(dlq.records) != 0 {
		t.Fatalf("DLQ should not be invoked on success")
	}
}

func TestRetryingNotifier_DLQOnExhaustion(t *testing.T) {
	inner := &stubNotifier{failFor: 5} // always fail within 3 attempts
	dlq := &recordingDLQ{}
	n := NewRetryingNotifier(inner, RetryConfig{MaxAttempts: 3, InitialDelay: time.Millisecond}, dlq)

	err := n.SendOrderStatusChanged(context.Background(), "o1", "a@x.com", "paid")
	if err == nil {
		t.Fatal("expected exhaustion error")
	}
	if inner.calls != 3 {
		t.Fatalf("expected 3 attempts, got %d", inner.calls)
	}
	if len(dlq.records) != 1 {
		t.Fatalf("DLQ should record once, got %d", len(dlq.records))
	}
}
