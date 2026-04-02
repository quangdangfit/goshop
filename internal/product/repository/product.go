package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"goshop/internal/product/domain"
	"goshop/internal/product/model"
	"goshop/pkg/config"
	"goshop/pkg/dbs"
	"goshop/pkg/paging"
)

//go:generate mockery --name=ProductRepository
type ProductRepository interface {
	Create(ctx context.Context, product *model.Product) error
	Update(ctx context.Context, product *model.Product) error
	ListProducts(ctx context.Context, req *domain.ListProductReq) ([]*model.Product, *paging.Pagination, error)
	GetProductByID(ctx context.Context, id string) (*model.Product, error)
	DecrementStock(ctx context.Context, id string, qty int) error
	UpdateRating(ctx context.Context, id string, avgRating float64, reviewCount int) error
}

type productRepo struct {
	db dbs.Database
}

func NewProductRepository(db dbs.Database) ProductRepository {
	return &productRepo{db: db}
}

func (r *productRepo) ListProducts(ctx context.Context, req *domain.ListProductReq) ([]*model.Product, *paging.Pagination, error) {
	ctx, cancel := context.WithTimeout(ctx, config.DatabaseTimeout)
	defer cancel()

	query := make([]dbs.Query, 0)
	if req.Name != "" {
		query = append(query, dbs.NewQuery("name LIKE ?", "%"+req.Name+"%"))
	}
	if req.Code != "" {
		query = append(query, dbs.NewQuery("code = ?", req.Code))
	}
	if req.CategoryID != "" {
		query = append(query, dbs.NewQuery("category_id = ?", req.CategoryID))
	}

	order := "created_at"
	if req.OrderBy != "" {
		order = req.OrderBy
		if req.OrderDesc {
			order += " DESC"
		}
	}

	var total int64
	if err := r.db.Count(ctx, &model.Product{}, &total, dbs.WithQuery(query...)); err != nil {
		return nil, nil, err
	}

	pagination := paging.New(req.Page, req.Limit, total)

	var products []*model.Product
	if err := r.db.Find(
		ctx,
		&products,
		dbs.WithQuery(query...),
		dbs.WithLimit(int(pagination.Limit)),
		dbs.WithOffset(int(pagination.Skip)),
		dbs.WithOrder(order),
	); err != nil {
		return nil, nil, err
	}

	return products, pagination, nil
}

func (r *productRepo) GetProductByID(ctx context.Context, id string) (*model.Product, error) {
	var product model.Product
	if err := r.db.FindById(ctx, id, &product); err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepo) Create(ctx context.Context, product *model.Product) error {
	return r.db.Create(ctx, product)
}

func (r *productRepo) Update(ctx context.Context, product *model.Product) error {
	return r.db.Update(ctx, product)
}

func (r *productRepo) DecrementStock(ctx context.Context, id string, qty int) error {
	result := r.db.GetDB().WithContext(ctx).Model(&model.Product{}).
		Where("id = ? AND stock_quantity >= ?", id, qty).
		UpdateColumn("stock_quantity", gorm.Expr("stock_quantity - ?", qty))
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("insufficient stock")
	}
	return nil
}

func (r *productRepo) UpdateRating(ctx context.Context, id string, avgRating float64, reviewCount int) error {
	return r.db.GetDB().WithContext(ctx).Model(&model.Product{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{"avg_rating": avgRating, "review_count": reviewCount}).Error
}
