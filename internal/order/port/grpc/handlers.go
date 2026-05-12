package grpc

import (
	"context"

	"github.com/quangdangfit/gocommon/logger"

	"goshop/internal/order/domain"
	"goshop/internal/order/service"
	"goshop/pkg/apperror"
	pb "goshop/proto/gen/go/order"
)

type OrderHandler struct {
	pb.UnimplementedOrderServiceServer

	service service.OrderService
}

func NewOrderHandler(service service.OrderService) *OrderHandler {
	return &OrderHandler{
		service: service,
	}
}

func (h *OrderHandler) PlaceOrder(ctx context.Context, req *pb.PlaceOrderReq) (*pb.PlaceOrderRes, error) {
	userID, _ := ctx.Value("userId").(string)
	if userID == "" {
		return nil, apperror.ErrUnauthorized.GRPCStatus()
	}

	var lines []domain.PlaceOrderLineReq
	for _, l := range req.Lines {
		lines = append(lines, domain.PlaceOrderLineReq{
			ProductID: l.ProductId,
			Quantity:  uint(l.Quantity),
		})
	}

	order, err := h.service.PlaceOrder(ctx, &domain.PlaceOrderReq{
		UserID: userID,
		Lines:  lines,
	})
	if err != nil {
		logger.Error("Failed to place order ", err)
		return nil, apperror.ToGRPCStatus(err)
	}

	return &pb.PlaceOrderRes{Order: orderInfoFromModel(order)}, nil
}

func (h *OrderHandler) GetOrderByID(ctx context.Context, req *pb.GetOrderByIDReq) (*pb.GetOrderByIDRes, error) {
	userID, _ := ctx.Value("userId").(string)
	if userID == "" {
		return nil, apperror.ErrUnauthorized.GRPCStatus()
	}

	if req.Id == "" {
		return nil, apperror.WrapMessage(apperror.ErrBadRequest, nil, "ID is required").GRPCStatus()
	}

	order, err := h.service.GetOrderByID(ctx, req.Id)
	if err != nil {
		logger.Error("Failed to get order ", err)
		return nil, apperror.ToGRPCStatus(err)
	}

	return &pb.GetOrderByIDRes{Order: orderInfoFromModel(order)}, nil
}

func (h *OrderHandler) GetMyOrders(ctx context.Context, req *pb.GetMyOrdersReq) (*pb.GetMyOrdersRes, error) {
	userID, _ := ctx.Value("userId").(string)
	if userID == "" {
		return nil, apperror.ErrUnauthorized.GRPCStatus()
	}

	orders, pagination, err := h.service.GetMyOrders(ctx, &domain.ListOrderReq{
		UserID:    userID,
		Status:    req.Status,
		Page:      req.Page,
		Limit:     req.Limit,
		OrderBy:   req.OrderBy,
		OrderDesc: req.OrderDesc,
	})
	if err != nil {
		logger.Error("Failed to get orders ", err)
		return nil, apperror.ToGRPCStatus(err)
	}

	res := &pb.GetMyOrdersRes{Orders: ordersInfoFromModel(orders)}
	if pagination != nil {
		res.Total = pagination.Total
		res.CurrentPage = pagination.CurrentPage
		res.Limit = pagination.Limit
	}
	return res, nil
}

func (h *OrderHandler) CancelOrder(ctx context.Context, req *pb.CancelOrderReq) (*pb.CancelOrderRes, error) {
	userID, _ := ctx.Value("userId").(string)
	if userID == "" {
		return nil, apperror.ErrUnauthorized.GRPCStatus()
	}

	if req.Id == "" {
		return nil, apperror.WrapMessage(apperror.ErrBadRequest, nil, "ID is required").GRPCStatus()
	}

	order, err := h.service.CancelOrder(ctx, req.Id, userID)
	if err != nil {
		logger.Error("Failed to cancel order ", err)
		return nil, apperror.ToGRPCStatus(err)
	}

	return &pb.CancelOrderRes{Order: orderInfoFromModel(order)}, nil
}
