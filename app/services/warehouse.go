package services

import (
	"context"

	"goshop/app/models"
	"goshop/app/repositories"
	"goshop/app/schema"
)

type IWarehouseSerivce interface {
	GetWarehouses(ctx context.Context, query *schema.WarehouseQueryParam) (*[]models.Warehouse, error)
	GetWarehouseByID(ctx context.Context, uuid string) (*models.Warehouse, error)
	CreateWarehouse(ctx context.Context, item *schema.WarehouseBodyParam) (*models.Warehouse, error)
	UpdateWarehouse(ctx context.Context, uuid string, items *schema.WarehouseBodyParam) (*models.Warehouse, error)
}

type warehouse struct {
	repo repositories.WarehouseRepository
}

func NewWarehouseService(repo repositories.WarehouseRepository) IWarehouseSerivce {
	return &warehouse{repo: repo}
}

func (w *warehouse) GetWarehouses(ctx context.Context, query *schema.WarehouseQueryParam) (*[]models.Warehouse, error) {
	warehouses, err := w.repo.GetWarehouses(query)
	if err != nil {
		return nil, err
	}

	return warehouses, nil
}

func (w *warehouse) GetWarehouseByID(ctx context.Context, uuid string) (*models.Warehouse, error) {
	warehouse, err := w.repo.GetWarehouseByID(uuid)
	if err != nil {
		return nil, err
	}

	return warehouse, nil
}

func (w *warehouse) CreateWarehouse(ctx context.Context, item *schema.WarehouseBodyParam) (*models.Warehouse, error) {
	warehouse, err := w.repo.CreateWarehouse(item)
	if err != nil {
		return nil, err
	}

	return warehouse, nil
}

func (w *warehouse) UpdateWarehouse(ctx context.Context, uuid string, item *schema.WarehouseBodyParam) (*models.Warehouse, error) {
	warehouse, err := w.repo.UpdateWarehouse(uuid, item)
	if err != nil {
		return nil, err
	}

	return warehouse, nil
}
