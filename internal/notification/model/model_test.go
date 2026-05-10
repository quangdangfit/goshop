package model

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPreferenceBeforeCreateGeneratesIDWhenEmpty(t *testing.T) {
	p := &Preference{}
	require.NoError(t, p.BeforeCreate(nil))
	require.NotEmpty(t, p.ID)
}

func TestPreferenceBeforeCreatePreservesExistingID(t *testing.T) {
	p := &Preference{ID: "fixed-id"}
	require.NoError(t, p.BeforeCreate(nil))
	require.Equal(t, "fixed-id", p.ID)
}

func TestDeadLetterBeforeCreateGeneratesIDWhenEmpty(t *testing.T) {
	d := &DeadLetterNotification{}
	require.NoError(t, d.BeforeCreate(nil))
	require.NotEmpty(t, d.ID)
}

func TestDeadLetterBeforeCreatePreservesExistingID(t *testing.T) {
	d := &DeadLetterNotification{ID: "fixed-dl"}
	require.NoError(t, d.BeforeCreate(nil))
	require.Equal(t, "fixed-dl", d.ID)
}
