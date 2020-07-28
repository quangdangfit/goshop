package repositories

import (
	"errors"

	"github.com/jinzhu/copier"
	"github.com/jinzhu/gorm"

	"goshop/app/models"
	"goshop/app/schema"
	"goshop/dbs"
)

type WarehouseRepository interface {
	GetWarehouses(query *schema.WarehouseQueryParam) (*[]models.Warehouse, error)
	GetWarehouseByID(uuid string) (*models.Warehouse, error)
	CreateWarehouse(item *schema.WarehouseBodyParam) (*models.Warehouse, error)
	UpdateWarehouse(uuid string, item *schema.WarehouseBodyParam) (*models.Warehouse, error)
}

type warehouseRepo struct {
	db *gorm.DB
}

func NewWarehouseRepository() WarehouseRepository {
	return &warehouseRepo{db: dbs.Database}
}

func (w *warehouseRepo) GetWarehouses(query *schema.WarehouseQueryParam) (*[]models.Warehouse, error) {
	var warehouses []models.Warehouse
	if w.db.Find(&warehouses, query).RecordNotFound() {
		return nil, nil
	}

	return &warehouses, nil
}

func (w *warehouseRepo) GetWarehouseByID(uuid string) (*models.Warehouse, error) {
	var warehouse models.Warehouse
	if w.db.Where("uuid = ?", uuid).First(&warehouse).RecordNotFound() {
		return nil, errors.New("not found warehouse")
	}

	return &warehouse, nil
}

func (w *warehouseRepo) CreateWarehouse(item *schema.WarehouseBodyParam) (*models.Warehouse, error) {
	var warehouse models.Warehouse
	copier.Copy(&warehouse, &item)

	if err := w.db.Create(&warehouse).Error; err != nil {
		return nil, err
	}

	return &warehouse, nil
}

func (w *warehouseRepo) UpdateWarehouse(uuid string, item *schema.WarehouseBodyParam) (*models.Warehouse, error) {
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
