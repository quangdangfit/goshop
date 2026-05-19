package service

import (
	"context"

	"github.com/quangdangfit/gocommon/logger"
	"github.com/quangdangfit/gocommon/validation"
	"golang.org/x/crypto/bcrypt"

	"goshop/internal/user/domain"
	"goshop/internal/user/model"
	"goshop/internal/user/repository"
	"goshop/pkg/apperror"
	"goshop/pkg/jtoken"
	"goshop/pkg/utils"
)

//go:generate mockery --name=UserService
type UserService interface {
	Login(ctx context.Context, req *domain.LoginReq) (*model.User, string, string, error)
	Register(ctx context.Context, req *domain.RegisterReq) (*model.User, error)
	GetUserByID(ctx context.Context, id string) (*model.User, error)
	RefreshToken(ctx context.Context, userID string) (string, error)
	ChangePassword(ctx context.Context, id string, req *domain.ChangePasswordReq) error
	UpsertUserFromOIDC(ctx context.Context, subject, email, name string) (*model.User, error)
}

type userService struct {
	validator validation.Validation
	repo      repository.UserRepository
}

func NewUserService(
	validator validation.Validation,
	repo repository.UserRepository) UserService {
	return &userService{
		validator: validator,
		repo:      repo,
	}
}

func (s *userService) Login(ctx context.Context, req *domain.LoginReq) (*model.User, string, string, error) {
	if err := s.validator.ValidateStruct(req); err != nil {
		return nil, "", "", err
	}

	user, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		logger.Errorf("Login.GetUserByEmail fail, email: %s, error: %s", req.Email, err)
		return nil, "", "", apperror.Wrap(apperror.ErrInvalidCredentials, err)
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		logger.Errorf("CompareHashAndPassword fail, email: %s, error: %s", req.Email, err)
		return nil, "", "", apperror.ErrInvalidCredentials
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

func (s *userService) Register(ctx context.Context, req *domain.RegisterReq) (*model.User, error) {
	if err := s.validator.ValidateStruct(req); err != nil {
		return nil, err
	}

	var user model.User
	if err := utils.Copy(&user, &req); err != nil {
		return nil, err
	}
	err := s.repo.Create(ctx, &user)
	if err != nil {
		logger.Errorf("Register.Create fail, email: %s, error: %s", req.Email, err)
		return nil, err
	}
	return &user, nil
}

func (s *userService) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	user, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		logger.Errorf("GetUserByID fail, id: %s, error: %s", id, err)
		return nil, err
	}

	return user, nil
}

func (s *userService) RefreshToken(ctx context.Context, userID string) (string, error) {
	user, err := s.repo.GetUserByID(ctx, userID)
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

func (s *userService) ChangePassword(ctx context.Context, id string, req *domain.ChangePasswordReq) error {
	if err := s.validator.ValidateStruct(req); err != nil {
		return err
	}
	user, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		logger.Errorf("ChangePassword.GetUserByID fail, id: %s, error: %s", id, err)
		return err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return apperror.ErrInvalidCredentials
	}

	user.Password = utils.HashAndSalt([]byte(req.NewPassword))
	err = s.repo.Update(ctx, user)
	if err != nil {
		logger.Errorf("ChangePassword.Update fail, id: %s, error: %s", id, err)
		return err
	}

	return nil
}

// UpsertUserFromOIDC creates or updates a local user record from the claims
// of an Authentik-issued ID token. The OIDC subject (sub) is used as the
// primary key so the link stays stable even if the user later changes their
// email address.
func (s *userService) UpsertUserFromOIDC(ctx context.Context, subject, email, name string) (*model.User, error) {
	// 1. Lookup by OIDC subject (our primary key for OIDC-provisioned users).
	user, err := s.repo.GetUserByID(ctx, subject)
	if err == nil {
		if user.Email != email {
			user.Email = email
			if err := s.repo.Update(ctx, user); err != nil {
				logger.Errorf("UpsertUserFromOIDC.Update fail, id: %s, error: %s", subject, err)
				return nil, err
			}
		}
		return user, nil
	}

	// 2. Fallback: a legacy JWT-era user with the same email exists. Keep the
	//    old ID and password hash so the user can still log in via either
	//    flow during the migration window.
	if userByEmail, err := s.repo.GetUserByEmail(ctx, email); err == nil {
		return userByEmail, nil
	}

	// 3. Brand-new user — provision with no password (Authentik owns credentials).
	newUser := &model.User{
		ID:    subject,
		Email: email,
		Role:  model.UserRoleCustomer,
	}
	if err := s.repo.Create(ctx, newUser); err != nil {
		logger.Errorf("UpsertUserFromOIDC.Create fail, email: %s, error: %s", email, err)
		return nil, err
	}
	return newUser, nil
}
