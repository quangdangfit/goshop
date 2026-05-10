package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	userModel "goshop/internal/user/model"
	userMocks "goshop/internal/user/repository/mocks"
)

func TestUserRepoLookup_ReturnsID(t *testing.T) {
	repo := userMocks.NewUserRepository(t)
	repo.On("GetUserByEmail", mock.Anything, "x@example.com").
		Return(&userModel.User{ID: "uid-1"}, nil).Once()

	id, err := NewUserRepoLookup(repo).GetUserIDByEmail(context.Background(), "x@example.com")
	require.NoError(t, err)
	require.Equal(t, "uid-1", id)
}

func TestUserRepoLookup_PropagatesError(t *testing.T) {
	repo := userMocks.NewUserRepository(t)
	repo.On("GetUserByEmail", mock.Anything, "x@example.com").
		Return(nil, errors.New("db down")).Once()

	id, err := NewUserRepoLookup(repo).GetUserIDByEmail(context.Background(), "x@example.com")
	require.Error(t, err)
	require.Empty(t, id)
}

func TestUserRepoLookup_NilUserReturnsEmpty(t *testing.T) {
	repo := userMocks.NewUserRepository(t)
	repo.On("GetUserByEmail", mock.Anything, "x@example.com").Return(nil, nil).Once()

	id, err := NewUserRepoLookup(repo).GetUserIDByEmail(context.Background(), "x@example.com")
	require.NoError(t, err)
	require.Empty(t, id)
}
