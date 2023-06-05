package services

import (
	"context"

	"goshop/app/models"
	"goshop/app/repositories"
	"goshop/app/schema"
)

type IQuantityService interface {
	GetQuantities(ctx context.Context, query *schema.QuantityQueryParam) (*[]models.Quantity, error)
	GetQuantityByID(ctx context.Context, uuid string) (*models.Quantity, error)
	CreateQuantity(ctx context.Context, item *schema.QuantityBodyParam) (*models.Quantity, error)
	UpdateQuantity(ctx context.Context, uuid string, item *schema.QuantityBodyParam) (*models.Quantity, error)
}

type quantity struct {
	repo repositories.IQuantityRepository
}

func NewQuantityService(repo repositories.IQuantityRepository) IQuantityService {
	return &quantity{repo: repo}
}

func (q *quantity) GetQuantities(ctx context.Context, query *schema.QuantityQueryParam) (*[]models.Quantity, error) {
	quantities, err := q.repo.GetQuantities(query)
	if err != nil {
		return nil, err
	}

	return quantities, nil
}

func (q *quantity) GetQuantityByID(ctx context.Context, uuid string) (*models.Quantity, error) {
	quantity, err := q.repo.GetQuantityByID(uuid)
	if err != nil {
		return nil, err
	}

	return quantity, nil
}

func (q *quantity) CreateQuantity(ctx context.Context, item *schema.QuantityBodyParam) (*models.Quantity, error) {
	quantity, err := q.repo.CreateQuantity(item)
	if err != nil {
		return nil, err
	}

	return quantity, nil
}

func (q *quantity) UpdateQuantity(ctx context.Context, uuid string, item *schema.QuantityBodyParam) (*models.Quantity, error) {
	quantity, err := q.repo.UpdateQuantity(uuid, item)
	if err != nil {
		return nil, err
	}

	return quantity, nil
}
