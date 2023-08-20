package service

import (
	"context"

	"github.com/quangdangfit/gocommon/logger"
	"github.com/quangdangfit/gocommon/validation"

	"goshop/internal/product/dto"
	"goshop/internal/product/model"
	"goshop/internal/product/repository"
	"goshop/pkg/paging"
	"goshop/pkg/utils"
)

//go:generate mockery --name=IProductService
type IProductService interface {
	ListProducts(c context.Context, req *dto.ListProductReq) ([]*model.Product, *paging.Pagination, error)
	GetProductByID(ctx context.Context, id string) (*model.Product, error)
	Create(ctx context.Context, req *dto.CreateProductReq) (*model.Product, error)
	Update(ctx context.Context, id string, req *dto.UpdateProductReq) (*model.Product, error)
}

type ProductService struct {
	validator validation.Validation
	repo      repository.IProductRepository
}

func NewProductService(
	validator validation.Validation,
	repo repository.IProductRepository,
) *ProductService {
	return &ProductService{
		validator: validator,
		repo:      repo,
	}
}

func (p *ProductService) GetProductByID(ctx context.Context, id string) (*model.Product, error) {
	product, err := p.repo.GetProductByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (p *ProductService) ListProducts(ctx context.Context, req *dto.ListProductReq) ([]*model.Product, *paging.Pagination, error) {
	products, pagination, err := p.repo.ListProducts(ctx, req)
	if err != nil {
		return nil, nil, err
	}

	return products, pagination, nil
}

func (p *ProductService) Create(ctx context.Context, req *dto.CreateProductReq) (*model.Product, error) {
	if err := p.validator.ValidateStruct(req); err != nil {
		return nil, err
	}

	var product model.Product
	utils.Copy(&product, req)

	err := p.repo.Create(ctx, &product)
	if err != nil {
		logger.Errorf("Create fail, error: %s", err)
		return nil, err
	}

	return &product, nil
}

func (p *ProductService) Update(ctx context.Context, id string, req *dto.UpdateProductReq) (*model.Product, error) {
	if err := p.validator.ValidateStruct(req); err != nil {
		return nil, err
	}

	product, err := p.repo.GetProductByID(ctx, id)
	if err != nil {
		logger.Errorf("Update.GetUserByID fail, id: %s, error: %s", id, err)
		return nil, err
	}

	utils.Copy(product, req)
	err = p.repo.Update(ctx, product)
	if err != nil {
		logger.Errorf("Update fail, id: %s, error: %s", id, err)
		return nil, err
	}

	return product, nil
}
