package services

import (
	"context"

	"goshop/app/models"
	"goshop/app/repositories"
	"goshop/app/schema"
)

type IProductService interface {
	GetProducts(c context.Context, params schema.ProductQueryParam) (*[]models.Product, error)
	GetProductByID(ctx context.Context, uuid string) (*models.Product, error)
	CreateProduct(ctx context.Context, item *schema.ProductBodyParam) (*models.Product, error)
	UpdateProduct(ctx context.Context, uuid string, item *schema.ProductBodyParam) (*models.Product, error)
	GetProductByCategoryID(ctx context.Context, uuid string) (*[]models.Product, error)
}

type product struct {
	repo repositories.IProductRepository
}

func NewProductService(repo repositories.IProductRepository) IProductService {
	return &product{repo: repo}
}

func (p *product) GetProductByID(ctx context.Context, uuid string) (*models.Product, error) {
	product, err := p.repo.GetProductByID(uuid)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (p *product) GetProducts(ctx context.Context, params schema.ProductQueryParam) (*[]models.Product, error) {
	products, err := p.repo.GetProducts(params)
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (p *product) GetProductByCategoryID(ctx context.Context, uuid string) (*[]models.Product, error) {
	products, err := p.repo.GetProductByCategoryID(uuid)
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (p *product) CreateProduct(ctx context.Context, item *schema.ProductBodyParam) (*models.Product, error) {
	product, err := p.repo.CreateProduct(item)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (p *product) UpdateProduct(ctx context.Context, uuid string, item *schema.ProductBodyParam) (*models.Product, error) {
	product, err := p.repo.UpdateProduct(uuid, item)
	if err != nil {
		return nil, err
	}

	return product, nil
}
