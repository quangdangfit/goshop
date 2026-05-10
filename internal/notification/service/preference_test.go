package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"goshop/internal/notification/model"
)

type stubPrefRepo struct {
	listFn   func(ctx context.Context, userID string) ([]*model.Preference, error)
	getFn    func(ctx context.Context, userID, eventType, channel string) (*model.Preference, error)
	upsertFn func(ctx context.Context, p *model.Preference) error
}

func (s *stubPrefRepo) ListByUser(ctx context.Context, u string) ([]*model.Preference, error) {
	return s.listFn(ctx, u)
}
func (s *stubPrefRepo) Get(ctx context.Context, u, e, c string) (*model.Preference, error) {
	return s.getFn(ctx, u, e, c)
}
func (s *stubPrefRepo) Upsert(ctx context.Context, p *model.Preference) error {
	return s.upsertFn(ctx, p)
}

func TestPreferenceServiceList(t *testing.T) {
	repo := &stubPrefRepo{
		listFn: func(_ context.Context, u string) ([]*model.Preference, error) {
			require.Equal(t, "u1", u)
			return []*model.Preference{{ID: "p1"}}, nil
		},
	}
	svc := NewPreferenceService(repo)
	got, err := svc.List(context.Background(), "u1")
	require.NoError(t, err)
	require.Len(t, got, 1)
}

func TestPreferenceServiceSet(t *testing.T) {
	called := false
	repo := &stubPrefRepo{
		upsertFn: func(_ context.Context, p *model.Preference) error {
			called = true
			require.Equal(t, "u1", p.UserID)
			require.Equal(t, "OrderPaid", p.EventType)
			require.Equal(t, "email", p.Channel)
			require.True(t, p.Enabled)
			return nil
		},
		getFn: func(_ context.Context, u, e, c string) (*model.Preference, error) {
			return &model.Preference{UserID: u, EventType: e, Channel: c, Enabled: true}, nil
		},
	}
	svc := NewPreferenceService(repo)
	got, err := svc.Set(context.Background(), "u1", "OrderPaid", "email", true)
	require.NoError(t, err)
	require.True(t, called)
	require.True(t, got.Enabled)
}

func TestPreferenceServiceSetUpsertError(t *testing.T) {
	repo := &stubPrefRepo{
		upsertFn: func(_ context.Context, _ *model.Preference) error {
			return errors.New("upsert failed")
		},
	}
	svc := NewPreferenceService(repo)
	_, err := svc.Set(context.Background(), "u1", "e", "c", false)
	require.Error(t, err)
}
