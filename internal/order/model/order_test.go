package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrder_BeforeCreate(t *testing.T) {
	tests := []struct {
		name           string
		order          *Order
		expectedStatus OrderStatus
	}{
		{
			name:           "DefaultStatus",
			order:          &Order{},
			expectedStatus: OrderStatusNew,
		},
		{
			name:           "ExistingStatus",
			order:          &Order{Status: OrderStatusInProgress},
			expectedStatus: OrderStatusInProgress,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.order.BeforeCreate(nil)
			assert.NoError(t, err)
			assert.NotEmpty(t, tc.order.ID)
			assert.NotEmpty(t, tc.order.Code)
			assert.Equal(t, tc.expectedStatus, tc.order.Status)
		})
	}
}

func TestOrderLine_BeforeCreate(t *testing.T) {
	line := &OrderLine{}
	err := line.BeforeCreate(nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, line.ID)
}

func TestOrderStatus_IsValid(t *testing.T) {
	tests := []struct {
		name   string
		status OrderStatus
		want   bool
	}{
		{name: "new", status: OrderStatusNew, want: true},
		{name: "in-progress", status: OrderStatusInProgress, want: true},
		{name: "done", status: OrderStatusDone, want: true},
		{name: "cancelled", status: OrderStatusCancelled, want: true},
		{name: "empty string", status: OrderStatus(""), want: false},
		{name: "unknown", status: OrderStatus("shipped"), want: false},
		{name: "wrong casing", status: OrderStatus("NEW"), want: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, tc.status.IsValid())
		})
	}
}
