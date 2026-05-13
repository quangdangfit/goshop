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

func TestRetryConfig_WithDefaults(t *testing.T) {
	got := RetryConfig{}.withDefaults()
	if got.MaxAttempts != 3 || got.InitialDelay != 100*time.Millisecond {
		t.Fatalf("zero-value defaults wrong: %+v", got)
	}
	custom := RetryConfig{MaxAttempts: 5, InitialDelay: 7 * time.Millisecond}.withDefaults()
	if custom.MaxAttempts != 5 || custom.InitialDelay != 7*time.Millisecond {
		t.Fatalf("custom values overridden: %+v", custom)
	}
}

func TestRetryingNotifier_NilDLQ_DoesNotPanic(t *testing.T) {
	inner := &stubNotifier{failFor: 5}
	n := NewRetryingNotifier(inner, RetryConfig{MaxAttempts: 2, InitialDelay: time.Millisecond}, nil)
	if err := n.SendOrderPlaced(context.Background(), "o", "u@x"); err == nil {
		t.Fatal("expected exhaustion error")
	}
}

func TestRetryingNotifier_ContextCancelled_StopsEarly(t *testing.T) {
	inner := &stubNotifier{failFor: 100}
	dlq := &recordingDLQ{}
	n := NewRetryingNotifier(inner, RetryConfig{MaxAttempts: 10, InitialDelay: 50 * time.Millisecond}, dlq)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(10 * time.Millisecond)
		cancel()
	}()
	err := n.SendOrderPlaced(ctx, "o", "u@x")
	if err == nil {
		t.Fatal("expected error")
	}
	if inner.calls >= 10 {
		t.Fatalf("ctx cancellation should stop retries early, got %d calls", inner.calls)
	}
	if len(dlq.records) != 1 {
		t.Fatalf("DLQ should still record on early termination, got %d", len(dlq.records))
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
