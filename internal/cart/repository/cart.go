package repository

import (
	"context"

	"goshop/internal/cart/model"
	"goshop/pkg/dbs"
)

//go:generate mockery --name=CartRepository
type CartRepository interface {
	Create(ctx context.Context, cart *model.Cart) error
	Update(ctx context.Context, cart *model.Cart) error
	GetCartByUserID(ctx context.Context, userID string) (*model.Cart, error)
}

type cartRepo struct {
	db dbs.Database
}

func NewCartRepository(db dbs.Database) CartRepository {
	return &cartRepo{db: db}
}

func (r *cartRepo) Create(ctx context.Context, cart *model.Cart) error {
	return r.db.Create(ctx, cart)
}

func (r *cartRepo) Update(ctx context.Context, cart *model.Cart) error {
	return r.db.Update(ctx, cart)
}

func (r *cartRepo) GetCartByUserID(ctx context.Context, userID string) (*model.Cart, error) {
	var order model.Cart
	opts := []dbs.FindOption{
		dbs.WithQuery(dbs.NewQuery("user_id = ?", userID)),
	}
	opts = append(opts, dbs.WithPreload([]string{"User", "Lines.Product"}))

	if err := r.db.FindOne(ctx, &order, opts...); err != nil {
		return nil, err
	}

	return &order, nil
}
