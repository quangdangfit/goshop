package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"goshop/app/dbs"
	"goshop/config"
	"goshop/internal/user/model"
)

type IUserRepository interface {
	Create(ctx context.Context, user *model.User) error
	Update(ctx context.Context, user *model.User) error
	GetUserByID(ctx context.Context, id string) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
}

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepository() *UserRepo {
	return &UserRepo{db: dbs.Database}
}

func (u *UserRepo) Create(ctx context.Context, user *model.User) error {
	ctx, cancel := context.WithTimeout(ctx, config.DatabaseTimeout)
	defer cancel()

	if err := u.db.Create(&user).Error; err != nil {
		return err
	}

	return nil
}

func (u *UserRepo) Update(ctx context.Context, user *model.User) error {
	ctx, cancel := context.WithTimeout(ctx, config.DatabaseTimeout)
	defer cancel()

	if err := u.db.Save(&user).Error; err != nil {
		return err
	}

	return nil
}

func (u *UserRepo) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	ctx, cancel := context.WithTimeout(ctx, config.DatabaseTimeout)
	defer cancel()

	var user model.User
	if err := dbs.Database.Where("id = ? ", id).First(&user).Error; err != nil {
		return nil, errors.New("user not found")
	}

	return &user, nil
}

func (u *UserRepo) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	ctx, cancel := context.WithTimeout(ctx, config.DatabaseTimeout)
	defer cancel()

	var user model.User
	if err := dbs.Database.Where("email = ? ", email).First(&user).Error; err != nil {
		return nil, errors.New("user not found")
	}

	return &user, nil
}
