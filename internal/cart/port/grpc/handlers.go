package grpc

import (
	"context"
	"errors"

	"github.com/quangdangfit/gocommon/logger"

	"goshop/internal/cart/dto"
	"goshop/internal/cart/service"
	"goshop/pkg/utils"
	pb "goshop/proto/gen/go/cart"
)

type CartHandler struct {
	pb.UnimplementedCartServiceServer

	service service.ICartService
}

func NewCartHandler(service service.ICartService) *CartHandler {
	return &CartHandler{
		service: service,
	}
}

func (h *CartHandler) AddProduct(ctx context.Context, req *pb.AddProductReq) (*pb.AddProductRes, error) {
	userID, _ := ctx.Value("userId").(string)
	if userID == "" {
		return nil, errors.New("unauthorized")
	}

	cart, err := h.service.AddProduct(ctx, &dto.AddProductReq{
		UserID: userID,
		Line: &dto.CartLineReq{
			ProductID: req.ProductId,
			Quantity:  uint(req.Quantity),
		},
	})
	if err != nil {
		logger.Error("Failed to add product ", err)
		return nil, err
	}

	var res pb.AddProductRes
	utils.Copy(&res.Cart, &cart)
	return &res, nil
}

func (h *CartHandler) RemoveProduct(ctx context.Context, req *pb.RemoveProductReq) (*pb.RemoveProductRes, error) {
	userID, _ := ctx.Value("userId").(string)
	if userID == "" {
		return nil, errors.New("unauthorized")
	}

	cart, err := h.service.RemoveProduct(ctx, &dto.RemoveProductReq{
		UserID:    userID,
		ProductID: req.ProductId,
	})
	if err != nil {
		logger.Error("Failed to remove product ", err)
		return nil, err
	}

	var res pb.RemoveProductRes
	utils.Copy(&res.Cart, &cart)
	return &res, nil
}

func (h *CartHandler) GetCart(ctx context.Context, req *pb.GetCartReq) (*pb.GetCartRes, error) {
	userID, _ := ctx.Value("userId").(string)
	if userID == "" {
		return nil, errors.New("unauthorized")
	}

	cart, err := h.service.GetCartByUserID(ctx, userID)
	if err != nil {
		logger.Error("Failed to get cart ", err)
		return nil, err
	}

	var res pb.GetCartRes
	utils.Copy(&res.Cart, &cart)
	return &res, nil
}
