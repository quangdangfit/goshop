package notification

import (
	"context"

	"github.com/quangdangfit/gocommon/logger"
)

// PreferenceChecker decides whether a user wants to receive a given event on a given channel.
// The order/notification stack stays decoupled from a concrete preferences table by depending
// on this interface; production wiring plugs in a DB-backed implementation.
type PreferenceChecker interface {
	IsEnabled(ctx context.Context, userEmail, eventType, channel string) (bool, error)
}

// AlwaysOnPreferences is the default checker: every user opts into every channel. Useful in
// tests and as a safe fallback when the preferences subsystem is unavailable.
type AlwaysOnPreferences struct{}

func (AlwaysOnPreferences) IsEnabled(_ context.Context, _, _, _ string) (bool, error) {
	return true, nil
}

const (
	channelEmail = "email"

	eventOrderPlaced  = "order_placed"
	eventOrderChanged = "order_status_changed"
)

type emailNotifier struct {
	sender EmailSender
	prefs  PreferenceChecker
}

// NewEmailNotifier returns a Notifier that sends each event as an email. Honors per-user
// preferences via the supplied checker; pass AlwaysOnPreferences{} to disable filtering.
func NewEmailNotifier(sender EmailSender, prefs PreferenceChecker) Notifier {
	if prefs == nil {
		prefs = AlwaysOnPreferences{}
	}
	return &emailNotifier{sender: sender, prefs: prefs}
}

func (n *emailNotifier) SendOrderPlaced(ctx context.Context, orderID, userEmail string) error {
	return n.send(ctx, eventOrderPlaced, userEmail, map[string]string{"OrderID": orderID})
}

func (n *emailNotifier) SendOrderStatusChanged(ctx context.Context, orderID, userEmail, newStatus string) error {
	return n.send(ctx, eventOrderChanged, userEmail, map[string]string{"OrderID": orderID, "Status": newStatus})
}

func (n *emailNotifier) send(ctx context.Context, event, userEmail string, data map[string]string) error {
	enabled, err := n.prefs.IsEnabled(ctx, userEmail, event, channelEmail)
	if err != nil {
		// Don't block delivery on a preferences read failure — log and proceed.
		logger.Warnf("notification preferences lookup failed for %s/%s: %s", userEmail, event, err)
	}
	if err == nil && !enabled {
		return nil
	}
	subject, body, err := renderTemplate(event, data)
	if err != nil {
		return err
	}
	return n.sender.Send(ctx, userEmail, subject, body)
}

// MultiNotifier fans an event out to several Notifier implementations (e.g. logger + email).
// Per-channel failures are logged but do not stop other channels from firing.
type MultiNotifier struct {
	children []Notifier
}

func NewMultiNotifier(children ...Notifier) Notifier {
	return &MultiNotifier{children: children}
}

func (m *MultiNotifier) SendOrderPlaced(ctx context.Context, orderID, userEmail string) error {
	for _, c := range m.children {
		if err := c.SendOrderPlaced(ctx, orderID, userEmail); err != nil {
			logger.Warnf("notifier child failed (order_placed): %s", err)
		}
	}
	return nil
}

func (m *MultiNotifier) SendOrderStatusChanged(ctx context.Context, orderID, userEmail, newStatus string) error {
	for _, c := range m.children {
		if err := c.SendOrderStatusChanged(ctx, orderID, userEmail, newStatus); err != nil {
			logger.Warnf("notifier child failed (order_status_changed): %s", err)
		}
	}
	return nil
}
