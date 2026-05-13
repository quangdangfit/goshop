package notification

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

type fakeSender struct {
	to, subject, body string
	err               error
	called            int
}

func (f *fakeSender) Send(_ context.Context, to, subject, body string) error {
	f.called++
	f.to, f.subject, f.body = to, subject, body
	return f.err
}

type stubPrefs struct {
	enabled bool
	err     error
}

func (s stubPrefs) IsEnabled(_ context.Context, _, _, _ string) (bool, error) {
	return s.enabled, s.err
}

func TestEmailNotifier_OrderPlaced_RendersTemplate(t *testing.T) {
	sender := &fakeSender{}
	n := NewEmailNotifier(sender, AlwaysOnPreferences{})
	require.NoError(t, n.SendOrderPlaced(context.Background(), "ord_1", "user@example.com"))
	require.Equal(t, 1, sender.called)
	require.Equal(t, "user@example.com", sender.to)
	require.Equal(t, "Order #ord_1 received", sender.subject)
	require.Contains(t, sender.body, "ord_1")
}

func TestEmailNotifier_PreferenceDisabled_Skips(t *testing.T) {
	sender := &fakeSender{}
	n := NewEmailNotifier(sender, stubPrefs{enabled: false})
	require.NoError(t, n.SendOrderPlaced(context.Background(), "ord_1", "user@example.com"))
	require.Equal(t, 0, sender.called)
}

func TestEmailNotifier_PreferenceLookupError_StillSends(t *testing.T) {
	sender := &fakeSender{}
	n := NewEmailNotifier(sender, stubPrefs{err: errors.New("db down")})
	require.NoError(t, n.SendOrderPlaced(context.Background(), "ord_1", "user@example.com"))
	require.Equal(t, 1, sender.called)
}

func TestEmailNotifier_StatusChanged(t *testing.T) {
	sender := &fakeSender{}
	n := NewEmailNotifier(sender, AlwaysOnPreferences{})
	require.NoError(t, n.SendOrderStatusChanged(context.Background(), "ord_1", "u@e.com", "paid"))
	require.Contains(t, sender.subject, "paid")
	require.Contains(t, sender.body, "paid")
}

func TestMultiNotifier_FansOut(t *testing.T) {
	a := &fakeSender{}
	b := &fakeSender{}
	notifier := NewMultiNotifier(
		NewEmailNotifier(a, AlwaysOnPreferences{}),
		NewEmailNotifier(b, AlwaysOnPreferences{}),
	)
	require.NoError(t, notifier.SendOrderPlaced(context.Background(), "ord_1", "u@e.com"))
	require.Equal(t, 1, a.called)
	require.Equal(t, 1, b.called)
}

func TestBuildDefault_NoSMTP_LoggerOnly(t *testing.T) {
	n := BuildDefault(Settings{})
	// Logger notifier returns nil — exercising the path keeps coverage and ensures the type exists.
	require.NoError(t, n.SendOrderPlaced(context.Background(), "ord", "u@e.com"))
}

func TestBuildDefault_WithSMTP_WrapsMulti(t *testing.T) {
	n := BuildDefault(Settings{
		SMTPHost:  "localhost",
		SMTPPort:  1025,
		EmailFrom: "from@x",
	})
	require.NotNil(t, n)
	_, ok := n.(*MultiNotifier)
	require.True(t, ok, "with SMTP set, BuildDefault should return a MultiNotifier")
}

type recordingDLQForFactory struct{ calls int }

func (d *recordingDLQForFactory) Record(_ context.Context, _, _, _ string, _ error) error {
	d.calls++
	return nil
}

func TestBuildDefault_WithDLQ_WrapsRetrying(t *testing.T) {
	n := BuildDefault(Settings{DLQ: &recordingDLQForFactory{}})
	require.NotNil(t, n)
	_, ok := n.(*RetryingNotifier)
	require.True(t, ok, "with a DLQ, BuildDefault should wrap in RetryingNotifier")
}

type alwaysFailingNotifier struct{ calls int }

func (a *alwaysFailingNotifier) SendOrderPlaced(_ context.Context, _, _ string) error {
	a.calls++
	return errors.New("boom")
}

func (a *alwaysFailingNotifier) SendOrderStatusChanged(_ context.Context, _, _, _ string) error {
	a.calls++
	return errors.New("boom")
}

func TestMultiNotifier_StatusChanged_LogsButDoesNotShortCircuit(t *testing.T) {
	bad := &alwaysFailingNotifier{}
	good := &fakeSender{}
	m := NewMultiNotifier(bad, NewEmailNotifier(good, AlwaysOnPreferences{}))
	require.NoError(t, m.SendOrderStatusChanged(context.Background(), "o", "u@e.com", "paid"))
	require.Equal(t, 1, bad.calls)
	require.Equal(t, 1, good.called)
}

func TestMultiNotifier_OrderPlaced_LogsButDoesNotShortCircuit(t *testing.T) {
	bad := &alwaysFailingNotifier{}
	good := &fakeSender{}
	m := NewMultiNotifier(bad, NewEmailNotifier(good, AlwaysOnPreferences{}))
	require.NoError(t, m.SendOrderPlaced(context.Background(), "o", "u@e.com"))
	require.Equal(t, 1, bad.calls)
	require.Equal(t, 1, good.called)
}
