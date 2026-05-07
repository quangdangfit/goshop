package service

import (
	"context"

	"goshop/internal/notification/repository"
	"goshop/pkg/notification"
)

// UserLookup resolves a user email to a user ID. Decoupled from internal/user to avoid
// importing that package directly (and to keep the adapter testable).
type UserLookup interface {
	GetUserIDByEmail(ctx context.Context, email string) (string, error)
}

// dbPreferenceChecker is the production PreferenceChecker: looks up the user by email,
// then queries the notification.Preference table for the (event,channel) row.
//
// Default policy: opt-in. Missing user, missing preference row, and any DB error fall back
// to "enabled=true" so a transient outage never silences delivery.
type dbPreferenceChecker struct {
	users UserLookup
	prefs repository.PreferenceRepository
}

func NewDBPreferenceChecker(users UserLookup, prefs repository.PreferenceRepository) notification.PreferenceChecker {
	return &dbPreferenceChecker{users: users, prefs: prefs}
}

func (c *dbPreferenceChecker) IsEnabled(ctx context.Context, userEmail, eventType, channel string) (bool, error) {
	userID, err := c.users.GetUserIDByEmail(ctx, userEmail)
	if err != nil || userID == "" {
		return true, err
	}
	pref, err := c.prefs.Get(ctx, userID, eventType, channel)
	if err != nil {
		return true, err
	}
	if pref == nil {
		return true, nil
	}
	return pref.Enabled, nil
}
