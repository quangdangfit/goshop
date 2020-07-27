package repositories

import (
	"errors"

	"github.com/jinzhu/copier"
	"github.com/jinzhu/gorm"

	"goshop/app/models"
	"goshop/app/schema"
	"goshop/dbs"
)

type IProductRepository interface {
	GetProducts(params schema.ProductQueryParams) (*[]models.Product, error)
	GetProductByID(uuid string) (*models.Product, error)
	GetProductByCategory(uuid string, active bool) (*[]models.Product, error)
	CreateProduct(req *models.ProductRequest) (*models.Product, error)
	UpdateProduct(uuid string, req *models.ProductRequest) (*models.Product, error)
}

type productRepo struct {
	db *gorm.DB
}

func NewProductRepository() IProductRepository {
	return &productRepo{db: dbs.Database}
}

func (r *productRepo) GetProducts(params schema.ProductQueryParams) (*[]models.Product, error) {
	var products []models.Product
	if r.db.Where(params).Find(&products).RecordNotFound() {
		return nil, nil
	}

	return &products, nil
}

func (r *productRepo) GetProductByCategory(categUUID string, active bool) (*[]models.Product, error) {
	var products []models.Product
	if r.db.Where("active = ? AND categ_uuid = ?", active, categUUID).Find(&products).RecordNotFound() {
		return nil, nil
	}

	return &products, nil
}

func (r *productRepo) GetProductByID(uuid string) (*models.Product, error) {
	var product models.Product
	if r.db.Where("uuid = ?", uuid).Find(&product).RecordNotFound() {
		return nil, errors.New("not found product")
	}

	return &product, nil
}

func (r *productRepo) CreateProduct(req *models.ProductRequest) (*models.Product, error) {
	var product models.Product
	copier.Copy(&product, &req)

	if err := r.db.Create(&product).Error; err != nil {
		return nil, err
	}

	return &product, nil
}

func (r *productRepo) UpdateProduct(uuid string, req *models.ProductRequest) (*models.Product, error) {
	var product models.Product
	if r.db.Where("uuid = ? ", uuid).First(&product).RecordNotFound() {
		return nil, errors.New("not found product")
	}

	copier.Copy(&product, &req)
	if err := r.db.Save(&product).Error; err != nil {
		return nil, err
	}

	return &product, nil
}
