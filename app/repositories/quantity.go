package repositories

import (
	"errors"

	"github.com/jinzhu/copier"
	"gorm.io/gorm"

	"goshop/app/models"
	"goshop/app/serializers"
	"goshop/dbs"
)

type IQuantityRepository interface {
	GetQuantities(query *serializers.QuantityQueryParam) (*[]models.Quantity, error)
	GetQuantityByID(uuid string) (*models.Quantity, error)
	GetQuantityProductID(productID string) (*models.Quantity, error)
	CreateQuantity(item *serializers.QuantityBodyParam) (*models.Quantity, error)
	UpdateQuantity(uuid string, item *serializers.QuantityBodyParam) (*models.Quantity, error)
}

type QuantityRepo struct {
	db *gorm.DB
}

func NewQuantityRepository() *QuantityRepo {
	return &QuantityRepo{db: dbs.Database}
}

func (r *QuantityRepo) GetQuantities(query *serializers.QuantityQueryParam) (*[]models.Quantity, error) {
	var quantities []models.Quantity
	if err := r.db.Find(&quantities, query).Error; err != nil {
		return nil, err
	}

	return &quantities, nil
}

func (r *QuantityRepo) GetQuantityByID(uuid string) (*models.Quantity, error) {
	var quantity models.Quantity
	if err := r.db.Where("uuid = ?", uuid).First(&quantity).Error; err != nil {
		return nil, errors.New("not found quantity")
	}

	return &quantity, nil
}

func (r *QuantityRepo) CreateQuantity(item *serializers.QuantityBodyParam) (*models.Quantity, error) {
	var quantity models.Quantity
	copier.Copy(&quantity, &item)

	if err := r.db.Create(&quantity).Error; err != nil {
		return nil, err
	}

	return &quantity, nil
}

func (r *QuantityRepo) UpdateQuantity(uuid string, item *serializers.QuantityBodyParam) (*models.Quantity, error) {
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

func (r *QuantityRepo) GetQuantityProductID(productID string) (*models.Quantity, error) {
	var quantity models.Quantity
	if err := r.db.Where("product_id = ?", productID).First(&quantity).Error; err != nil {
		return nil, errors.New("not found quantity")
	}

	return &quantity, nil
}
