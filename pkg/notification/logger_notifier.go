package notification

import (
	"context"
	"fmt"

	"github.com/quangdangfit/gocommon/logger"
)

type loggerNotifier struct{}

func NewLoggerNotifier() Notifier {
	return &loggerNotifier{}
}

func (n *loggerNotifier) SendOrderPlaced(ctx context.Context, orderID, userEmail string) error {
	logger.Info(fmt.Sprintf("[Notification] Order placed: orderID=%s, user=%s", orderID, userEmail))
	return nil
}

func (n *loggerNotifier) SendOrderStatusChanged(ctx context.Context, orderID, userEmail, newStatus string) error {
	logger.Info(fmt.Sprintf("[Notification] Order status changed: orderID=%s, user=%s, status=%s", orderID, userEmail, newStatus))
	return nil
}
