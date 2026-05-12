package service

import (
	"context"
	"errors"

	"github.com/quangdangfit/gocommon/logger"
	"github.com/quangdangfit/gocommon/validation"

	"goshop/internal/product/domain"
	"goshop/internal/product/model"
	"goshop/internal/product/repository"
	"goshop/pkg/paging"
)

var errInvalidStockQty = errors.New("stock quantity must be positive")

//go:generate mockery --name=ProductService
type ProductService interface {
	ListProducts(c context.Context, req *domain.ListProductReq) ([]*model.Product, *paging.Pagination, error)
	GetProductByID(ctx context.Context, id string) (*model.Product, error)
	Create(ctx context.Context, req *domain.CreateProductReq) (*model.Product, error)
	Update(ctx context.Context, id string, req *domain.UpdateProductReq) (*model.Product, error)
	AddStock(ctx context.Context, id string, qty int, adminUserID string) (*model.Product, error)
}

type productSvc struct {
	validator validation.Validation
	repo      repository.ProductRepository
}

func NewProductService(
	validator validation.Validation,
	repo repository.ProductRepository,
) ProductService {
	return &productSvc{
		validator: validator,
		repo:      repo,
	}
}

func (p *productSvc) GetProductByID(ctx context.Context, id string) (*model.Product, error) {
	product, err := p.repo.GetProductByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (p *productSvc) ListProducts(ctx context.Context, req *domain.ListProductReq) ([]*model.Product, *paging.Pagination, error) {
	products, pagination, err := p.repo.ListProducts(ctx, req)
	if err != nil {
		return nil, nil, err
	}

	return products, pagination, nil
}

func (p *productSvc) Create(ctx context.Context, req *domain.CreateProductReq) (*model.Product, error) {
	if err := p.validator.ValidateStruct(req); err != nil {
		return nil, err
	}

	product := model.Product{
		Name:          req.Name,
		Description:   req.Description,
		Price:         req.Price,
		StockQuantity: req.StockQuantity,
		Images:        req.Images,
	}
	if req.CategoryID != "" {
		cid := req.CategoryID
		product.CategoryID = &cid
	}

	if err := p.repo.Create(ctx, &product); err != nil {
		logger.Errorf("Create fail, error: %s", err)
		return nil, err
	}

	return &product, nil
}

// AddStock atomically increases a product's stock_quantity. The adminUserID is included in
// the audit log so we can trace who restocked what; not persisted yet (a stock_audit table
// is a follow-up).
func (p *productSvc) AddStock(ctx context.Context, id string, qty int, adminUserID string) (*model.Product, error) {
	if qty <= 0 {
		return nil, errInvalidStockQty
	}
	if err := p.repo.AddStock(ctx, id, qty); err != nil {
		return nil, err
	}
	logger.Infof("admin restock: admin=%s product=%s qty=+%d", adminUserID, id, qty)
	return p.repo.GetProductByID(ctx, id)
}

func (p *productSvc) Update(ctx context.Context, id string, req *domain.UpdateProductReq) (*model.Product, error) {
	if err := p.validator.ValidateStruct(req); err != nil {
		return nil, err
	}

	product, err := p.repo.GetProductByID(ctx, id)
	if err != nil {
		logger.Errorf("Update.GetUserByID fail, id: %s, error: %s", id, err)
		return nil, err
	}

	if req.Name != "" {
		product.Name = req.Name
	}
	if req.Description != "" {
		product.Description = req.Description
	}
	if req.Price != 0 {
		product.Price = req.Price
	}
	if req.StockQuantity != nil {
		product.StockQuantity = *req.StockQuantity
	}
	if req.Images != nil {
		product.Images = req.Images
	}
	if req.CategoryID != "" {
		cid := req.CategoryID
		product.CategoryID = &cid
	}
	err = p.repo.Update(ctx, product)
	if err != nil {
		logger.Errorf("Update fail, id: %s, error: %s", id, err)
		return nil, err
	}

	return product, nil
}
