package repositories

import (
	"context"

	"gorm.io/gorm"

	"goshop/app/dbs"
	"goshop/app/models"
	"goshop/app/serializers"
	"goshop/config"
	"goshop/pkg/paging"
)

type IProductRepository interface {
	Create(ctx context.Context, product *models.Product) error
	Update(ctx context.Context, product *models.Product) error
	ListProducts(ctx context.Context, req serializers.ListProductReq) ([]*models.Product, *paging.Pagination, error)
	GetProductByID(ctx context.Context, id string) (*models.Product, error)
}

type ProductRepo struct {
	db *gorm.DB
}

func NewProductRepository() *ProductRepo {
	return &ProductRepo{db: dbs.Database}
}

func (r *ProductRepo) ListProducts(ctx context.Context, req serializers.ListProductReq) ([]*models.Product, *paging.Pagination, error) {
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
	if err := query.Model(&models.Product{}).Count(&total).Error; err != nil {
		return nil, nil, err
	}

	pagination := paging.New(req.Page, req.Limit, total)

	var products []*models.Product
	if err := query.
		Limit(int(pagination.Limit)).
		Offset(int(pagination.Skip)).
		Order(order).
		Find(&products).Error; err != nil {
		return nil, nil, nil
	}

	return products, pagination, nil
}

func (r *ProductRepo) GetProductByID(ctx context.Context, id string) (*models.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, config.DatabaseTimeout)
	defer cancel()

	var product models.Product
	if err := r.db.Where("id = ?", id).First(&product).Error; err != nil {
		return nil, err
	}

	return &product, nil
}

func (r *ProductRepo) Create(ctx context.Context, product *models.Product) error {
	ctx, cancel := context.WithTimeout(ctx, config.DatabaseTimeout)
	defer cancel()

	if err := r.db.Create(&product).Error; err != nil {
		return err
	}

	return nil
}

func (r *ProductRepo) Update(ctx context.Context, product *models.Product) error {
	ctx, cancel := context.WithTimeout(ctx, config.DatabaseTimeout)
	defer cancel()

	if err := r.db.Save(&product).Error; err != nil {
		return err
	}

	return nil
}
