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

func TestLoggerNotifier_SendOrderPlaced(t *testing.T) {
	n := NewLoggerNotifier()
	err := n.SendOrderPlaced(context.Background(), "order-123", "user@example.com")
	assert.NoError(t, err)
}

func TestLoggerNotifier_SendOrderStatusChanged(t *testing.T) {
	n := NewLoggerNotifier()
	err := n.SendOrderStatusChanged(context.Background(), "order-123", "user@example.com", "done")
	assert.NoError(t, err)
}
