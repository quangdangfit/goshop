package model

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCanTransitionTo(t *testing.T) {
	cases := []struct {
		from, to OrderStatus
		ok       bool
	}{
		{OrderStatusNew, OrderStatusNew, true}, // idempotent
		{OrderStatusNew, OrderStatusInProgress, true},
		{OrderStatusNew, OrderStatusCancelled, true},
		{OrderStatusNew, OrderStatusPendingPayment, false},
		{OrderStatusNew, OrderStatusDone, false},
		{OrderStatusPendingPayment, OrderStatusPaid, true},
		{OrderStatusPendingPayment, OrderStatusPaymentFailed, true},
		{OrderStatusPendingPayment, OrderStatusCancelled, true},
		{OrderStatusPendingPayment, OrderStatusDone, false},
		{OrderStatusPaid, OrderStatusInProgress, true},
		{OrderStatusPaid, OrderStatusCancelled, true},
		{OrderStatusInProgress, OrderStatusDone, true},
		{OrderStatusInProgress, OrderStatusCancelled, true},
		{OrderStatusPaymentFailed, OrderStatusCancelled, true},
		{OrderStatusPaymentFailed, OrderStatusPendingPayment, true},
		{OrderStatusDone, OrderStatusCancelled, false},
		{OrderStatusCancelled, OrderStatusNew, false},
		{OrderStatus("bogus"), OrderStatusNew, false},
	}
	for _, tc := range cases {
		require.Equalf(t, tc.ok, tc.from.CanTransitionTo(tc.to), "from=%s to=%s", tc.from, tc.to)
	}
}
