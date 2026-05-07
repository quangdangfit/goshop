package service

import (
	"context"

	userRepo "goshop/internal/user/repository"
)

// userRepoLookup adapts internal/user/repository.UserRepository to the UserLookup contract
// expected by dbPreferenceChecker. Kept here so pkg/notification stays free of user-domain
// imports.
type userRepoLookup struct {
	users userRepo.UserRepository
}

func NewUserRepoLookup(users userRepo.UserRepository) UserLookup {
	return &userRepoLookup{users: users}
}

func (l *userRepoLookup) GetUserIDByEmail(ctx context.Context, email string) (string, error) {
	u, err := l.users.GetUserByEmail(ctx, email)
	if err != nil || u == nil {
		return "", err
	}
	return u.ID, nil
}
