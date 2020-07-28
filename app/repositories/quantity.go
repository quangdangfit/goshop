package repositories

import (
	"errors"

	"github.com/jinzhu/copier"
	"github.com/jinzhu/gorm"

	"goshop/app/models"
	"goshop/app/schema"
	"goshop/dbs"
)

type QuantityRepository interface {
	GetQuantities(query *schema.QuantityQueryParam) (*[]models.Quantity, error)
	GetQuantityByID(uuid string) (*models.Quantity, error)
	GetQuantityProductID(productUUID string) (*models.Quantity, error)
	CreateQuantity(item *schema.QuantityBodyParam) (*models.Quantity, error)
	UpdateQuantity(uuid string, item *schema.QuantityBodyParam) (*models.Quantity, error)
}

type quantityRepo struct {
	db *gorm.DB
}

func NewQuantityRepository() QuantityRepository {
	return &quantityRepo{db: dbs.Database}
}

func (r *quantityRepo) GetQuantities(query *schema.QuantityQueryParam) (*[]models.Quantity, error) {
	var quantities []models.Quantity
	if r.db.Find(&quantities, query).RecordNotFound() {
		return nil, nil
	}

	return &quantities, nil
}

func (r *quantityRepo) GetQuantityByID(uuid string) (*models.Quantity, error) {
	var quantity models.Quantity
	if r.db.Where("uuid = ?", uuid).First(&quantity).RecordNotFound() {
		return nil, errors.New("not found quantity")
	}

	return &quantity, nil
}

func (r *quantityRepo) CreateQuantity(item *schema.QuantityBodyParam) (*models.Quantity, error) {
	var quantity models.Quantity
	copier.Copy(&quantity, &item)

	if err := r.db.Create(&quantity).Error; err != nil {
		return nil, err
	}

	return &quantity, nil
}

func (r *quantityRepo) UpdateQuantity(uuid string, item *schema.QuantityBodyParam) (*models.Quantity, error) {
	quantity, err := r.GetQuantityByID(uuid)
	if err != nil {
		return nil, err
	}

	copier.Copy(quantity, &item)
	if err := r.db.Save(&quantity).Error; err != nil {
		return nil, err
	}

	return quantity, nil
}

func (r *quantityRepo) GetQuantityProductID(productUUID string) (*models.Quantity, error) {
	var quantity models.Quantity
	if r.db.Where("product_uuid = ?", productUUID).First(&quantity).RecordNotFound() {
		return nil, errors.New("not found quantity")
	}

	return &quantity, nil
}
