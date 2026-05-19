-- Initial schema for goshop. Mirrors the AutoMigrate output that previously ran on
-- app startup. Every CREATE is guarded with IF NOT EXISTS so this migration is safe
-- to apply (or re-apply) on a fresh database or one that already has the tables.

-- ============================================================
-- 1. Tables
-- ============================================================

CREATE TABLE IF NOT EXISTS addresses (
    id text NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    user_id text NOT NULL,
    name text,
    phone text,
    street text,
    city text,
    country text,
    is_default boolean DEFAULT false
);

CREATE TABLE IF NOT EXISTS categories (
    id text NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    name text NOT NULL,
    slug text NOT NULL,
    description text
);

CREATE TABLE IF NOT EXISTS coupons (
    id text NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    code text NOT NULL,
    discount_type text,
    discount_value numeric,
    min_order_amount numeric DEFAULT 0,
    max_usage bigint DEFAULT 0,
    used_count bigint DEFAULT 0,
    expires_at timestamp with time zone
);

CREATE TABLE IF NOT EXISTS dead_letter_notifications (
    id text NOT NULL,
    created_at timestamp with time zone,
    event_type character varying(64) NOT NULL,
    user_email character varying(255) NOT NULL,
    payload text,
    last_error text
);

CREATE TABLE IF NOT EXISTS order_lines (
    id text NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    order_id text,
    product_id text,
    quantity bigint,
    price numeric
);

CREATE TABLE IF NOT EXISTS orders (
    id text NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    code text,
    user_id text,
    total_price numeric,
    discount_amount numeric DEFAULT 0,
    final_price numeric DEFAULT 0,
    coupon_code text,
    status text
);

CREATE TABLE IF NOT EXISTS payments (
    id text NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    order_id text NOT NULL,
    provider text NOT NULL,
    provider_intent_id text NOT NULL,
    amount bigint NOT NULL,
    currency text NOT NULL,
    status text NOT NULL
);

CREATE TABLE IF NOT EXISTS preferences (
    id text NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    user_id text NOT NULL,
    event_type character varying(64) NOT NULL,
    channel character varying(32) NOT NULL,
    enabled boolean DEFAULT true NOT NULL
);

CREATE TABLE IF NOT EXISTS products (
    id text NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    code text,
    name text,
    description text,
    price numeric,
    active boolean DEFAULT true,
    stock_quantity bigint DEFAULT 0,
    reserved_quantity bigint DEFAULT 0,
    avg_rating numeric DEFAULT 0,
    review_count bigint DEFAULT 0,
    images text,
    category_id text,
    CONSTRAINT chk_products_reserved_quantity CHECK ((reserved_quantity >= 0)),
    CONSTRAINT chk_products_stock_quantity CHECK ((stock_quantity >= 0)),
    CONSTRAINT chk_products_reserved_lte_stock CHECK ((reserved_quantity <= stock_quantity))
);

CREATE TABLE IF NOT EXISTS provider_events (
    created_at timestamp with time zone,
    provider character varying(32) NOT NULL,
    event_id character varying(128) NOT NULL
);

CREATE TABLE IF NOT EXISTS reviews (
    id text NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    user_id text NOT NULL,
    product_id text NOT NULL,
    rating bigint,
    comment text
);

CREATE TABLE IF NOT EXISTS stock_reservations (
    id text NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    order_id text NOT NULL,
    product_id text NOT NULL,
    quantity bigint NOT NULL,
    status text NOT NULL,
    expires_at timestamp with time zone NOT NULL
);

CREATE TABLE IF NOT EXISTS users (
    id text NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    email text NOT NULL,
    password text,
    role text
);

CREATE TABLE IF NOT EXISTS wishlists (
    id text NOT NULL,
    created_at timestamp with time zone,
    user_id text NOT NULL,
    product_id text NOT NULL
);

-- ============================================================
-- 2. Primary keys and unique constraints
-- ============================================================

ALTER TABLE ONLY dead_letter_notifications
    ADD CONSTRAINT dead_letter_notifications_pkey PRIMARY KEY (id);

ALTER TABLE ONLY payments
    ADD CONSTRAINT payments_pkey PRIMARY KEY (id);

ALTER TABLE ONLY preferences
    ADD CONSTRAINT preferences_pkey PRIMARY KEY (id);

ALTER TABLE ONLY provider_events
    ADD CONSTRAINT provider_events_pkey PRIMARY KEY (provider, event_id);

ALTER TABLE ONLY stock_reservations
    ADD CONSTRAINT stock_reservations_pkey PRIMARY KEY (id);

ALTER TABLE ONLY addresses
    ADD CONSTRAINT uni_addresses_id PRIMARY KEY (id);

ALTER TABLE ONLY categories
    ADD CONSTRAINT uni_categories_id PRIMARY KEY (id);

ALTER TABLE ONLY coupons
    ADD CONSTRAINT uni_coupons_id PRIMARY KEY (id);

ALTER TABLE ONLY order_lines
    ADD CONSTRAINT uni_order_lines_id PRIMARY KEY (id);

ALTER TABLE ONLY orders
    ADD CONSTRAINT uni_orders_id PRIMARY KEY (id);

ALTER TABLE ONLY products
    ADD CONSTRAINT uni_products_id PRIMARY KEY (id);

ALTER TABLE ONLY reviews
    ADD CONSTRAINT uni_reviews_id PRIMARY KEY (id);

