-- Per-user opt-in/out per (event_type, channel). Absence of a row means "enabled" — see
-- internal/notification/service/preference_checker.go.

CREATE TABLE IF NOT EXISTS notification_preferences (
    id          VARCHAR(36)  PRIMARY KEY,
    user_id     VARCHAR(36)  NOT NULL,
    event_type  VARCHAR(64)  NOT NULL,
    channel     VARCHAR(32)  NOT NULL,
    enabled     BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_pref_user_event_channel UNIQUE (user_id, event_type, channel)
);
