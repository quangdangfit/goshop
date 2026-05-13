-- Consolidated initial migration for the Stripe + reservations + notifications work.
-- Applies cleanly on a fresh database after AutoMigrate has created the base tables
-- (users, products, orders, order_lines, categories, etc.). Each statement is
-- idempotent (`IF NOT EXISTS` / `IF EXISTS`) so this file is safe to re-run.

-- ============================================================
-- 1. Drop server-side cart (cart now lives client-side, localStorage).
--    Run only after the deprecated /cart gRPC service is no longer registered.
-- ============================================================
DROP TABLE IF EXISTS cart_lines;
DROP TABLE IF EXISTS carts;

-- ============================================================
-- 2. Stripe payments + provider event idempotency log.
--    `payments` is one-to-one with orders (latest active record);
--    `provider_events` PK (provider, event_id) prevents duplicate webhook side effects.
-- ============================================================
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

CREATE TABLE IF NOT EXISTS provider_events (
    provider    VARCHAR(32)  NOT NULL,
    event_id    VARCHAR(128) NOT NULL,
    received_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    PRIMARY KEY (provider, event_id)
);

-- ============================================================
-- 3. Strict inventory: reserved-but-not-committed counters + CHECK guard.
--    The CHECK constraints are the safety net — even with a buggy app the DB
--    rejects writes that would drive a counter negative or oversell.
-- ============================================================
ALTER TABLE products
    ADD COLUMN IF NOT EXISTS reserved_quantity INT NOT NULL DEFAULT 0;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'chk_products_stock_nonneg') THEN
        ALTER TABLE products ADD CONSTRAINT chk_products_stock_nonneg       CHECK (stock_quantity    >= 0);
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'chk_products_reserved_nonneg') THEN
        ALTER TABLE products ADD CONSTRAINT chk_products_reserved_nonneg    CHECK (reserved_quantity >= 0);
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'chk_products_reserved_lte_stock') THEN
        ALTER TABLE products ADD CONSTRAINT chk_products_reserved_lte_stock CHECK (reserved_quantity <= stock_quantity);
    END IF;
END $$;

CREATE TABLE IF NOT EXISTS stock_reservations (
    id          VARCHAR(36)  PRIMARY KEY,
    order_id    VARCHAR(36)  NOT NULL,
    product_id  VARCHAR(36)  NOT NULL,
    quantity    INT          NOT NULL CHECK (quantity > 0),
    status      VARCHAR(16)  NOT NULL,
    expires_at  TIMESTAMPTZ  NOT NULL,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_reservations_order      ON stock_reservations (order_id);
CREATE INDEX IF NOT EXISTS idx_reservations_expires_at ON stock_reservations (expires_at) WHERE status = 'active';

-- ============================================================
-- 4. Notification preferences (per-user opt toggle per event_type x channel).
--    Absence of a row = enabled. See internal/notification/service/preference_checker.go.
-- ============================================================
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

-- ============================================================
-- 5. Dead-letter audit for notifications that exhausted retry budget.
-- ============================================================
CREATE TABLE IF NOT EXISTS dead_letter_notifications (
    id          VARCHAR(36)  PRIMARY KEY,
    event_type  VARCHAR(64)  NOT NULL,
    user_email  VARCHAR(255) NOT NULL,
    payload     TEXT,
    last_error  TEXT,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_dlq_event_type ON dead_letter_notifications (event_type);

-- ============================================================
-- 6. Order status enum values introduced by Stripe + reservation flow.
--    No schema change — orders.status is VARCHAR; the application enforces
--    the state machine (see internal/order/model/order.go CanTransitionTo):
--      pending_payment - awaiting Stripe payment intent confirmation
--      paid            - webhook payment_intent.succeeded received
--      payment_failed  - webhook payment_intent.payment_failed received
-- ============================================================
