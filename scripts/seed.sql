-- Seed data for GoShop
-- Run: psql -U postgres -d goshop -f scripts/seed.sql

-- ============================================================
-- Categories
-- ============================================================
INSERT INTO categories (id, name, slug, description, created_at, updated_at)
VALUES
  ('cat-electronics',   'Electronics',       'electronics',       'Phones, laptops, and gadgets',           NOW(), NOW()),
  ('cat-clothing',      'Clothing',          'clothing',          'Men and women fashion',                  NOW(), NOW()),
  ('cat-home',          'Home & Kitchen',    'home-kitchen',      'Furniture, appliances, and cookware',    NOW(), NOW()),
  ('cat-sports',        'Sports & Outdoors', 'sports-outdoors',   'Equipment and activewear',               NOW(), NOW()),
  ('cat-books',         'Books',             'books',             'Fiction, non-fiction, and textbooks',    NOW(), NOW()),
  ('cat-beauty',        'Beauty & Health',   'beauty-health',     'Skincare, makeup, and wellness',         NOW(), NOW())
ON CONFLICT DO NOTHING;

-- ============================================================
-- Products
-- ============================================================
INSERT INTO products (id, code, name, description, price, stock_quantity, avg_rating, review_count, active, images, category_id, created_at, updated_at)
VALUES

-- Electronics
  ('prod-001', 'P2024001', 'iPhone 15 Pro',
   'Apple iPhone 15 Pro with A17 Pro chip, 48MP camera system, and titanium design.',
   999.99, 50, 4.8, 120, true,
   '["https://images.unsplash.com/photo-1695048133142-1a20484d2569?w=800","https://images.unsplash.com/photo-1592750475338-74b7b21085ab?w=800"]',
   'cat-electronics', NOW(), NOW()),

  ('prod-002', 'P2024002', 'Samsung Galaxy S24 Ultra',
   'Samsung flagship with 200MP camera, built-in S Pen, and AI-powered features.',
   1199.99, 35, 4.7, 95, true,
   '["https://images.unsplash.com/photo-1706439136399-b1e1e36d8c65?w=800"]',
   'cat-electronics', NOW(), NOW()),

  ('prod-003', 'P2024003', 'MacBook Air M3',
   '15-inch MacBook Air powered by Apple M3 chip, 18-hour battery life.',
   1299.99, 20, 4.9, 60, true,
   '["https://images.unsplash.com/photo-1611186871525-a00b10ab5847?w=800","https://images.unsplash.com/photo-1517336714731-489689fd1ca8?w=800"]',
   'cat-electronics', NOW(), NOW()),

  ('prod-004', 'P2024004', 'Sony WH-1000XM5',
   'Industry-leading noise canceling wireless headphones with 30-hour battery.',
   349.99, 80, 4.7, 200, true,
   '["https://images.unsplash.com/photo-1505740420928-5e560c06d30e?w=800","https://images.unsplash.com/photo-1484704849700-f032a568e944?w=800"]',
   'cat-electronics', NOW(), NOW()),

  ('prod-005', 'P2024005', 'iPad Pro 12.9"',
   'Apple iPad Pro with M2 chip, Liquid Retina XDR display, and Apple Pencil support.',
   1099.99, 30, 4.6, 45, true,
   '["https://images.unsplash.com/photo-1544244015-0df4b3ffc6b0?w=800"]',
   'cat-electronics', NOW(), NOW()),

  ('prod-006', 'P2024006', 'Dell XPS 15',
   'Dell XPS 15 with Intel Core i7, OLED display, and NVIDIA RTX 4060 GPU.',
   1599.99, 15, 4.5, 38, true,
   '["https://images.unsplash.com/photo-1593642632559-0c6d3fc62b89?w=800"]',
   'cat-electronics', NOW(), NOW()),

-- Clothing
  ('prod-007', 'P2024007', 'Classic White T-Shirt',
   '100% organic cotton unisex t-shirt, pre-shrunk and machine washable.',
   29.99, 200, 4.5, 310, true,
   '["https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w=800","https://images.unsplash.com/photo-1583743814966-8936f5b7be1a?w=800"]',
   'cat-clothing', NOW(), NOW()),

  ('prod-008', 'P2024008', 'Slim Fit Jeans',
   'Stretch denim slim fit jeans available in multiple washes.',
   59.99, 150, 4.3, 180, true,
   '["https://images.unsplash.com/photo-1542272604-787c3835535d?w=800"]',
   'cat-clothing', NOW(), NOW()),

  ('prod-009', 'P2024009', 'Waterproof Jacket',
   'Lightweight 3-layer waterproof and windproof outdoor jacket.',
   129.99, 75, 4.6, 92, true,
   '["https://images.unsplash.com/photo-1544923246-77307dd654cb?w=800"]',
   'cat-clothing', NOW(), NOW()),

  ('prod-010', 'P2024010', 'Running Shorts',
   'Moisture-wicking 5-inch running shorts with built-in liner.',
   39.99, 120, 4.4, 65, true,
   '["https://images.unsplash.com/photo-1591195853828-11db59a44f43?w=800"]',
   'cat-clothing', NOW(), NOW()),

