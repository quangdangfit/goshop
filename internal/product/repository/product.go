package repository

import (
	"context"

	"goshop/internal/product/dto"
	"goshop/internal/product/model"
	"goshop/pkg/config"
	"goshop/pkg/dbs"
	"goshop/pkg/paging"
)

//go:generate mockery --name=IProductRepository
type IProductRepository interface {
	Create(ctx context.Context, product *model.Product) error
	Update(ctx context.Context, product *model.Product) error
	ListProducts(ctx context.Context, req *dto.ListProductReq) ([]*model.Product, *paging.Pagination, error)
	GetProductByID(ctx context.Context, id string) (*model.Product, error)
}

type ProductRepo struct {
	db dbs.IDatabase
}

func NewProductRepository(db dbs.IDatabase) *ProductRepo {
	return &ProductRepo{db: db}
}

func (r *ProductRepo) ListProducts(ctx context.Context, req *dto.ListProductReq) ([]*model.Product, *paging.Pagination, error) {
	ctx, cancel := context.WithTimeout(ctx, config.DatabaseTimeout)
	defer cancel()

	query := make([]dbs.Query, 0)
	if req.Name != "" {
		query = append(query, dbs.NewQuery("name LIKE ?", "%"+req.Name+"%"))
	}
	if req.Code != "" {
		query = append(query, dbs.NewQuery("code = ?", req.Code))
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

func (r *ProductRepo) GetProductByID(ctx context.Context, id string) (*model.Product, error) {
	var product model.Product
	if err := r.db.FindById(ctx, id, &product); err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *ProductRepo) Create(ctx context.Context, product *model.Product) error {
	return r.db.Create(ctx, product)
}

func (r *ProductRepo) Update(ctx context.Context, product *model.Product) error {
	return r.db.Update(ctx, product)
}
