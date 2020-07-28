package services

import (
	"context"

	jwtMiddle "goshop/app/middleware/jwt"
	"goshop/app/models"
	"goshop/app/repositories"
	"goshop/app/schema"
)

type IUserService interface {
	Login(ctx context.Context, item *schema.Login) (*models.User, string, error)
	Register(ctx context.Context, item *schema.Register) (*models.User, string, error)
	GetUserByID(ctx context.Context, uuid string) (*models.User, error)
}

type user struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) IUserService {
	return &user{repo: repo}
}

func (u *user) checkPermission(uuid string, data map[string]interface{}) bool {
	return data["uuid"] == uuid
}

func (u *user) Login(ctx context.Context, item *schema.Login) (*models.User, string, error) {
	user, err := u.repo.Login(item)
	if err != nil {
		return nil, "", err
	}

	token := jwtMiddle.GenerateToken(user)
	return user, token, nil
}

func (u *user) Register(ctx context.Context, item *schema.Register) (*models.User, string, error) {
	user, err := u.repo.Register(item)
	if err != nil {
		return nil, "", err
	}

	token := jwtMiddle.GenerateToken(user)
	return user, token, nil
}

func (u *user) GetUserByID(ctx context.Context, uuid string) (*models.User, error) {
	user, err := u.repo.GetUserByID(uuid)
	if err != nil {
		return nil, err
	}

	return user, nil
}
