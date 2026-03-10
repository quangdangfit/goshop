package repository

import (
	"context"

	"goshop/internal/user/model"
	"goshop/pkg/dbs"
)

//go:generate mockery --name=WishlistRepository
type WishlistRepository interface {
	GetWishlist(ctx context.Context, userID string) ([]*model.Wishlist, error)
	Add(ctx context.Context, userID, productID string) error
	Remove(ctx context.Context, userID, productID string) error
}

type wishlistRepo struct {
	db dbs.Database
}

func NewWishlistRepository(db dbs.Database) WishlistRepository {
	return &wishlistRepo{db: db}
}

func (r *wishlistRepo) GetWishlist(ctx context.Context, userID string) ([]*model.Wishlist, error) {
	var items []*model.Wishlist
	if err := r.db.Find(ctx, &items, dbs.WithQuery(dbs.NewQuery("user_id = ?", userID))); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *wishlistRepo) Add(ctx context.Context, userID, productID string) error {
	item := &model.Wishlist{UserID: userID, ProductID: productID}
	return r.db.Create(ctx, item)
}

func (r *wishlistRepo) Remove(ctx context.Context, userID, productID string) error {
	return r.db.Delete(ctx, &model.Wishlist{},
		dbs.WithQuery(dbs.NewQuery("user_id = ? AND product_id = ?", userID, productID)),
	)
}
