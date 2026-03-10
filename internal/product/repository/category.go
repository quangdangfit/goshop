package repository

import (
	"context"

	"goshop/internal/product/model"
	"goshop/pkg/dbs"
)

//go:generate mockery --name=CategoryRepository
type CategoryRepository interface {
	List(ctx context.Context) ([]*model.Category, error)
	GetByID(ctx context.Context, id string) (*model.Category, error)
	Create(ctx context.Context, category *model.Category) error
	Update(ctx context.Context, category *model.Category) error
	Delete(ctx context.Context, id string) error
}

type categoryRepo struct {
	db dbs.Database
}

func NewCategoryRepository(db dbs.Database) CategoryRepository {
	return &categoryRepo{db: db}
}

func (r *categoryRepo) List(ctx context.Context) ([]*model.Category, error) {
	var categories []*model.Category
	if err := r.db.Find(ctx, &categories); err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *categoryRepo) GetByID(ctx context.Context, id string) (*model.Category, error) {
	var category model.Category
	if err := r.db.FindById(ctx, id, &category); err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepo) Create(ctx context.Context, category *model.Category) error {
	return r.db.Create(ctx, category)
}

func (r *categoryRepo) Update(ctx context.Context, category *model.Category) error {
	return r.db.Update(ctx, category)
}

func (r *categoryRepo) Delete(ctx context.Context, id string) error {
	return r.db.Delete(ctx, &model.Category{}, dbs.WithQuery(dbs.NewQuery("id = ?", id)))
}
