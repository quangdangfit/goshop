package repository

import (
	"context"

	"gorm.io/gorm"

	"goshop/config"
	"goshop/internal/product/dto"
	"goshop/internal/product/model"
	"goshop/pkg/paging"
)

//go:generate mockery --name=IProductRepository
type IProductRepository interface {
	Create(ctx context.Context, product *model.Product) error
	Update(ctx context.Context, product *model.Product) error
	ListProducts(ctx context.Context, req *dto.ListProductReq) ([]*model.Product, *paging.Pagination, error)
	GetProductByID(ctx context.Context, id string) (*model.Product, error)
}

type ProductRepo struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepo {
	_ = db.AutoMigrate(&model.Product{})
	return &ProductRepo{db: db}
}

func (r *ProductRepo) ListProducts(ctx context.Context, req *dto.ListProductReq) ([]*model.Product, *paging.Pagination, error) {
	ctx, cancel := context.WithTimeout(ctx, config.DatabaseTimeout)
	defer cancel()

	query := r.db
	order := "created_at"
	if req.Name != "" {
		query = query.Where("name LIKE ?", "%"+req.Name+"%")
	}
	if req.Code != "" {
		query = query.Where("code = ?", req.Code)
	}
	if req.OrderBy != "" {
		order = req.OrderBy
		if req.OrderDesc {
			order += " DESC"
		}
	}
	var total int64
	if err := query.Model(&model.Product{}).Count(&total).Error; err != nil {
		return nil, nil, err
	}

	pagination := paging.New(req.Page, req.Limit, total)

	var products []*model.Product
	if err := query.
		Limit(int(pagination.Limit)).
		Offset(int(pagination.Skip)).
		Order(order).
		Find(&products).Error; err != nil {
		return nil, nil, nil
	}

	return products, pagination, nil
}

func (r *ProductRepo) GetProductByID(ctx context.Context, id string) (*model.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, config.DatabaseTimeout)
	defer cancel()

	var product model.Product
	if err := r.db.Where("id = ?", id).First(&product).Error; err != nil {
		return nil, err
	}

	return &product, nil
}

func (r *ProductRepo) Create(ctx context.Context, product *model.Product) error {
	ctx, cancel := context.WithTimeout(ctx, config.DatabaseTimeout)
	defer cancel()

	if err := r.db.Create(&product).Error; err != nil {
		return err
	}

	return nil
}

func (r *ProductRepo) Update(ctx context.Context, product *model.Product) error {
	ctx, cancel := context.WithTimeout(ctx, config.DatabaseTimeout)
	defer cancel()

	if err := r.db.Save(&product).Error; err != nil {
		return err
	}

	return nil
}
