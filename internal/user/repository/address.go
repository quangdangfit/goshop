package repository

import (
	"context"

	"goshop/internal/user/model"
	"goshop/pkg/dbs"
)

//go:generate mockery --name=AddressRepository
type AddressRepository interface {
	ListByUser(ctx context.Context, userID string) ([]*model.Address, error)
	GetByID(ctx context.Context, id, userID string) (*model.Address, error)
	Create(ctx context.Context, address *model.Address) error
	Update(ctx context.Context, address *model.Address) error
	Delete(ctx context.Context, id, userID string) error
	SetDefault(ctx context.Context, id, userID string) error
}

type addressRepo struct {
	db dbs.Database
}

func NewAddressRepository(db dbs.Database) AddressRepository {
	return &addressRepo{db: db}
}

func (r *addressRepo) ListByUser(ctx context.Context, userID string) ([]*model.Address, error) {
	var addresses []*model.Address
	if err := r.db.Find(ctx, &addresses, dbs.WithQuery(dbs.NewQuery("user_id = ?", userID))); err != nil {
		return nil, err
	}
	return addresses, nil
}

func (r *addressRepo) GetByID(ctx context.Context, id, userID string) (*model.Address, error) {
	var address model.Address
	if err := r.db.FindOne(ctx, &address,
		dbs.WithQuery(dbs.NewQuery("id = ? AND user_id = ?", id, userID)),
	); err != nil {
		return nil, err
	}
	return &address, nil
}

func (r *addressRepo) Create(ctx context.Context, address *model.Address) error {
	return r.db.Create(ctx, address)
}

func (r *addressRepo) Update(ctx context.Context, address *model.Address) error {
	return r.db.Update(ctx, address)
}

func (r *addressRepo) Delete(ctx context.Context, id, userID string) error {
	if _, err := r.GetByID(ctx, id, userID); err != nil {
		return err
	}
	return r.db.Delete(ctx, &model.Address{}, dbs.WithQuery(dbs.NewQuery("id = ? AND user_id = ?", id, userID)))
}

func (r *addressRepo) SetDefault(ctx context.Context, id, userID string) error {
	address, err := r.GetByID(ctx, id, userID)
	if err != nil {
		return err
	}

	return r.db.WithTransaction(func() error {
		if err := r.db.GetDB().Model(&model.Address{}).
			Where("user_id = ?", userID).
			Update("is_default", false).Error; err != nil {
			return err
		}
		address.IsDefault = true
		return r.db.Update(ctx, address)
	})
}
