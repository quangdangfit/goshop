package notification

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/smtp"
	"strings"
	"text/template"
)

// EmailSender is the minimal contract an email transport must satisfy. Real implementations
// dial an SMTP server; tests can supply an in-memory fake to capture sent messages.
type EmailSender interface {
	Send(ctx context.Context, to, subject, body string) error
}

// SMTPConfig is the dial information for plain SMTP/STARTTLS. Auth is optional — leave User
// and Password empty for unauthenticated relays (e.g. MailHog in tests).
type SMTPConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	From     string
}

type smtpSender struct {
	cfg SMTPConfig
}

func NewSMTPSender(cfg SMTPConfig) EmailSender {
	return &smtpSender{cfg: cfg}
}

// Send dispatches one message synchronously. Callers that need fan-out, retries, or async
// delivery should wrap this in a higher-level worker; this is intentionally bare-bones so it
// can be exercised against a real MTA.
func (s *smtpSender) Send(ctx context.Context, to, subject, body string) error {
	if s.cfg.Host == "" || s.cfg.From == "" {
		return errors.New("smtp not configured")
	}
	addr := fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.Port)
	msg := buildRFC5322Message(s.cfg.From, to, subject, body)

	var auth smtp.Auth
	if s.cfg.User != "" {
		auth = smtp.PlainAuth("", s.cfg.User, s.cfg.Password, s.cfg.Host)
	}
	// Honor ctx cancellation by short-circuiting before dial.
	if err := ctx.Err(); err != nil {
		return err
	}
	return smtp.SendMail(addr, auth, s.cfg.From, []string{to}, msg)
}

func buildRFC5322Message(from, to, subject, body string) []byte {
	var b bytes.Buffer
	fmt.Fprintf(&b, "From: %s\r\n", from)
	fmt.Fprintf(&b, "To: %s\r\n", to)
	fmt.Fprintf(&b, "Subject: %s\r\n", subject)
	b.WriteString("MIME-Version: 1.0\r\n")
	b.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
	b.WriteString("\r\n")
	b.WriteString(body)
	return b.Bytes()
}

// emailTemplates holds the rendered subject/body for each event type. Kept inline so changes
// don't require a build-tag-protected embed; long term these should move to template files.
var emailTemplates = map[string]struct {
	Subject string
	Body    string
}{
	"order_placed": {
		Subject: "Order #{{.OrderID}} received",
		Body:    "Hi,\n\nWe've received your order {{.OrderID}} and reserved your items for 15 minutes while we process payment.\n\nThanks for shopping with GoShop.",
	},
	"order_status_changed": {
		Subject: "Order #{{.OrderID}}: now {{.Status}}",
		Body:    "Hi,\n\nYour order {{.OrderID}} is now in status: {{.Status}}.\n\nThanks,\nGoShop",
	},
}

func renderTemplate(name string, data any) (subject, body string, err error) {
	t, ok := emailTemplates[name]
	if !ok {
		return "", "", fmt.Errorf("unknown template %q", name)
	}
	subject, err = renderOne(t.Subject, data)
	if err != nil {
		return "", "", err
	}
	body, err = renderOne(t.Body, data)
	return subject, body, err
}

func renderOne(tpl string, data any) (string, error) {
	t, err := template.New("").Parse(tpl)
	if err != nil {
		return "", err
	}
	var buf strings.Builder
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
