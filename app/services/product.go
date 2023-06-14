package services

import (
	"context"

	"github.com/quangdangfit/gocommon/logger"

	"goshop/app/models"
	"goshop/app/repositories"
	"goshop/app/serializers"
	"goshop/pkg/paging"
	"goshop/pkg/utils"
)

type IProductService interface {
	ListProducts(c context.Context, req *serializers.ListProductReq) ([]*models.Product, *paging.Pagination, error)
	GetProductByID(ctx context.Context, id string) (*models.Product, error)
	Create(ctx context.Context, req *serializers.CreateProductReq) (*models.Product, error)
	Update(ctx context.Context, id string, req *serializers.UpdateProductReq) (*models.Product, error)
}

type ProductService struct {
	repo repositories.IProductRepository
}

func NewProductService(repo repositories.IProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (p *ProductService) GetProductByID(ctx context.Context, id string) (*models.Product, error) {
	product, err := p.repo.GetProductByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (p *ProductService) ListProducts(ctx context.Context, req *serializers.ListProductReq) ([]*models.Product, *paging.Pagination, error) {
	products, pagination, err := p.repo.ListProducts(ctx, req)
	if err != nil {
		return nil, nil, err
	}

	return products, pagination, nil
}

func (p *ProductService) Create(ctx context.Context, req *serializers.CreateProductReq) (*models.Product, error) {
	var product models.Product
	utils.Copy(&product, req)

	err := p.repo.Create(ctx, &product)
	if err != nil {
		logger.Errorf("Create fail, error: %s", err)
		return nil, err
	}

	return &product, nil
}

func (p *ProductService) Update(ctx context.Context, id string, req *serializers.UpdateProductReq) (*models.Product, error) {
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
