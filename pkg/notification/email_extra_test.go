package notification

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewSMTPSender_ReturnsNonNil(t *testing.T) {
	s := NewSMTPSender(SMTPConfig{Host: "h", Port: 25, From: "f@x"})
	assert.NotNil(t, s)
}

func TestSMTPSender_Send_NotConfigured(t *testing.T) {
	s := NewSMTPSender(SMTPConfig{})
	err := s.Send(context.Background(), "to@x", "s", "b")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "smtp not configured")
}

func TestSMTPSender_Send_ContextCancelled(t *testing.T) {
	s := NewSMTPSender(SMTPConfig{Host: "h", Port: 25, From: "f@x", User: "u", Password: "p"})
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := s.Send(ctx, "to@x", "s", "b")
	assert.ErrorIs(t, err, context.Canceled)
}

func TestBuildRFC5322Message_Format(t *testing.T) {
	msg := string(buildRFC5322Message("from@x", "to@y", "subj", "hello"))
	assert.Contains(t, msg, "From: from@x\r\n")
	assert.Contains(t, msg, "To: to@y\r\n")
	assert.Contains(t, msg, "Subject: subj\r\n")
	assert.Contains(t, msg, "MIME-Version: 1.0\r\n")
	assert.True(t, strings.HasSuffix(msg, "hello"))
}

func TestRenderTemplate_UnknownTemplate(t *testing.T) {
	_, _, err := renderTemplate("does_not_exist", nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unknown template")
}

func TestRenderOne_ParseError(t *testing.T) {
	_, err := renderOne("{{ .X ", nil)
	require.Error(t, err)
}

func TestRenderOne_ExecuteError(t *testing.T) {
	// Calling a method that doesn't exist on the data triggers an execute-time error.
	_, err := renderOne("{{ .Missing.Field }}", struct{}{})
	require.Error(t, err)
}
