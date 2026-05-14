package grpc

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"goshop/internal/user/model"
)

func TestUserInfoFromModel_Nil(t *testing.T) {
	require.Nil(t, userInfoFromModel(nil))
}

func TestUserInfoFromModel_MapsFields(t *testing.T) {
	now := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	got := userInfoFromModel(&model.User{
		ID:        "u1",
		Email:     "a@b.c",
		CreatedAt: now,
		UpdatedAt: now,
	})
	require.NotNil(t, got)
	require.Equal(t, "u1", got.Id)
	require.Equal(t, "a@b.c", got.Email)
	require.Equal(t, now.Format(time.RFC3339), got.CreatedAt)
	require.Equal(t, now.Format(time.RFC3339), got.UpdatedAt)
}
