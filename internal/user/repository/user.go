package repository

import (
	"context"

	"goshop/internal/user/model"
	"goshop/pkg/dbs"
)

//go:generate mockery --name=IUserRepository
type IUserRepository interface {
	Create(ctx context.Context, user *model.User) error
	Update(ctx context.Context, user *model.User) error
	GetUserByID(ctx context.Context, id string) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
}

type UserRepo struct {
	db dbs.IDatabase
}

func NewUserRepository(db dbs.IDatabase) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(ctx context.Context, user *model.User) error {
	return r.db.Create(ctx, user)
}

func (r *UserRepo) Update(ctx context.Context, user *model.User) error {
	return r.db.Update(ctx, user)
}

func (r *UserRepo) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	var user model.User
	if err := r.db.FindById(ctx, id, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepo) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	query := dbs.NewQuery("email = ?", email)
	if err := r.db.FindOne(ctx, &user, dbs.WithQuery(query)); err != nil {
		return nil, err
	}

	return &user, nil
}
