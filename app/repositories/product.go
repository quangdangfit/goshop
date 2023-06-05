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
	GetProducts(params schema.ProductQueryParam) (*[]models.Product, error)
	GetProductByID(uuid string) (*models.Product, error)
	GetProductByCategoryID(uuid string) (*[]models.Product, error)
	CreateProduct(item *schema.ProductBodyParam) (*models.Product, error)
	UpdateProduct(uuid string, item *schema.ProductBodyParam) (*models.Product, error)
}

type ProductRepo struct {
	db *gorm.DB
}

func NewProductRepository() *ProductRepo {
	return &ProductRepo{db: dbs.Database}
}

func (r *ProductRepo) GetProducts(params schema.ProductQueryParam) (*[]models.Product, error) {
	var products []models.Product
	if r.db.Where(params).Find(&products).RecordNotFound() {
		return nil, nil
	}

	return &products, nil
}

func (r *ProductRepo) GetProductByCategoryID(uuid string) (*[]models.Product, error) {
	var products []models.Product
	if r.db.Where("categ_uuid = ?", uuid).Find(&products).RecordNotFound() {
		return nil, nil
	}

	return &products, nil
}

func (r *ProductRepo) GetProductByID(uuid string) (*models.Product, error) {
	var product models.Product
	if r.db.Where("uuid = ?", uuid).Find(&product).RecordNotFound() {
		return nil, errors.New("not found product")
	}

	return &product, nil
}

func (r *ProductRepo) CreateProduct(item *schema.ProductBodyParam) (*models.Product, error) {
	var product models.Product
	copier.Copy(&product, &item)

	if err := r.db.Create(&product).Error; err != nil {
		return nil, err
	}

	return &product, nil
}

func (r *ProductRepo) UpdateProduct(uuid string, item *schema.ProductBodyParam) (*models.Product, error) {
	var product models.Product
	if r.db.Where("uuid = ? ", uuid).First(&product).RecordNotFound() {
		return nil, errors.New("not found product")
	}

	copier.Copy(&product, &item)
	if err := r.db.Save(&product).Error; err != nil {
		return nil, err
	}

	return &product, nil
}
