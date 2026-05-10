package model

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPaymentBeforeCreateGeneratesID(t *testing.T) {
	p := &Payment{}
	require.NoError(t, p.BeforeCreate(nil))
	require.NotEmpty(t, p.ID)
	require.Equal(t, PaymentStatusPending, p.Status)
}

func TestPaymentBeforeCreatePreservesExistingFields(t *testing.T) {
	p := &Payment{ID: "given", Status: PaymentStatusSucceeded}
	require.NoError(t, p.BeforeCreate(nil))
	require.Equal(t, "given", p.ID)
	require.Equal(t, PaymentStatusSucceeded, p.Status)
}
