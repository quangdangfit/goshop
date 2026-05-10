package eventbus

const (
	TopicOrderCreated   = "order.created"
	TopicOrderPaid      = "order.paid"
	TopicOrderCancelled = "order.cancelled"
	TopicLowStock       = "inventory.low_stock"
)

type OrderCreated struct {
	OrderID   string
	UserID    string
	UserEmail string
}

func (OrderCreated) Topic() string { return TopicOrderCreated }

type OrderPaid struct {
	OrderID   string
	UserID    string
	UserEmail string
}

func (OrderPaid) Topic() string { return TopicOrderPaid }

type OrderCancelled struct {
	OrderID   string
	UserID    string
	UserEmail string
	Reason    string
}

func (OrderCancelled) Topic() string { return TopicOrderCancelled }

// LowStock fires when a product's available stock (stock - reserved) crosses below
// the configured threshold. Carries enough context for an admin notification.
type LowStock struct {
	ProductID string
	Available int
	Threshold int
}

func (LowStock) Topic() string { return TopicLowStock }
