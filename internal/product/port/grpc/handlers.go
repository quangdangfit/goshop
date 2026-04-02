package grpc

import (
	"context"

	"github.com/quangdangfit/gocommon/logger"

	"goshop/internal/product/domain"
	"goshop/internal/product/service"
	"goshop/pkg/apperror"
	"goshop/pkg/utils"
	pb "goshop/proto/gen/go/product"
)

type ProductHandler struct {
	pb.UnimplementedProductServiceServer

	service service.ProductService
}

func NewProductHandler(service service.ProductService) *ProductHandler {
	return &ProductHandler{
		service: service,
	}
}

func (h *ProductHandler) GetProductByID(ctx context.Context, req *pb.GetProductByIDReq) (*pb.GetProductByIDRes, error) {
	if req.Id == "" {
		return nil, apperror.WrapMessage(apperror.ErrBadRequest, nil, "ID is required").GRPCStatus()
	}

	product, err := h.service.GetProductByID(ctx, req.Id)
	if err != nil {
		logger.Error("Failed to get product ", err)
		return nil, apperror.ToGRPCStatus(err)
	}

	var res pb.GetProductByIDRes
	utils.Copy(&res.Product, &product)
	return &res, nil
}

func (h *ProductHandler) ListProducts(ctx context.Context, req *pb.ListProductsReq) (*pb.ListProductsRes, error) {
	products, pagination, err := h.service.ListProducts(ctx, &domain.ListProductReq{
		Name:      req.Name,
		Code:      req.Code,
		Page:      req.Page,
		Limit:     req.Limit,
		OrderBy:   req.OrderBy,
		OrderDesc: req.OrderDesc,
	})
	if err != nil {
		logger.Error("Failed to list products ", err)
		return nil, apperror.ToGRPCStatus(err)
	}

	var res pb.ListProductsRes
	utils.Copy(&res.Products, &products)
	if pagination != nil {
		res.Total = pagination.Total
		res.CurrentPage = pagination.CurrentPage
		res.Limit = pagination.Limit
	}
	return &res, nil
}

func (h *ProductHandler) CreateProduct(ctx context.Context, req *pb.CreateProductReq) (*pb.CreateProductRes, error) {
	product, err := h.service.Create(ctx, &domain.CreateProductReq{
		Name:        req.Name,
		Description: req.Description,
		Price:       float64(req.Price),
	})
	if err != nil {
		logger.Error("Failed to create product ", err)
		return nil, apperror.ToGRPCStatus(err)
	}

	var res pb.CreateProductRes
	utils.Copy(&res.Product, &product)
	return &res, nil
}

func (h *ProductHandler) UpdateProduct(ctx context.Context, req *pb.UpdateProductReq) (*pb.UpdateProductRes, error) {
	if req.Id == "" {
		return nil, apperror.WrapMessage(apperror.ErrBadRequest, nil, "ID is required").GRPCStatus()
	}

	product, err := h.service.Update(ctx, req.Id, &domain.UpdateProductReq{
		Name:        req.Name,
		Description: req.Description,
		Price:       float64(req.Price),
	})
	if err != nil {
		logger.Error("Failed to update product ", err)
		return nil, apperror.ToGRPCStatus(err)
	}

	var res pb.UpdateProductRes
	utils.Copy(&res.Product, &product)
	return &res, nil
}
