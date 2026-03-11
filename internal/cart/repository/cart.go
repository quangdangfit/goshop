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
	DeleteLine(ctx context.Context, cartID, productID string) error
	ClearCart(ctx context.Context, userID string) error
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
	for _, line := range cart.Lines {
		line.CartID = cart.ID
		// Nil out Product to avoid GORM cascade-saving the nested product
		product := line.Product
		line.Product = nil
		if err := r.db.Update(ctx, line); err != nil {
			line.Product = product
			return err
		}
		line.Product = product
	}
	return nil
}

func (r *cartRepo) DeleteLine(ctx context.Context, cartID, productID string) error {
	return r.db.Delete(ctx, &model.CartLine{},
		dbs.WithQuery(dbs.NewQuery("cart_id = ? AND product_id = ?", cartID, productID)))
}

func (r *cartRepo) ClearCart(ctx context.Context, userID string) error {
	cart, err := r.GetCartByUserID(ctx, userID)
	if err != nil {
		return nil // cart doesn't exist, nothing to clear
	}
	return r.db.Delete(ctx, &model.CartLine{},
		dbs.WithQuery(dbs.NewQuery("cart_id = ?", cart.ID)))
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
