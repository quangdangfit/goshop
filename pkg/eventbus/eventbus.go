// Package eventbus is a lightweight in-process pub/sub for decoupling domains
// from notification transports. Domains publish typed events; subscribers (e.g. the
// notification service) react asynchronously. The interface is intentionally
// minimal so we can swap in NATS/Kafka later without touching publishers.
package eventbus

import (
	"context"
	"sync"
)

// Event is any value carrying a Topic. Concrete events are defined in events.go.
type Event interface {
	Topic() string
}

// Handler reacts to a single event. It runs on a background goroutine; errors
// should be logged by the handler — the bus will not surface them to the publisher.
type Handler func(ctx context.Context, ev Event)

// Bus is the public contract. Publish never blocks the caller waiting on subscribers.
type Bus interface {
	Subscribe(topic string, h Handler)
	Publish(ctx context.Context, ev Event)
}

type inproc struct {
	mu   sync.RWMutex
	subs map[string][]Handler
}

// New returns an in-process bus. Handlers run on detached goroutines so a slow
// subscriber cannot block the publisher.
func New() Bus {
	return &inproc{subs: make(map[string][]Handler)}
}

func (b *inproc) Subscribe(topic string, h Handler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.subs[topic] = append(b.subs[topic], h)
}

// Default returns the process-wide singleton bus, lazily initialized. Domains that
// want to publish without dragging the bus through every constructor can call this;
// main.go is responsible for subscribing handlers at startup.
func Default() Bus {
	defaultMu.Lock()
	defer defaultMu.Unlock()
	if defaultBus == nil {
		defaultBus = New()
	}
	return defaultBus
}

// SetDefault overrides the process-wide bus. Tests use this to install a fresh bus.
func SetDefault(b Bus) {
	defaultMu.Lock()
	defer defaultMu.Unlock()
	defaultBus = b
}

var (
	defaultMu  sync.Mutex
	defaultBus Bus
)

func (b *inproc) Publish(ctx context.Context, ev Event) {
	b.mu.RLock()
	handlers := append([]Handler(nil), b.subs[ev.Topic()]...)
	b.mu.RUnlock()

	for _, h := range handlers {
		h := h
		go h(ctx, ev)
	}
}
