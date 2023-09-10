package repository

import (
	"context"

	"goshop/internal/cart/model"
	"goshop/pkg/dbs"
)

//go:generate mockery --name=ICartRepository
type ICartRepository interface {
	Create(ctx context.Context, cart *model.Cart) error
	Update(ctx context.Context, cart *model.Cart) error
	GetCartByUserID(ctx context.Context, userID string) (*model.Cart, error)
}

type CartRepo struct {
	db dbs.IDatabase
}

func NewCartRepository(db dbs.IDatabase) *CartRepo {
	return &CartRepo{db: db}
}

func (r *CartRepo) Create(ctx context.Context, cart *model.Cart) error {
	return r.db.Create(ctx, cart)
}

func (r *CartRepo) Update(ctx context.Context, cart *model.Cart) error {
	return r.db.Update(ctx, cart)
}

func (r *CartRepo) GetCartByUserID(ctx context.Context, userID string) (*model.Cart, error) {
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
