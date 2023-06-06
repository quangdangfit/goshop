package services

import (
	"context"

	"github.com/quangdangfit/gocommon/logger"

	"goshop/app/models"
	"goshop/app/repositories"
	"goshop/app/serializers"
)

type ICategoryService interface {
	GetCategories(ctx context.Context, query *serializers.CategoryQueryParam) (*[]models.Category, error)
	GetCategoryByID(ctx context.Context, uuid string) (*models.Category, error)
	CreateCategory(cxt context.Context, item *serializers.Category) (*models.Category, error)
	UpdateCategory(ctx context.Context, uuid string, item *serializers.CategoryBodyParam) (*models.Category, error)
}

type category struct {
	repo repositories.ICategoryRepository
}

func NewCategoryService(repo repositories.ICategoryRepository) ICategoryService {
	return &category{repo: repo}
}

func (c *category) GetCategories(ctx context.Context, query *serializers.CategoryQueryParam) (*[]models.Category, error) {
	categories, err := c.repo.GetCategories(query)
	if err != nil {
		logger.Error("Failed to get categories: ", err)
		return nil, err
	}

	return categories, nil
}

func (c *category) GetCategoryByID(ctx context.Context, uuid string) (*models.Category, error) {
	category, err := c.repo.GetCategoryByID(uuid)
	if err != nil {
		logger.Error("Failed to get category: ", err)
		return nil, err
	}

	return category, nil
}

func (c *category) CreateCategory(cxt context.Context, item *serializers.Category) (*models.Category, error) {
	category, err := c.repo.CreateCategory(item)
	if err != nil {
		logger.Error("Failed to create category", err.Error())
		return nil, err
	}

	return category, nil
}

func (c *category) UpdateCategory(ctx context.Context, uuid string, item *serializers.CategoryBodyParam) (*models.Category, error) {
	category, err := c.repo.UpdateCategory(uuid, item)
	if err != nil {
		logger.Error("Failed to update category: ", err)
		return nil, err
	}

	return category, nil
}
