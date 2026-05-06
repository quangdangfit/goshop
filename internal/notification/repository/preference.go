package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"goshop/internal/notification/model"
	"goshop/pkg/dbs"
)

//go:generate mockery --name=PreferenceRepository
type PreferenceRepository interface {
	ListByUser(ctx context.Context, userID string) ([]*model.Preference, error)
	Get(ctx context.Context, userID, eventType, channel string) (*model.Preference, error)
	Upsert(ctx context.Context, pref *model.Preference) error
}

type preferenceRepo struct {
	db dbs.Database
}

func NewPreferenceRepository(db dbs.Database) PreferenceRepository {
	return &preferenceRepo{db: db}
}

func (r *preferenceRepo) ListByUser(ctx context.Context, userID string) ([]*model.Preference, error) {
	var rows []*model.Preference
	err := r.db.GetDB().WithContext(ctx).Where("user_id = ?", userID).Find(&rows).Error
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (r *preferenceRepo) Get(ctx context.Context, userID, eventType, channel string) (*model.Preference, error) {
	var p model.Preference
	err := r.db.GetDB().WithContext(ctx).
		Where("user_id = ? AND event_type = ? AND channel = ?", userID, eventType, channel).
		First(&p).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}

// Upsert writes the preference, replacing any existing row for the same (user,event,channel).
func (r *preferenceRepo) Upsert(ctx context.Context, pref *model.Preference) error {
	existing, err := r.Get(ctx, pref.UserID, pref.EventType, pref.Channel)
	if err != nil {
		return err
	}
	if existing == nil {
		return r.db.Create(ctx, pref)
	}
	existing.Enabled = pref.Enabled
	return r.db.Update(ctx, existing)
}
