package notification

import "context"

//go:generate mockery --name=Notifier
type Notifier interface {
	SendOrderPlaced(ctx context.Context, orderID, userEmail string) error
	SendOrderStatusChanged(ctx context.Context, orderID, userEmail, newStatus string) error
}
