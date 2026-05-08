//go:build integration

package notification_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"goshop/pkg/notification"
)

// TestSMTPSender_AgainstMailHog: spin up MailHog (which exposes both SMTP on :1025 and a
// REST API on :8025), send an email through SMTPSender, then assert the message landed in
// MailHog by polling its messages endpoint.
//
// This proves the full SMTP wire path including RFC 5322 framing — bugs that wouldn't show
// up against a hand-rolled fake sender.
func TestSMTPSender_AgainstMailHog(t *testing.T) {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "mailhog/mailhog:v1.0.1",
		ExposedPorts: []string{"1025/tcp", "8025/tcp"},
		WaitingFor:   wait.ForListeningPort("8025/tcp").WithStartupTimeout(60 * time.Second),
	}
	c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err)
	t.Cleanup(func() {
		shutdown, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		_ = c.Terminate(shutdown)
	})

	host, err := c.Host(ctx)
	require.NoError(t, err)
	smtpPort, err := c.MappedPort(ctx, "1025/tcp")
	require.NoError(t, err)
	apiPort, err := c.MappedPort(ctx, "8025/tcp")
	require.NoError(t, err)

	smtpPortInt, err := strconv.Atoi(smtpPort.Port())
	require.NoError(t, err)
	sender := notification.NewSMTPSender(notification.SMTPConfig{
		Host: host,
		Port: smtpPortInt,
		From: "shop@goshop.test",
	})

	require.NoError(t, sender.Send(ctx, "buyer@goshop.test", "Order #ABC", "your order is in"))

	// MailHog v2 API: /api/v2/messages — poll briefly to allow the SMTP handoff to land.
	apiURL := fmt.Sprintf("http://%s:%s/api/v2/messages", host, apiPort.Port())
	deadline := time.Now().Add(10 * time.Second)
	var found bool
	for time.Now().Before(deadline) && !found {
		resp, err := http.Get(apiURL) //nolint:gosec,noctx // localhost test container
		if err != nil {
			time.Sleep(200 * time.Millisecond)
			continue
		}
		body, _ := io.ReadAll(resp.Body)
		_ = resp.Body.Close()
		var payload struct {
			Items []struct {
				Content struct{ Body string }
				To      []struct{ Mailbox, Domain string }
			}
		}
		if err := json.Unmarshal(body, &payload); err == nil {
			for _, m := range payload.Items {
				if strings.Contains(m.Content.Body, "your order is in") &&
					len(m.To) == 1 && m.To[0].Mailbox == "buyer" {
					found = true
					break
				}
			}
		}
		if !found {
			time.Sleep(200 * time.Millisecond)
		}
	}
	require.True(t, found, "expected MailHog to record the email")
}
