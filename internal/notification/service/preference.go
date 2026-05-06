package service

import (
	"context"

	"goshop/internal/notification/model"
	"goshop/internal/notification/repository"
)

//go:generate mockery --name=PreferenceService
type PreferenceService interface {
	List(ctx context.Context, userID string) ([]*model.Preference, error)
	Set(ctx context.Context, userID, eventType, channel string, enabled bool) (*model.Preference, error)
}

type preferenceService struct {
	repo repository.PreferenceRepository
}

func NewPreferenceService(repo repository.PreferenceRepository) PreferenceService {
	return &preferenceService{repo: repo}
}

func (s *preferenceService) List(ctx context.Context, userID string) ([]*model.Preference, error) {
	return s.repo.ListByUser(ctx, userID)
}

func (s *preferenceService) Set(ctx context.Context, userID, eventType, channel string, enabled bool) (*model.Preference, error) {
	pref := &model.Preference{
		UserID:    userID,
		EventType: eventType,
		Channel:   channel,
		Enabled:   enabled,
	}
	if err := s.repo.Upsert(ctx, pref); err != nil {
		return nil, err
	}
	return s.repo.Get(ctx, userID, eventType, channel)
}
