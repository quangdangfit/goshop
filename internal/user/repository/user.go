package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"goshop/internal/user/model"
	"goshop/pkg/config"
)

//go:generate mockery --name=IUserRepository
type IUserRepository interface {
	Create(ctx context.Context, user *model.User) error
	Update(ctx context.Context, user *model.User) error
	GetUserByID(ctx context.Context, id string) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
}

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepo {
	_ = db.AutoMigrate(&model.User{})
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(ctx context.Context, user *model.User) error {
	ctx, cancel := context.WithTimeout(ctx, config.DatabaseTimeout)
	defer cancel()

	if err := r.db.Create(&user).Error; err != nil {
		return err
	}

	return nil
}

func (r *UserRepo) Update(ctx context.Context, user *model.User) error {
	ctx, cancel := context.WithTimeout(ctx, config.DatabaseTimeout)
	defer cancel()

	if err := r.db.Save(&user).Error; err != nil {
		return err
	}

	return nil
}

func (r *UserRepo) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	ctx, cancel := context.WithTimeout(ctx, config.DatabaseTimeout)
	defer cancel()

	var user model.User
	if err := r.db.Where("id = ? ", id).First(&user).Error; err != nil {
		return nil, errors.New("user not found")
	}

	return &user, nil
}

func (r *UserRepo) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	ctx, cancel := context.WithTimeout(ctx, config.DatabaseTimeout)
	defer cancel()

	var user model.User
	if err := r.db.Where("email = ? ", email).First(&user).Error; err != nil {
		return nil, errors.New("user not found")
	}

	return &user, nil
}
