package stripe

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"goshop/pkg/payment"
)

func sign(secret string, ts int64, body []byte) string {
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = fmt.Fprintf(mac, "%d.", ts)
	_, _ = mac.Write(body)
	return fmt.Sprintf("t=%d,v1=%s", ts, hex.EncodeToString(mac.Sum(nil)))
}

func TestVerifyWebhook_Valid(t *testing.T) {
	p := NewProvider(Config{WebhookSecret: "whsec_test"})
	now := time.Unix(1700000000, 0)
	p.now = func() time.Time { return now }

	body := []byte(`{"id":"evt_1","type":"payment_intent.succeeded","data":{"object":{"id":"pi_1","metadata":{"order_id":"ord_42"}}}}`)
	header := sign("whsec_test", now.Unix(), body)

	ev, err := p.VerifyWebhook(body, header)
	require.NoError(t, err)
	require.Equal(t, "evt_1", ev.ID)
	require.Equal(t, payment.EventPaymentSucceeded, ev.Type)
	require.Equal(t, "pi_1", ev.PaymentIntentID)
	require.Equal(t, "ord_42", ev.OrderID)
}

func TestVerifyWebhook_BadSignature(t *testing.T) {
	p := NewProvider(Config{WebhookSecret: "whsec_test"})
	p.now = func() time.Time { return time.Unix(1700000000, 0) }

	_, err := p.VerifyWebhook([]byte(`{}`), "t=1700000000,v1=deadbeef")
	require.ErrorIs(t, err, payment.ErrInvalidSignature)
}

func TestVerifyWebhook_StaleTimestamp(t *testing.T) {
	p := NewProvider(Config{WebhookSecret: "whsec_test"})
	p.now = func() time.Time { return time.Unix(1700000000, 0) }

	body := []byte(`{}`)
	stale := sign("whsec_test", 1700000000-int64(10*time.Minute/time.Second), body)
	_, err := p.VerifyWebhook(body, stale)
	require.ErrorIs(t, err, payment.ErrInvalidSignature)
}

func TestVerifyWebhook_EmptyHeader(t *testing.T) {
	p := NewProvider(Config{WebhookSecret: "whsec_test"})
	_, err := p.VerifyWebhook([]byte(`{}`), "")
	require.ErrorIs(t, err, payment.ErrInvalidSignature)
}

func TestCreateIntent_PostsExpectedForm(t *testing.T) {
	var captured struct {
		path string
		auth string
		idem string
		form map[string][]string
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		captured.path = r.URL.Path
		captured.auth = r.Header.Get("Authorization")
		captured.idem = r.Header.Get("Idempotency-Key")
		body, _ := io.ReadAll(r.Body)
		// reuse url.ParseQuery via testing-only helper:
		captured.form = parseFormBody(string(body))
		_ = json.NewEncoder(w).Encode(map[string]any{
			"id":            "pi_1",
			"client_secret": "pi_1_secret",
			"status":        "requires_payment_method",
			"amount":        1234,
			"currency":      "usd",
		})
	}))
	defer srv.Close()

	p := NewProvider(Config{SecretKey: "sk_test", APIBase: srv.URL})
	intent, err := p.CreateIntent(t.Context(), payment.CreateIntentParams{
		Amount:         1234,
		Currency:       "USD",
		OrderID:        "ord_1",
		IdempotencyKey: "order_ord_1",
	})
	require.NoError(t, err)
	require.Equal(t, "pi_1", intent.ID)
	require.Equal(t, "/v1/payment_intents", captured.path)
	require.Equal(t, "Bearer sk_test", captured.auth)
	require.Equal(t, "order_ord_1", captured.idem)
	require.Equal(t, []string{"1234"}, captured.form["amount"])
	require.Equal(t, []string{"usd"}, captured.form["currency"])
	require.Equal(t, []string{"ord_1"}, captured.form["metadata[order_id]"])
}

func parseFormBody(raw string) map[string][]string {
	out := make(map[string][]string)
	for _, kv := range strings.Split(raw, "&") {
		eq := strings.IndexByte(kv, '=')
		if eq < 0 {
			continue
		}
		k := kv[:eq]
		v := kv[eq+1:]
		// percent-decode "[" "]" which url.Values.Encode escapes
		k = strings.ReplaceAll(k, "%5B", "[")
		k = strings.ReplaceAll(k, "%5D", "]")
		out[k] = append(out[k], v)
	}
	return out
}
