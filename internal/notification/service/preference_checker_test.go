package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"goshop/internal/notification/model"
)

type stubLookup struct {
	id  string
	err error
}

func (s *stubLookup) GetUserIDByEmail(_ context.Context, _ string) (string, error) {
	return s.id, s.err
}

func TestDBPreferenceChecker_LookupErrorOptsIn(t *testing.T) {
	c := NewDBPreferenceChecker(&stubLookup{err: errors.New("boom")}, &stubPrefRepo{})
	enabled, err := c.IsEnabled(context.Background(), "x@example.com", "OrderPaid", "email")
	require.True(t, enabled)
	require.Error(t, err)
}

func TestDBPreferenceChecker_EmptyUserOptsIn(t *testing.T) {
	c := NewDBPreferenceChecker(&stubLookup{id: ""}, &stubPrefRepo{})
	enabled, err := c.IsEnabled(context.Background(), "x@example.com", "OrderPaid", "email")
	require.True(t, enabled)
	require.NoError(t, err)
}

func TestDBPreferenceChecker_PrefRepoErrorOptsIn(t *testing.T) {
	repo := &stubPrefRepo{
		getFn: func(_ context.Context, _, _, _ string) (*model.Preference, error) {
			return nil, errors.New("db down")
		},
	}
	c := NewDBPreferenceChecker(&stubLookup{id: "u1"}, repo)
	enabled, err := c.IsEnabled(context.Background(), "x@example.com", "OrderPaid", "email")
	require.True(t, enabled)
	require.Error(t, err)
}

func TestDBPreferenceChecker_NoPrefDefaultsEnabled(t *testing.T) {
	repo := &stubPrefRepo{
		getFn: func(_ context.Context, _, _, _ string) (*model.Preference, error) { return nil, nil },
	}
	c := NewDBPreferenceChecker(&stubLookup{id: "u1"}, repo)
	enabled, err := c.IsEnabled(context.Background(), "x@example.com", "OrderPaid", "email")
	require.True(t, enabled)
	require.NoError(t, err)
}

func TestDBPreferenceChecker_HonorsPrefValue(t *testing.T) {
	repo := &stubPrefRepo{
		getFn: func(_ context.Context, _, _, _ string) (*model.Preference, error) {
			return &model.Preference{Enabled: false}, nil
		},
	}
	c := NewDBPreferenceChecker(&stubLookup{id: "u1"}, repo)
	enabled, err := c.IsEnabled(context.Background(), "x@example.com", "OrderPaid", "email")
	require.False(t, enabled)
	require.NoError(t, err)
}
