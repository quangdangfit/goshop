package grpc

import (
	"context"
	"errors"

	"github.com/quangdangfit/gocommon/logger"

	"goshop/internal/order/dto"
	"goshop/internal/order/service"
	"goshop/pkg/utils"
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
		return nil, errors.New("unauthorized")
	}

	var lines []dto.PlaceOrderLineReq
	for _, l := range req.Lines {
		lines = append(lines, dto.PlaceOrderLineReq{
			ProductID: l.ProductId,
			Quantity:  uint(l.Quantity),
		})
	}

	order, err := h.service.PlaceOrder(ctx, &dto.PlaceOrderReq{
		UserID: userID,
		Lines:  lines,
	})
	if err != nil {
		logger.Error("Failed to place order ", err)
		return nil, err
	}

	var res pb.PlaceOrderRes
	utils.Copy(&res.Order, &order)
	return &res, nil
}

func (h *OrderHandler) GetOrderByID(ctx context.Context, req *pb.GetOrderByIDReq) (*pb.GetOrderByIDRes, error) {
	userID, _ := ctx.Value("userId").(string)
	if userID == "" {
		return nil, errors.New("unauthorized")
	}

	if req.Id == "" {
		return nil, errors.New("id is required")
	}

	order, err := h.service.GetOrderByID(ctx, req.Id)
	if err != nil {
		logger.Error("Failed to get order ", err)
		return nil, err
	}

	var res pb.GetOrderByIDRes
	utils.Copy(&res.Order, &order)
	return &res, nil
}

func (h *OrderHandler) GetMyOrders(ctx context.Context, req *pb.GetMyOrdersReq) (*pb.GetMyOrdersRes, error) {
	userID, _ := ctx.Value("userId").(string)
	if userID == "" {
		return nil, errors.New("unauthorized")
	}

	orders, pagination, err := h.service.GetMyOrders(ctx, &dto.ListOrderReq{
		UserID:    userID,
		Status:    req.Status,
		Page:      req.Page,
		Limit:     req.Limit,
		OrderBy:   req.OrderBy,
		OrderDesc: req.OrderDesc,
	})
	if err != nil {
		logger.Error("Failed to get orders ", err)
		return nil, err
	}

	var res pb.GetMyOrdersRes
	utils.Copy(&res.Orders, &orders)
	if pagination != nil {
		res.Total = pagination.Total
		res.CurrentPage = pagination.CurrentPage
		res.Limit = pagination.Limit
	}
	return &res, nil
}

func (h *OrderHandler) CancelOrder(ctx context.Context, req *pb.CancelOrderReq) (*pb.CancelOrderRes, error) {
	userID, _ := ctx.Value("userId").(string)
	if userID == "" {
		return nil, errors.New("unauthorized")
	}

	if req.Id == "" {
		return nil, errors.New("id is required")
	}

	order, err := h.service.CancelOrder(ctx, req.Id, userID)
	if err != nil {
		logger.Error("Failed to cancel order ", err)
		return nil, err
	}

	var res pb.CancelOrderRes
	utils.Copy(&res.Order, &order)
	return &res, nil
}
