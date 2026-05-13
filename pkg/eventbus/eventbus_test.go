package eventbus

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestPublishFansOutToAllSubscribers(t *testing.T) {
	bus := New()
	var wg sync.WaitGroup
	wg.Add(2)

	var got1, got2 string
	bus.Subscribe(TopicOrderCreated, func(_ context.Context, ev Event) {
		defer wg.Done()
		got1 = ev.(OrderCreated).OrderID
	})
	bus.Subscribe(TopicOrderCreated, func(_ context.Context, ev Event) {
		defer wg.Done()
		got2 = ev.(OrderCreated).OrderID
	})

	bus.Publish(context.Background(), OrderCreated{OrderID: "abc"})

	done := make(chan struct{})
	go func() { wg.Wait(); close(done) }()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for subscribers")
	}
	if got1 != "abc" || got2 != "abc" {
		t.Fatalf("got %q / %q, want abc / abc", got1, got2)
	}
}

func TestPublishWithNoSubscribersIsNoop(t *testing.T) {
	bus := New()
	bus.Publish(context.Background(), LowStock{ProductID: "x"})
}

func TestDefault_LazilyInitializesSingleton(t *testing.T) {
	SetDefault(nil)
	t.Cleanup(func() { SetDefault(nil) })

	a := Default()
	b := Default()
	if a == nil {
		t.Fatal("Default returned nil")
	}
	if a != b {
		t.Fatal("Default should return the same instance on repeated calls")
	}
}

func TestSetDefault_OverridesSingleton(t *testing.T) {
	SetDefault(nil)
	t.Cleanup(func() { SetDefault(nil) })

	custom := New()
	SetDefault(custom)
	if Default() != custom {
		t.Fatal("SetDefault did not install the supplied bus")
	}
}

func TestEventTopics(t *testing.T) {
	cases := []struct {
		ev   Event
		want string
	}{
		{OrderCreated{}, TopicOrderCreated},
		{OrderPaid{}, TopicOrderPaid},
		{OrderCancelled{}, TopicOrderCancelled},
		{LowStock{}, TopicLowStock},
	}
	for _, c := range cases {
		if got := c.ev.Topic(); got != c.want {
			t.Errorf("%T.Topic() = %q, want %q", c.ev, got, c.want)
		}
	}
}

func TestSlowSubscriberDoesNotBlockPublisher(t *testing.T) {
	bus := New()
	bus.Subscribe(TopicOrderPaid, func(_ context.Context, _ Event) {
		time.Sleep(200 * time.Millisecond)
	})

	start := time.Now()
	bus.Publish(context.Background(), OrderPaid{OrderID: "1"})
	if time.Since(start) > 50*time.Millisecond {
		t.Fatalf("Publish blocked on slow subscriber: %s", time.Since(start))
	}
}
