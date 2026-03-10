package repository

import (
	"context"

	"goshop/internal/product/model"
	"goshop/pkg/dbs"
	"goshop/pkg/paging"
)

//go:generate mockery --name=ReviewRepository
type ReviewRepository interface {
	ListByProduct(ctx context.Context, productID string, page, limit int64) ([]*model.Review, *paging.Pagination, error)
	GetByID(ctx context.Context, id string) (*model.Review, error)
	Create(ctx context.Context, review *model.Review) error
	Update(ctx context.Context, review *model.Review) error
	Delete(ctx context.Context, id, userID string) error
	GetAggregates(ctx context.Context, productID string) (avgRating float64, count int, err error)
}

type reviewRepo struct {
	db dbs.Database
}

func NewReviewRepository(db dbs.Database) ReviewRepository {
	return &reviewRepo{db: db}
}

func (r *reviewRepo) ListByProduct(ctx context.Context, productID string, page, limit int64) ([]*model.Review, *paging.Pagination, error) {
	var total int64
	if err := r.db.Count(ctx, &model.Review{}, &total, dbs.WithQuery(dbs.NewQuery("product_id = ?", productID))); err != nil {
		return nil, nil, err
	}
	pagination := paging.New(page, limit, total)
	var reviews []*model.Review
	if err := r.db.Find(ctx, &reviews,
		dbs.WithQuery(dbs.NewQuery("product_id = ?", productID)),
		dbs.WithLimit(int(pagination.Limit)),
		dbs.WithOffset(int(pagination.Skip)),
		dbs.WithOrder("created_at DESC"),
	); err != nil {
		return nil, nil, err
	}
	return reviews, pagination, nil
}

func (r *reviewRepo) GetByID(ctx context.Context, id string) (*model.Review, error) {
	var review model.Review
	if err := r.db.FindById(ctx, id, &review); err != nil {
		return nil, err
	}
	return &review, nil
}

func (r *reviewRepo) Create(ctx context.Context, review *model.Review) error {
	return r.db.Create(ctx, review)
}

func (r *reviewRepo) Update(ctx context.Context, review *model.Review) error {
	return r.db.Update(ctx, review)
}

func (r *reviewRepo) Delete(ctx context.Context, id, userID string) error {
	return r.db.Delete(ctx, &model.Review{},
		dbs.WithQuery(dbs.NewQuery("id = ? AND user_id = ?", id, userID)),
	)
}

type reviewAgg struct {
	Avg   float64
	Count int
}

func (r *reviewRepo) GetAggregates(ctx context.Context, productID string) (float64, int, error) {
	var result reviewAgg
	err := r.db.GetDB().WithContext(ctx).Model(&model.Review{}).
		Select("COALESCE(AVG(rating), 0) as avg, COUNT(*) as count").
		Where("product_id = ? AND deleted_at IS NULL", productID).
		Scan(&result).Error
	return result.Avg, result.Count, err
}