ALTER TABLE ONLY users
    ADD CONSTRAINT uni_users_email UNIQUE (email);

ALTER TABLE ONLY users
    ADD CONSTRAINT uni_users_id PRIMARY KEY (id);

ALTER TABLE ONLY wishlists
    ADD CONSTRAINT uni_wishlists_id PRIMARY KEY (id);

-- ============================================================
-- 3. Indexes
-- ============================================================

CREATE INDEX IF NOT EXISTS idx_addresses_deleted_at ON addresses USING btree (deleted_at);

CREATE INDEX IF NOT EXISTS idx_addresses_id ON addresses USING btree (id);

CREATE INDEX IF NOT EXISTS idx_addresses_user_id ON addresses USING btree (user_id);

CREATE INDEX IF NOT EXISTS idx_categories_deleted_at ON categories USING btree (deleted_at);

CREATE INDEX IF NOT EXISTS idx_categories_id ON categories USING btree (id);

CREATE UNIQUE INDEX IF NOT EXISTS idx_categories_name ON categories USING btree (name);

CREATE UNIQUE INDEX IF NOT EXISTS idx_categories_slug ON categories USING btree (slug);

CREATE UNIQUE INDEX IF NOT EXISTS idx_coupons_code ON coupons USING btree (code);

CREATE INDEX IF NOT EXISTS idx_coupons_deleted_at ON coupons USING btree (deleted_at);

CREATE INDEX IF NOT EXISTS idx_coupons_id ON coupons USING btree (id);

CREATE INDEX IF NOT EXISTS idx_dead_letter_notifications_event_type ON dead_letter_notifications USING btree (event_type);

CREATE INDEX IF NOT EXISTS idx_order_lines_deleted_at ON order_lines USING btree (deleted_at);

CREATE INDEX IF NOT EXISTS idx_order_lines_id ON order_lines USING btree (id);

CREATE INDEX IF NOT EXISTS idx_orders_deleted_at ON orders USING btree (deleted_at);

CREATE INDEX IF NOT EXISTS idx_orders_id ON orders USING btree (id);

CREATE INDEX IF NOT EXISTS idx_payments_deleted_at ON payments USING btree (deleted_at);

CREATE UNIQUE INDEX IF NOT EXISTS idx_payments_order_id ON payments USING btree (order_id);

CREATE INDEX IF NOT EXISTS idx_payments_provider_intent_id ON payments USING btree (provider_intent_id);

CREATE INDEX IF NOT EXISTS idx_payments_status ON payments USING btree (status);

CREATE UNIQUE INDEX IF NOT EXISTS idx_pref_user_event_channel ON preferences USING btree (user_id, event_type, channel);

CREATE UNIQUE INDEX IF NOT EXISTS idx_product_code ON products USING btree (code);

CREATE UNIQUE INDEX IF NOT EXISTS idx_product_name ON products USING btree (name);

CREATE INDEX IF NOT EXISTS idx_products_deleted_at ON products USING btree (deleted_at);

CREATE INDEX IF NOT EXISTS idx_products_id ON products USING btree (id);

CREATE UNIQUE INDEX IF NOT EXISTS idx_review_user_product ON reviews USING btree (user_id, product_id);

CREATE INDEX IF NOT EXISTS idx_reviews_deleted_at ON reviews USING btree (deleted_at);

CREATE INDEX IF NOT EXISTS idx_reviews_id ON reviews USING btree (id);

CREATE INDEX IF NOT EXISTS idx_stock_reservations_deleted_at ON stock_reservations USING btree (deleted_at);

-- Partial index — sweeper only scans active reservations, so the index doesn't bloat
-- with terminal-state rows.
CREATE INDEX IF NOT EXISTS idx_stock_reservations_expires_at ON stock_reservations USING btree (expires_at) WHERE status = 'active';

CREATE INDEX IF NOT EXISTS idx_stock_reservations_order_id ON stock_reservations USING btree (order_id);

CREATE INDEX IF NOT EXISTS idx_stock_reservations_product_id ON stock_reservations USING btree (product_id);

CREATE INDEX IF NOT EXISTS idx_stock_reservations_status ON stock_reservations USING btree (status);

CREATE INDEX IF NOT EXISTS idx_user_email ON users USING btree (email);

CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users USING btree (deleted_at);

CREATE INDEX IF NOT EXISTS idx_users_id ON users USING btree (id);

CREATE UNIQUE INDEX IF NOT EXISTS idx_wishlist_user_product ON wishlists USING btree (user_id, product_id);

CREATE INDEX IF NOT EXISTS idx_wishlists_id ON wishlists USING btree (id);

-- ============================================================
-- 4. Foreign keys
-- ============================================================

ALTER TABLE ONLY order_lines
    ADD CONSTRAINT fk_order_lines_product FOREIGN KEY (product_id) REFERENCES products(id);

ALTER TABLE ONLY order_lines
    ADD CONSTRAINT fk_orders_lines FOREIGN KEY (order_id) REFERENCES orders(id);

ALTER TABLE ONLY orders
    ADD CONSTRAINT fk_orders_user FOREIGN KEY (user_id) REFERENCES users(id);

ALTER TABLE ONLY products
    ADD CONSTRAINT fk_products_category FOREIGN KEY (category_id) REFERENCES categories(id);

