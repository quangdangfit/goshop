package domain

import "goshop/internal/order/model"

// CouponFromModel returns the API DTO for a coupon, hiding internal columns
// (deleted_at) and removing the unreachable utils.Copy error branches that
// JSON-roundtrip conversion produced.
func CouponFromModel(m *model.Coupon) *Coupon {
	if m == nil {
		return nil
	}
	return &Coupon{
		ID:             m.ID,
		Code:           m.Code,
		DiscountType:   string(m.DiscountType),
		DiscountValue:  m.DiscountValue,
		MinOrderAmount: m.MinOrderAmount,
		MaxUsage:       m.MaxUsage,
		UsedCount:      m.UsedCount,
		ExpiresAt:      m.ExpiresAt,
	}
}

// OrderFromModel returns the API DTO for an order.
func OrderFromModel(m *model.Order) *Order {
	if m == nil {
		return nil
	}
	return &Order{
		ID:             m.ID,
		Code:           m.Code,
		TotalPrice:     m.TotalPrice,
		DiscountAmount: m.DiscountAmount,
		FinalPrice:     m.FinalPrice,
		CouponCode:     m.CouponCode,
		Status:         string(m.Status),
		Lines:          OrderLinesFromModel(m.Lines),
	}
}

func OrdersFromModel(rows []*model.Order) []*Order {
	out := make([]*Order, len(rows))
	for i, r := range rows {
		out[i] = OrderFromModel(r)
	}
	return out
}

func OrderLineFromModel(m *model.OrderLine) *OrderLine {
	if m == nil {
		return nil
	}
	out := &OrderLine{
		Quantity: m.Quantity,
		Price:    m.Price,
	}
	if m.Product != nil {
		out.Product = Product{
			ID:    m.Product.ID,
			Code:  m.Product.Code,
			Name:  m.Product.Name,
			Price: m.Product.Price,
		}
	}
	return out
}

func OrderLinesFromModel(lines []*model.OrderLine) []*OrderLine {
	out := make([]*OrderLine, len(lines))
	for i, l := range lines {
		out[i] = OrderLineFromModel(l)
	}
	return out
}
