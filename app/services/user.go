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
	"goshop/pkg/utils"
)

type IUserService interface {
	Login(ctx context.Context, req *serializers.LoginReq) (*models.User, string, string, error)
	Register(ctx context.Context, req *serializers.RegisterReq) (*models.User, error)
	GetUserByID(ctx context.Context, id string) (*models.User, error)
	RefreshToken(ctx context.Context, userID string) (string, error)
	ChangePassword(ctx context.Context, id string, req *serializers.ChangePasswordReq) error
}

type UserService struct {
	repo repositories.IUserRepository
}

func NewUserService(repo repositories.IUserRepository) IUserService {
	return &UserService{repo: repo}
}

func (u *UserService) Login(ctx context.Context, req *serializers.LoginReq) (*models.User, string, string, error) {
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
		"role":  user.Role,
	}
	accessToken := jtoken.GenerateAccessToken(tokenData)
	refreshToken := jtoken.GenerateRefreshToken(tokenData)
	return user, accessToken, refreshToken, nil
}

func (u *UserService) Register(ctx context.Context, req *serializers.RegisterReq) (*models.User, error) {
	var user models.User
	utils.Copy(&user, &req)
	err := u.repo.Create(ctx, &user)
	if err != nil {
		logger.Errorf("Register.Create fail, email: %s, error: %s", req.Email, err)
		return nil, err
	}
	return &user, nil
}

func (u *UserService) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	user, err := u.repo.GetUserByID(ctx, id)
	if err != nil {
		logger.Errorf("GetUserByID fail, id: %s, error: %s", id, err)
		return nil, err
	}

	return user, nil
}

func (u *UserService) RefreshToken(ctx context.Context, userID string) (string, error) {
	user, err := u.repo.GetUserByID(ctx, userID)
	if err != nil {
		logger.Errorf("RefreshToken.GetUserByID fail, id: %s, error: %s", userID, err)
		return "", err
	}

	tokenData := map[string]interface{}{
		"id":    user.ID,
		"email": user.Email,
		"role":  user.Role,
	}
	accessToken := jtoken.GenerateAccessToken(tokenData)
	return accessToken, nil
}

func (u *UserService) ChangePassword(ctx context.Context, id string, req *serializers.ChangePasswordReq) error {
	user, err := u.repo.GetUserByID(ctx, id)
	if err != nil {
		logger.Errorf("ChangePassword.GetUserByID fail, id: %s, error: %s", id, err)
		return err
	}

	passErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if passErr == bcrypt.ErrMismatchedHashAndPassword && passErr != nil {
		return errors.New("wrong password")
	}

	user.Password = utils.HashAndSalt([]byte(req.NewPassword))
	err = u.repo.Update(ctx, user)
	if err != nil {
		logger.Errorf("ChangePassword.Update fail, id: %s, error: %s", id, err)
		return err
	}

	return nil
}
