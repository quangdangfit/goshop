-- Emergency rollback for 0001_init_schema.up.sql. Drops every table the initial
-- schema introduces. Use with care — production data goes with it.

DROP TABLE IF EXISTS wishlists       CASCADE;
DROP TABLE IF EXISTS reviews         CASCADE;
DROP TABLE IF EXISTS order_lines     CASCADE;
DROP TABLE IF EXISTS stock_reservations CASCADE;
DROP TABLE IF EXISTS payments        CASCADE;
DROP TABLE IF EXISTS provider_events CASCADE;
DROP TABLE IF EXISTS preferences     CASCADE;
DROP TABLE IF EXISTS dead_letter_notifications CASCADE;
DROP TABLE IF EXISTS orders          CASCADE;
DROP TABLE IF EXISTS coupons         CASCADE;
DROP TABLE IF EXISTS addresses       CASCADE;
DROP TABLE IF EXISTS products        CASCADE;
DROP TABLE IF EXISTS categories      CASCADE;
DROP TABLE IF EXISTS users           CASCADE;
