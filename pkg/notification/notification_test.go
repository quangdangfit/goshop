package notification

import (
	"context"
	"testing"

	"github.com/quangdangfit/gocommon/logger"
	"github.com/stretchr/testify/assert"

	"goshop/pkg/config"
)

func init() {
	logger.Initialize(config.ProductionEnv)
}

func TestNewLoggerNotifier(t *testing.T) {
	n := NewLoggerNotifier()
	assert.NotNil(t, n)
}

func TestLoggerNotifier_Send(t *testing.T) {
	tests := []struct {
		name string
		send func(n Notifier) error
	}{
		{
			name: "OrderPlaced",
			send: func(n Notifier) error {
				return n.SendOrderPlaced(context.Background(), "order-123", "user@example.com")
			},
		},
		{
			name: "OrderStatusChanged",
			send: func(n Notifier) error {
				return n.SendOrderStatusChanged(context.Background(), "order-123", "user@example.com", "done")
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			n := NewLoggerNotifier()
			err := tc.send(n)
			assert.NoError(t, err)
		})
	}
}
