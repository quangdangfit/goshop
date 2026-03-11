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
	db := r.db.GetDB().WithContext(ctx)
	// Only upsert the CartLine rows (no nested Product), skip associations
	for _, line := range cart.Lines {
		line.CartID = cart.ID
		if err := db.Omit("Product").Save(line).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *cartRepo) DeleteLine(ctx context.Context, cartID, productID string) error {
	return r.db.GetDB().WithContext(ctx).
		Where("cart_id = ? AND product_id = ?", cartID, productID).
		Delete(&model.CartLine{}).Error
}

func (r *cartRepo) ClearCart(ctx context.Context, userID string) error {
	cart, err := r.GetCartByUserID(ctx, userID)
	if err != nil {
		return nil // cart doesn't exist, nothing to clear
	}
	return r.db.GetDB().WithContext(ctx).
		Where("cart_id = ?", cart.ID).
		Delete(&model.CartLine{}).Error
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
