package repositories

import (
	"context"
	"errors"
	"time"

	"github.com/jinzhu/copier"
	"gorm.io/gorm"

	"goshop/app/models"
	"goshop/app/serializers"
	"goshop/dbs"
)

type IProductRepository interface {
	Create(ctx context.Context, req *models.Product) error

	ListProducts(ctx context.Context, req serializers.ListProductReq) (*[]models.Product, error)
	GetProductByID(ctx context.Context, id string) (*models.Product, error)
	Update(ctx context.Context, id string, req *serializers.CreateProductReq) (*models.Product, error)
}

type ProductRepo struct {
	db *gorm.DB
}

func NewProductRepository() *ProductRepo {
	return &ProductRepo{db: dbs.Database}
}

func (r *ProductRepo) ListProducts(ctx context.Context, req serializers.ListProductReq) (*[]models.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var products []models.Product
	if err := r.db.Where(req).Find(&products).Error; err != nil {
		return nil, nil
	}

	return &products, nil
}

func (r *ProductRepo) GetProductByID(ctx context.Context, id string) (*models.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var product models.Product
	if err := r.db.Where("id = ?", id).Find(&product).Error; err != nil {
		return nil, errors.New("not found product")
	}

	return &product, nil
}

func (r *ProductRepo) Create(ctx context.Context, product *models.Product) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := r.db.Create(&product).Error; err != nil {
		return err
	}

	return nil
}

func (r *ProductRepo) Update(ctx context.Context, id string, req *serializers.CreateProductReq) (*models.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var product models.Product
	if err := r.db.Where("id = ? ", id).First(&product).Error; err != nil {
		return nil, errors.New("not found product")
	}

	copier.Copy(&product, &req)
	if err := r.db.Save(&product).Error; err != nil {
		return nil, err
	}

	return &product, nil
}
