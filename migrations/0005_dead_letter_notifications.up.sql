-- Audit log for notifications that exhausted retry budget. Operators triage from this
-- table and can replay manually. Keep it append-only.

CREATE TABLE IF NOT EXISTS dead_letter_notifications (
    id          VARCHAR(36)  PRIMARY KEY,
    event_type  VARCHAR(64)  NOT NULL,
    user_email  VARCHAR(255) NOT NULL,
    payload     TEXT,
    last_error  TEXT,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_dlq_event_type ON dead_letter_notifications (event_type);
