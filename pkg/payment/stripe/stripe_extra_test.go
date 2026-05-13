package stripe

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"goshop/pkg/payment"
)

func TestNewProvider_AppliesDefaults(t *testing.T) {
	p := NewProvider(Config{})
	assert.Equal(t, defaultAPIBase, p.apiBase)
	require.NotNil(t, p.httpClient)
	assert.Equal(t, 10*time.Second, p.httpClient.Timeout)
}

func TestNewProvider_TrimsTrailingSlashFromAPIBase(t *testing.T) {
	p := NewProvider(Config{APIBase: "https://example.com/"})
	assert.Equal(t, "https://example.com", p.apiBase)
}

func TestCreateIntent_HTTPError_DecodesStripeError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":{"type":"invalid_request_error","code":"parameter_missing","message":"amount required"}}`))
	}))
	defer srv.Close()

	p := NewProvider(Config{SecretKey: "sk_test", APIBase: srv.URL})
	_, err := p.CreateIntent(t.Context(), payment.CreateIntentParams{Amount: 100, Currency: "USD", OrderID: "o"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "status=400")
	assert.Contains(t, err.Error(), "amount required")
}

func TestCreateIntent_TransportError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	srv.Close() // close immediately so Do() fails

	p := NewProvider(Config{SecretKey: "sk_test", APIBase: srv.URL})
	_, err := p.CreateIntent(t.Context(), payment.CreateIntentParams{Amount: 100, Currency: "USD", OrderID: "o"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "stripe create intent")
}

func TestCreateIntent_DecodeError_ReturnsError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`not json`))
	}))
	defer srv.Close()

	p := NewProvider(Config{SecretKey: "sk_test", APIBase: srv.URL})
	_, err := p.CreateIntent(t.Context(), payment.CreateIntentParams{Amount: 100, Currency: "USD", OrderID: "o"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "decode body")
}

func TestVerifyWebhook_BadBodyAfterValidSignature(t *testing.T) {
	p := NewProvider(Config{WebhookSecret: "whsec_test"})
	now := time.Unix(1700000000, 0)
	p.now = func() time.Time { return now }

	body := []byte(`not json`)
	header := sign("whsec_test", now.Unix(), body)

	_, err := p.VerifyWebhook(body, header)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "decode body")
}

func TestParseStripeSignature_BadTimestamp(t *testing.T) {
	_, _, err := parseStripeSignature("t=notanumber,v1=abc")
	assert.ErrorIs(t, err, payment.ErrInvalidSignature)
}

func TestParseStripeSignature_MissingTimestamp(t *testing.T) {
	_, _, err := parseStripeSignature("v1=abc")
	assert.ErrorIs(t, err, payment.ErrInvalidSignature)
}

func TestParseStripeSignature_MissingV1(t *testing.T) {
	_, _, err := parseStripeSignature("t=1700000000")
	assert.ErrorIs(t, err, payment.ErrInvalidSignature)
}

func TestParseStripeSignature_IgnoresMalformedParts(t *testing.T) {
	now := int64(1700000000)
	ts, sigs, err := parseStripeSignature("t=1700000000,foo,v1=abc")
	require.NoError(t, err)
	assert.Equal(t, now, ts)
	assert.Equal(t, []string{"abc"}, sigs)
}
