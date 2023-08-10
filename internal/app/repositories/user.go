package repositories

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"goshop/internal/app/dbs"
	"goshop/internal/app/models"
	"goshop/internal/config"
)

type IUserRepository interface {
	Create(ctx context.Context, user *models.User) error
	Update(ctx context.Context, user *models.User) error
	GetUserByID(ctx context.Context, id string) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
}

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepository() *UserRepo {
	return &UserRepo{db: dbs.Database}
}

func (u *UserRepo) Create(ctx context.Context, user *models.User) error {
	ctx, cancel := context.WithTimeout(ctx, config.DatabaseTimeout)
	defer cancel()

	if err := u.db.Create(&user).Error; err != nil {
		return err
	}

	return nil
}

func (u *UserRepo) Update(ctx context.Context, user *models.User) error {
	ctx, cancel := context.WithTimeout(ctx, config.DatabaseTimeout)
	defer cancel()

	if err := u.db.Save(&user).Error; err != nil {
		return err
	}

	return nil
}

func (u *UserRepo) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, config.DatabaseTimeout)
	defer cancel()

	var user models.User
	if err := dbs.Database.Where("id = ? ", id).First(&user).Error; err != nil {
		return nil, errors.New("user not found")
	}

	return &user, nil
}

func (u *UserRepo) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, config.DatabaseTimeout)
	defer cancel()

	var user models.User
	if err := dbs.Database.Where("email = ? ", email).First(&user).Error; err != nil {
		return nil, errors.New("user not found")
	}

	return &user, nil
}
