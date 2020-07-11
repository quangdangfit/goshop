package repositories

import (
	"errors"

	"github.com/jinzhu/copier"
	"github.com/jinzhu/gorm"

	"goshop/dbs"
	"goshop/models"
)

type QuantityRepository interface {
	GetQuantities(map[string]interface{}) (*[]models.Quantity, error)
	GetQuantityByID(uuid string) (*models.Quantity, error)
	CreateQuantity(req *models.QuantityBodyRequest) (*models.Quantity, error)
	UpdateQuantity(uuid string, req *models.QuantityBodyRequest) (*models.Quantity, error)
}

type quantityRepo struct {
	db *gorm.DB
}

func NewQuantityRepository() QuantityRepository {
	return &quantityRepo{db: dbs.Database}
}

func (r *quantityRepo) GetQuantities(query map[string]interface{}) (*[]models.Quantity, error) {
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

func (r *quantityRepo) CreateQuantity(req *models.QuantityBodyRequest) (*models.Quantity, error) {
	var quantity models.Quantity
	copier.Copy(&quantity, &req)

	if err := r.db.Create(&quantity).Error; err != nil {
		return nil, err
	}

	return &quantity, nil
}

func (r *quantityRepo) UpdateQuantity(uuid string, req *models.QuantityBodyRequest) (*models.Quantity, error) {
	quantity, err := r.GetQuantityByID(uuid)
	if err != nil {
		return nil, err
	}

	copier.Copy(quantity, &req)
	if err := r.db.Save(&quantity).Error; err != nil {
		return nil, err
	}

	return quantity, nil
}
