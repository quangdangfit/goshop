package repositories

import (
	"errors"

	"github.com/jinzhu/copier"
	"gorm.io/gorm"

	"goshop/app/models"
	"goshop/app/serializers"
	"goshop/dbs"
)

type IWarehouseRepository interface {
	GetWarehouses(query *serializers.WarehouseQueryParam) (*[]models.Warehouse, error)
	GetWarehouseByID(uuid string) (*models.Warehouse, error)
	CreateWarehouse(item *serializers.WarehouseBodyParam) (*models.Warehouse, error)
	UpdateWarehouse(uuid string, item *serializers.WarehouseBodyParam) (*models.Warehouse, error)
}

type WarehouseRepo struct {
	db *gorm.DB
}

func NewWarehouseRepository() *WarehouseRepo {
	return &WarehouseRepo{db: dbs.Database}
}

func (w *WarehouseRepo) GetWarehouses(query *serializers.WarehouseQueryParam) (*[]models.Warehouse, error) {
	var warehouses []models.Warehouse
	if err := w.db.Find(&warehouses, query).Error; err != nil {
		return nil, nil
	}

	return &warehouses, nil
}

func (w *WarehouseRepo) GetWarehouseByID(uuid string) (*models.Warehouse, error) {
	var warehouse models.Warehouse
	if err := w.db.Where("uuid = ?", uuid).First(&warehouse).Error; err != nil {
		return nil, errors.New("not found warehouse")
	}

	return &warehouse, nil
}

func (w *WarehouseRepo) CreateWarehouse(item *serializers.WarehouseBodyParam) (*models.Warehouse, error) {
	var warehouse models.Warehouse
	copier.Copy(&warehouse, &item)

	if err := w.db.Create(&warehouse).Error; err != nil {
		return nil, err
	}

	return &warehouse, nil
}

func (w *WarehouseRepo) UpdateWarehouse(uuid string, item *serializers.WarehouseBodyParam) (*models.Warehouse, error) {
	warehouse, err := w.GetWarehouseByID(uuid)
	if err != nil {
		return nil, err
	}

	copier.Copy(warehouse, &item)
	if err := w.db.Save(&warehouse).Error; err != nil {
		return nil, err
	}

	return warehouse, nil
}
