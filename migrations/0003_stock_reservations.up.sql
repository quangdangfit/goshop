-- Strict inventory: track reserved-but-not-yet-committed stock so concurrent checkouts
-- cannot oversell. The CHECK is the safety net — even with a buggy app the DB will
-- reject any write that would drive a counter negative.

ALTER TABLE products
    ADD COLUMN IF NOT EXISTS reserved_quantity INT NOT NULL DEFAULT 0;

ALTER TABLE products
    ADD CONSTRAINT chk_products_stock_nonneg       CHECK (stock_quantity >= 0),
    ADD CONSTRAINT chk_products_reserved_nonneg    CHECK (reserved_quantity >= 0),
    ADD CONSTRAINT chk_products_reserved_lte_stock CHECK (reserved_quantity <= stock_quantity);

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