-- Home & Kitchen
  ('prod-011', 'P2024011', 'Instant Pot Duo 7-in-1',
   'Electric pressure cooker, slow cooker, rice cooker, steamer, sauté, and warmer. 6-quart.',
   89.99, 60, 4.7, 520, true,
   '["https://images.unsplash.com/photo-1585515320310-259814833e62?w=800"]',
   'cat-home', NOW(), NOW()),

  ('prod-012', 'P2024012', 'Nespresso Vertuo Plus',
   'Coffee and espresso machine with 5 cup sizes, milk frother included.',
   179.99, 45, 4.5, 280, true,
   '["https://images.unsplash.com/photo-1495474472287-4d71bcdd2085?w=800","https://images.unsplash.com/photo-1511920170033-f8396924c348?w=800"]',
   'cat-home', NOW(), NOW()),

  ('prod-013', 'P2024013', 'Dyson V15 Detect',
   'Cordless vacuum with laser dust detection and LCD screen showing debris count.',
   699.99, 25, 4.6, 145, true,
   '["https://images.unsplash.com/photo-1558618666-fcd25c85cd64?w=800"]',
   'cat-home', NOW(), NOW()),

  ('prod-014', 'P2024014', 'Cast Iron Skillet 12"',
   'Pre-seasoned cast iron skillet, oven-safe up to 500°F, works on all cooktops.',
   44.99, 100, 4.8, 430, true,
   '["https://images.unsplash.com/photo-1618160702438-9b02ab6515c9?w=800"]',
   'cat-home', NOW(), NOW()),

  ('prod-015', 'P2024015', 'Bamboo Cutting Board Set',
   'Set of 3 eco-friendly bamboo cutting boards with juice groove.',
   34.99, 90, 4.4, 215, true,
   '["https://images.unsplash.com/photo-1590794056226-79ef3a8147e1?w=800"]',
   'cat-home', NOW(), NOW()),

-- Sports & Outdoors
  ('prod-016', 'P2024016', 'Yoga Mat Premium',
   'Non-slip 6mm thick yoga mat with alignment lines, carrying strap included.',
   49.99, 110, 4.6, 340, true,
   '["https://images.unsplash.com/photo-1601925228126-54c5171f3bce?w=800","https://images.unsplash.com/photo-1518611012118-696072aa579a?w=800"]',
   'cat-sports', NOW(), NOW()),

  ('prod-017', 'P2024017', 'Adjustable Dumbbell Set',
   'Space-saving adjustable dumbbells, 5-52.5 lbs per dumbbell, quick-change mechanism.',
   349.99, 30, 4.7, 185, true,
   '["https://images.unsplash.com/photo-1571019613454-1cb2f99b2d8b?w=800"]',
   'cat-sports', NOW(), NOW()),

  ('prod-018', 'P2024018', 'Trail Running Shoes',
   'Lightweight trail running shoes with rock plate and grippy outsole.',
   119.99, 65, 4.5, 98, true,
   '["https://images.unsplash.com/photo-1542291026-7eec264c27ff?w=800","https://images.unsplash.com/photo-1608231387042-66d1773070a5?w=800"]',
   'cat-sports', NOW(), NOW()),

  ('prod-019', 'P2024019', 'Hydration Backpack 15L',
   '15L trail running pack with 2L water reservoir, chest and hip straps.',
   79.99, 55, 4.4, 72, true,
   '["https://images.unsplash.com/photo-1553062407-98eeb64c6a62?w=800"]',
   'cat-sports', NOW(), NOW()),

-- Books
  ('prod-020', 'P2024020', 'The Pragmatic Programmer',
   'Classic software engineering book covering best practices and developer philosophy. 20th Anniversary Edition.',
   44.99, 200, 4.9, 850, true,
   '["https://images.unsplash.com/photo-1544716278-ca5e3f4abd8c?w=800"]',
   'cat-books', NOW(), NOW()),

  ('prod-021', 'P2024021', 'Clean Code',
   'A handbook of agile software craftsmanship by Robert C. Martin.',
   39.99, 180, 4.8, 720, true,
   '["https://images.unsplash.com/photo-1532012197267-da84d127e765?w=800"]',
   'cat-books', NOW(), NOW()),

  ('prod-022', 'P2024022', 'Designing Data-Intensive Applications',
   'The big ideas behind reliable, scalable, and maintainable systems by Martin Kleppmann.',
   54.99, 150, 4.9, 610, true,
   '["https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=800"]',
   'cat-books', NOW(), NOW()),

  ('prod-023', 'P2024023', 'Atomic Habits',
   'An easy and proven way to build good habits and break bad ones by James Clear.',
   19.99, 300, 4.8, 1200, true,
   '["https://images.unsplash.com/photo-1589829085413-56de8ae18c73?w=800"]',
   'cat-books', NOW(), NOW()),

