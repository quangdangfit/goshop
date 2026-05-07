package repository

import (
	"context"

	"goshop/internal/notification/model"
	"goshop/pkg/dbs"
	"goshop/pkg/notification"
)

// dbDeadLetterSink persists exhausted notifications to the dead_letter_notifications table.
// Implements notification.DeadLetterSink so the runtime notifier can be wrapped without the
// pkg/notification package needing to know about GORM.
type dbDeadLetterSink struct {
	db dbs.Database
}

func NewDeadLetterSink(db dbs.Database) notification.DeadLetterSink {
	return &dbDeadLetterSink{db: db}
}

func (s *dbDeadLetterSink) Record(ctx context.Context, eventType, userEmail, payload string, lastErr error) error {
	row := &model.DeadLetterNotification{
		EventType: eventType,
		UserEmail: userEmail,
		Payload:   payload,
	}
	if lastErr != nil {
		row.LastError = lastErr.Error()
	}
	return s.db.Create(ctx, row)
}
