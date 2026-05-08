-- Stripe payment integration. `payments` is a one-to-one with orders (latest active
-- payment record); `provider_events` is the idempotency log used by the webhook handler.

CREATE TABLE IF NOT EXISTS payments (
    id                  VARCHAR(36)  PRIMARY KEY,
    order_id            VARCHAR(36)  NOT NULL,
    provider            VARCHAR(32)  NOT NULL,
    provider_intent_id  VARCHAR(128) NOT NULL,
    amount              BIGINT       NOT NULL,
    currency            VARCHAR(8)   NOT NULL,
    status              VARCHAR(32)  NOT NULL,
    created_at          TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_payments_order_id ON payments (order_id);
CREATE INDEX IF NOT EXISTS idx_payments_status   ON payments (status);

-- Composite PK: a provider may legitimately reuse event ids if you migrate providers,
-- but within one provider event ids are globally unique. (provider, event_id) prevents
-- duplicate side-effects on retried webhooks.
CREATE TABLE IF NOT EXISTS provider_events (
    provider    VARCHAR(32)  NOT NULL,
    event_id    VARCHAR(128) NOT NULL,
    received_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    PRIMARY KEY (provider, event_id)
);
