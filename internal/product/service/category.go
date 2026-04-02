package service

import (
	"context"

	"github.com/quangdangfit/gocommon/validation"

	"goshop/internal/product/domain"
	"goshop/internal/product/model"
	"goshop/internal/product/repository"
	"goshop/pkg/utils"
)

//go:generate mockery --name=CategoryService
type CategoryService interface {
	ListCategories(ctx context.Context) ([]*model.Category, error)
	GetCategoryByID(ctx context.Context, id string) (*model.Category, error)
	Create(ctx context.Context, req *domain.CreateCategoryReq) (*model.Category, error)
	Update(ctx context.Context, id string, req *domain.UpdateCategoryReq) (*model.Category, error)
	Delete(ctx context.Context, id string) error
}

type categorySvc struct {
	validator validation.Validation
	repo      repository.CategoryRepository
}

func NewCategoryService(validator validation.Validation, repo repository.CategoryRepository) CategoryService {
	return &categorySvc{validator: validator, repo: repo}
}

func (s *categorySvc) ListCategories(ctx context.Context) ([]*model.Category, error) {
	return s.repo.List(ctx)
}

func (s *categorySvc) GetCategoryByID(ctx context.Context, id string) (*model.Category, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *categorySvc) Create(ctx context.Context, req *domain.CreateCategoryReq) (*model.Category, error) {
	if err := s.validator.ValidateStruct(req); err != nil {
		return nil, err
	}
	var category model.Category
	utils.Copy(&category, req)
	if err := s.repo.Create(ctx, &category); err != nil {
		return nil, err
	}
	return &category, nil
}

func (s *categorySvc) Update(ctx context.Context, id string, req *domain.UpdateCategoryReq) (*model.Category, error) {
	category, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	utils.Copy(category, req)
	if err := s.repo.Update(ctx, category); err != nil {
		return nil, err
	}
	return category, nil
}

func (s *categorySvc) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