-- Beauty & Health
  ('prod-024', 'P2024024', 'Vitamin C Serum',
   '20% Vitamin C + Hyaluronic Acid + Vitamin E serum for brightening and anti-aging.',
   24.99, 140, 4.5, 390, true,
   '["https://images.unsplash.com/photo-1556228578-8c89e6adf883?w=800","https://images.unsplash.com/photo-1620916566398-39f1143ab7be?w=800"]',
   'cat-beauty', NOW(), NOW()),

  ('prod-025', 'P2024025', 'Electric Toothbrush',
   'Rechargeable electric toothbrush with 3 modes, 2-minute timer, and 30-day battery.',
   49.99, 85, 4.6, 260, true,
   '["https://images.unsplash.com/photo-1559591937-3ae0d5e5b5b7?w=800"]',
   'cat-beauty', NOW(), NOW())

ON CONFLICT DO NOTHING;

-- ============================================================
-- Users
-- Passwords are bcrypt(cost=10). Both accounts use the same plain-text
-- password "password123" / "admin123" for local dev only — DO NOT use in prod.
-- ============================================================
INSERT INTO users (id, email, password, role, created_at, updated_at)
VALUES
  ('user-customer', 'customer@goshop.local',
   '$2a$10$1IHKAaa2QA7sBwVfMtWwm.5pLTeWrxnkFSSj.nLGV.BDCzSAw5EjC', -- password123
   'customer', NOW(), NOW()),
  ('user-admin',    'admin@goshop.local',
   '$2a$10$31sruWVh64GtIJ9pQj..LuqBtTVqDEmkUNBHelX3/JhXSM9XqItYy', -- admin123
   'admin',    NOW(), NOW())
ON CONFLICT DO NOTHING;

-- ============================================================
-- Addresses (one default for the customer)
-- ============================================================
INSERT INTO addresses (id, user_id, name, phone, street, city, country, is_default, created_at, updated_at)
VALUES
  ('addr-001', 'user-customer', 'Test Customer', '+84-900-000-001',
   '123 Đường Lê Lợi', 'Hồ Chí Minh', 'VN', true, NOW(), NOW())
ON CONFLICT DO NOTHING;

-- ============================================================
-- Wishlist
-- ============================================================
INSERT INTO wishlists (id, user_id, product_id, created_at)
VALUES
  ('wish-001', 'user-customer', 'prod-001', NOW()),
  ('wish-002', 'user-customer', 'prod-020', NOW())
ON CONFLICT DO NOTHING;

-- ============================================================
-- Coupon
-- ============================================================
INSERT INTO coupons (id, code, discount_type, discount_value, min_order_amount, max_usage, used_count, expires_at, created_at, updated_at)
VALUES
  ('coupon-welcome10', 'WELCOME10', 'percent', 10, 0, 1000, 0,
   NOW() + INTERVAL '30 days', NOW(), NOW())
ON CONFLICT DO NOTHING;

-- ============================================================
-- Sample order in pending_payment (with reserved stock)
-- Demonstrates the full Stripe + reservation flow.
-- ============================================================
INSERT INTO orders (id, code, user_id, total_price, final_price, discount_amount, status, created_at, updated_at)
VALUES
  ('order-demo-1', 'SO20260513001', 'user-customer', 1349.98, 1349.98, 0, 'pending_payment', NOW(), NOW())
ON CONFLICT DO NOTHING;

INSERT INTO order_lines (id, order_id, product_id, quantity, price, created_at, updated_at)
VALUES
  ('ol-demo-1a', 'order-demo-1', 'prod-001', 1, 999.99, NOW(), NOW()),
  ('ol-demo-1b', 'order-demo-1', 'prod-004', 1, 349.99, NOW(), NOW())
ON CONFLICT DO NOTHING;

-- Reserve stock for the pending order. Bump products.reserved_quantity to match.
INSERT INTO stock_reservations (id, order_id, product_id, quantity, status, expires_at, created_at, updated_at)
VALUES
  ('res-demo-1a', 'order-demo-1', 'prod-001', 1, 'active', NOW() + INTERVAL '15 minutes', NOW(), NOW()),
  ('res-demo-1b', 'order-demo-1', 'prod-004', 1, 'active', NOW() + INTERVAL '15 minutes', NOW(), NOW())
ON CONFLICT DO NOTHING;

UPDATE products SET reserved_quantity = reserved_quantity + 1 WHERE id IN ('prod-001', 'prod-004');

-- Stripe payment record (pending — webhook will flip status to 'succeeded')
INSERT INTO payments (id, order_id, provider, provider_intent_id, amount, currency, status, created_at, updated_at)
VALUES
  ('pay-demo-1', 'order-demo-1', 'stripe', 'pi_demo_seed_1', 134998, 'usd', 'pending', NOW(), NOW())
ON CONFLICT DO NOTHING;

-- ============================================================
-- Notification preferences (customer opts out of order_status_changed email)
-- ============================================================
INSERT INTO notification_preferences (id, user_id, event_type, channel, enabled, created_at, updated_at)
VALUES
  ('pref-001', 'user-customer', 'order_status_changed', 'email', false, NOW(), NOW())
ON CONFLICT DO NOTHING;
