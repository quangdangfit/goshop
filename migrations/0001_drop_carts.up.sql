-- Drop the server-side cart table. Cart now persists client-side (localStorage).
-- Run only after the deprecated /cart gRPC service is no longer registered (see
-- internal/cart removal).

DROP TABLE IF EXISTS cart_lines;
DROP TABLE IF EXISTS carts;
