package repository

import (
	"context"

	"goshop/internal/order/model"
	"goshop/pkg/dbs"
)

type UserRepository interface {
	GetUserByID(ctx context.Context, id string) (*model.User, error)
}

type userRepo struct {
	db dbs.Database
}

func NewUserRepository(db dbs.Database) UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	var user model.User
	if err := r.db.FindById(ctx, id, &user); err != nil {
		return nil, err
	}
	return &user, nil
}
