package grpc

import (
	"goshop/internal/order/model"
	pb "goshop/proto/gen/go/order"
)

// orderInfoFromModel returns the gRPC OrderInfo for a storage order. Removes the
// unreachable utils.Copy error branches that a JSON-roundtrip conversion would
// produce; all field copies here are total functions.
func orderInfoFromModel(m *model.Order) *pb.OrderInfo {
	if m == nil {
		return nil
	}
	return &pb.OrderInfo{
		Id:         m.ID,
		Code:       m.Code,
		UserId:     m.UserID,
		Lines:      orderLineInfosFromModel(m.Lines),
		TotalPrice: float32(m.TotalPrice),
		Status:     string(m.Status),
	}
}

func ordersInfoFromModel(rows []*model.Order) []*pb.OrderInfo {
	out := make([]*pb.OrderInfo, len(rows))
	for i, r := range rows {
		out[i] = orderInfoFromModel(r)
	}
	return out
}

func orderLineInfosFromModel(lines []*model.OrderLine) []*pb.OrderLineInfo {
	out := make([]*pb.OrderLineInfo, len(lines))
	for i, l := range lines {
		info := &pb.OrderLineInfo{
			ProductId: l.ProductID,
			Quantity:  uint32(l.Quantity), //nolint:gosec // quantity is a small positive integer
			Price:     float32(l.Price),
		}
		if l.Product != nil {
			info.ProductName = l.Product.Name
		}
		out[i] = info
	}
	return out
}
