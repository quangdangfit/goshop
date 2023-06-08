package services

import (
	"context"
	"errors"

	"github.com/quangdangfit/gocommon/logger"
	"golang.org/x/crypto/bcrypt"

	"goshop/app/models"
	"goshop/app/repositories"
	"goshop/app/serializers"
	"goshop/pkg/jtoken"
)

type IUserService interface {
	Login(ctx context.Context, req *serializers.LoginReq) (*models.User, string, string, error)
	Register(ctx context.Context, req *serializers.RegisterReq) (*models.User, error)
	GetUserByID(ctx context.Context, id string) (*models.User, error)
	RefreshToken(ctx context.Context, userID string) (string, error)
}

type user struct {
	repo repositories.IUserRepository
}

func NewUserService(repo repositories.IUserRepository) IUserService {
	return &user{repo: repo}
}

func (u *user) Login(ctx context.Context, req *serializers.LoginReq) (*models.User, string, string, error) {
	user, err := u.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		logger.Errorf("Login.GetUserByEmail fail, email: %s, error: %s", req.Email, err)
		return nil, "", "", err
	}

	if user == nil {
		return nil, "", "", errors.New("user not found")
	}

	passErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if passErr == bcrypt.ErrMismatchedHashAndPassword && passErr != nil {
		return nil, "", "", errors.New("wrong password")
	}

	tokenData := map[string]interface{}{
		"id":    user.ID,
		"email": user.Email,
	}
	accessToken := jtoken.GenerateAccessToken(tokenData)
	refreshToken := jtoken.GenerateRefreshToken(tokenData)
	return user, accessToken, refreshToken, nil
}

func (u *user) Register(ctx context.Context, req *serializers.RegisterReq) (*models.User, error) {
	user, err := u.repo.Create(ctx, req)
	if err != nil {
		logger.Errorf("Register.Create fail, email: %s, error: %s", req.Email, err)
		return nil, err
	}
	return user, nil
}

func (u *user) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	user, err := u.repo.GetUserByID(ctx, id)
	if err != nil {
		logger.Errorf("GetUserByID fail, id: %s, error: %s", id, err)
		return nil, err
	}

	return user, nil
}

func (u *user) RefreshToken(ctx context.Context, userID string) (string, error) {
	user, err := u.repo.GetUserByID(ctx, userID)
	if err != nil {
		logger.Errorf("RefreshToken.GetUserByID fail, id: %s, error: %s", userID, err)
		return "", err
	}

	tokenData := map[string]interface{}{
		"id":    user.ID,
		"email": user.Email,
	}
	accessToken := jtoken.GenerateAccessToken(tokenData)
	return accessToken, nil
}
