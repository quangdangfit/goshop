package services

import (
	"context"

	"github.com/jinzhu/copier"

	"goshop/app/models"
	"goshop/app/repositories"
	"goshop/app/serializers"
)

type IProductService interface {
	ListProducts(c context.Context, req serializers.ListProductReq) (*[]models.Product, error)
	GetProductByID(ctx context.Context, id string) (*models.Product, error)
	Create(ctx context.Context, req *serializers.CreateProductReq) (*models.Product, error)
	Update(ctx context.Context, id string, req *serializers.CreateProductReq) (*models.Product, error)
}

type ProductRepo struct {
	repo repositories.IProductRepository
}

func NewProductService(repo repositories.IProductRepository) IProductService {
	return &ProductRepo{repo: repo}
}

func (p *ProductRepo) GetProductByID(ctx context.Context, id string) (*models.Product, error) {
	product, err := p.repo.GetProductByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (p *ProductRepo) ListProducts(ctx context.Context, req serializers.ListProductReq) (*[]models.Product, error) {
	products, err := p.repo.ListProducts(ctx, req)
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (p *ProductRepo) Create(ctx context.Context, req *serializers.CreateProductReq) (*models.Product, error) {
	var product models.Product
	err := copier.Copy(&product, req)
	if err != nil {
		return nil, err
	}

	err = p.repo.Create(ctx, &product)
	if err != nil {
		return nil, err
	}

	return &product, nil
}

func (p *ProductRepo) Update(ctx context.Context, id string, req *serializers.CreateProductReq) (*models.Product, error) {
	product, err := p.repo.Update(ctx, id, req)
	if err != nil {
		return nil, err
	}

	return product, nil
}
