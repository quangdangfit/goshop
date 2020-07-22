package repositories

import (
	"errors"

	"github.com/jinzhu/copier"
	"github.com/jinzhu/gorm"

	"goshop/app/models"
	"goshop/dbs"
)

type WarehouseRepository interface {
	GetWarehouses(map[string]interface{}) (*[]models.Warehouse, error)
	GetWarehouseByID(uuid string) (*models.Warehouse, error)
	CreateWarehouse(req *models.WarehouseBodyRequest) (*models.Warehouse, error)
	UpdateWarehouse(uuid string, req *models.WarehouseBodyRequest) (*models.Warehouse, error)
}

type warehouseRepo struct {
	db *gorm.DB
}

func NewWarehouseRepository() WarehouseRepository {
	return &warehouseRepo{db: dbs.Database}
}

func (r *warehouseRepo) GetWarehouses(query map[string]interface{}) (*[]models.Warehouse, error) {
	var warehouses []models.Warehouse
	if r.db.Find(&warehouses, query).RecordNotFound() {
		return nil, nil
	}

	return &warehouses, nil
}

func (r *warehouseRepo) GetWarehouseByID(uuid string) (*models.Warehouse, error) {
	var warehouse models.Warehouse
	if r.db.Where("uuid = ?", uuid).First(&warehouse).RecordNotFound() {
		return nil, errors.New("not found warehouse")
	}

	return &warehouse, nil
}

func (r *warehouseRepo) CreateWarehouse(req *models.WarehouseBodyRequest) (*models.Warehouse, error) {
	var warehouse models.Warehouse
	copier.Copy(&warehouse, &req)

	if err := r.db.Create(&warehouse).Error; err != nil {
		return nil, err
	}

	return &warehouse, nil
}

func (r *warehouseRepo) UpdateWarehouse(uuid string, req *models.WarehouseBodyRequest) (*models.Warehouse, error) {
	warehouse, err := r.GetWarehouseByID(uuid)
	if err != nil {
		return nil, err
	}

	copier.Copy(warehouse, &req)
	if err := r.db.Save(&warehouse).Error; err != nil {
		return nil, err
	}

	return warehouse, nil
}
