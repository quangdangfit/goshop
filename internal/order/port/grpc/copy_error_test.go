package grpc

import (
	"context"
	"math"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"goshop/internal/order/model"
	srvMocks "goshop/internal/order/service/mocks"
	"goshop/pkg/paging"
	pb "goshop/proto/gen/go/order"
)

func nanOrder() *model.Order {
	return &model.Order{ID: "o1", TotalPrice: math.NaN(), Status: model.OrderStatusNew}
}

func ctxWithUser() context.Context {
	return context.WithValue(context.Background(), "userId", "u1") //nolint:staticcheck
}

func TestGRPC_PlaceOrder_CopyError(t *testing.T) {
	svc := srvMocks.NewOrderService(t)
	svc.On("PlaceOrder", mock.Anything, mock.Anything).Return(nanOrder(), nil).Once()
	h := NewOrderHandler(svc)
	_, err := h.PlaceOrder(ctxWithUser(), &pb.PlaceOrderReq{
		Lines: []*pb.PlaceOrderLineReq{{ProductId: "p1", Quantity: 1}},
	})
	require.Error(t, err)
}

func TestGRPC_GetOrderByID_CopyError(t *testing.T) {
	svc := srvMocks.NewOrderService(t)
	svc.On("GetOrderByID", mock.Anything, "o1").Return(nanOrder(), nil).Once()
	h := NewOrderHandler(svc)
	_, err := h.GetOrderByID(ctxWithUser(), &pb.GetOrderByIDReq{Id: "o1"})
	require.Error(t, err)
}

func TestGRPC_GetMyOrders_CopyError(t *testing.T) {
	svc := srvMocks.NewOrderService(t)
	svc.On("GetMyOrders", mock.Anything, mock.Anything).
		Return([]*model.Order{nanOrder()}, &paging.Pagination{}, nil).Once()
	h := NewOrderHandler(svc)
	_, err := h.GetMyOrders(ctxWithUser(), &pb.GetMyOrdersReq{})
	require.Error(t, err)
}

func TestGRPC_CancelOrder_CopyError(t *testing.T) {
	svc := srvMocks.NewOrderService(t)
	svc.On("CancelOrder", mock.Anything, "o1", "u1").Return(nanOrder(), nil).Once()
	h := NewOrderHandler(svc)
	_, err := h.CancelOrder(ctxWithUser(), &pb.CancelOrderReq{Id: "o1"})
	require.Error(t, err)
}
