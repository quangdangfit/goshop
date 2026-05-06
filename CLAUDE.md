# CLAUDE.md

Project-specific rules are in `.claude/rules/`.

## Quick Reference

```bash
# Build
go build -o main cmd/api

# Test
make unittest                    # all tests with coverage
go test ./internal/product/service/... -v -run TestProductServiceTestSuite/TestCreateSuccess

# Lint
golangci-lint run

# Mocks
mockery                          # regenerate all mocks (uses .mockery.yml)

# Swagger
make doc                         # swag fmt && swag init -g ./cmd/api/main.go

# Proto
cd proto && buf generate

# Verify before commit
go vet ./... && golangci-lint run && go test ./... && go mod tidy
```

## Rules Index

| Rule | Scope | Description |
|------|-------|-------------|
| [architecture](/.claude/rules/architecture.md) | `cmd/**`, `internal/**` | Domain structure, dual server setup |
| [go-conventions](/.claude/rules/go-conventions.md) | `**/*.go` | Error handling, naming, idioms |
| [api-design](/.claude/rules/api-design.md) | `**/port/**`, `proto/**` | REST + gRPC standards |
| [database](/.claude/rules/database.md) | `**/repository/**`, `**/model/**` | GORM, transactions, query patterns |
| [security](/.claude/rules/security.md) | `**/*.go`, `pkg/middleware/**` | Input validation, auth, secrets |
| [performance](/.claude/rules/performance.md) | `**/*.go` | Concurrency, caching, optimization |
| [testing](/.claude/rules/testing.md) | `**/*_test.go`, `**/mocks/**` | Testify suites, mocks, handler tests |
| [logging-observability](/.claude/rules/logging-observability.md) | `**/*.go` | Structured logging, tracing, metrics |
| [code-review](/.claude/rules/code-review.md) | always | Quality checklist, PR discipline |
| [git-workflow](/.claude/rules/git-workflow.md) | always | Commit conventions, branch strategy |
| [commands](/.claude/rules/commands.md) | always | Build, test, and dev commands |
| [configuration](/.claude/rules/configuration.md) | `pkg/config/**` | Environment variables, config files |
| [grpc-proto](/.claude/rules/grpc-proto.md) | `proto/**` | Proto generation, plugin versions |
| [shared-packages](/.claude/rules/shared-packages.md) | `pkg/**` | dbs, redis, jtoken, middleware, utils |

---

# Roadmap

This roadmap covers four initiatives. Phases run in this order: **BE → FE → Tests** (unit + integration with testcontainers). Each initiative below is broken down accordingly.

## 1. Client-side cart (drop server-side cart)

**Proposal:** Persist cart only on the client (localStorage / IndexedDB). The server stops storing carts. On `POST /orders`, the client sends the full line items; the server validates products, prices, and stock, then creates the order.

**Assessment — is this a good solution?**

Pros:
- Removes an entire write-heavy domain from the backend (cart table, cart service, cart gRPC). Less code, fewer bugs, lower DB load.
- Cart UX is instant (no network round-trips for add/remove/update qty).
- Works offline; no auth required to build a cart (guest checkout becomes natural).
- Scales trivially — cart state lives with the user.

Cons / risks:
- **No cross-device sync.** A user adding items on mobile won't see them on desktop unless logged in and we add an opt-in sync endpoint.
- **Price/stock drift.** The cart can be hours/days old; products may have changed price, gone out of stock, or been deleted. The server MUST re-validate every line at order time and surface a clear "cart updated" response so the FE can reconcile.
- **No server-side analytics on abandoned carts** unless the FE periodically pings a lightweight `/cart/snapshot` telemetry endpoint (optional, decoupled).
- **Trust boundary.** Never trust client-supplied prices — always re-fetch authoritative price from `products` at order creation. (The current server-side cart already has this risk; client-side just makes it explicit.)

**Verdict:** Good fit for this codebase. The cart domain is currently gRPC-only and lightly used; removing it simplifies the system. The drift and sync concerns are solvable with strict server-side re-validation and an optional logged-in sync endpoint later. **Recommended.**

### Phase 1A — Backend
- Remove `internal/cart` HTTP/gRPC registration from `cmd/api/main.go`.
- Drop `CartService` from gRPC server; remove proto service (mark fields/methods `reserved` rather than renumber).
- Delete `internal/cart/` package and its mocks. Drop `carts` table via a new versioned migration (do not just remove the model — write a migration that drops the table after a deprecation window).
- Extend `OrderService.CreateOrder` to accept a list of `{product_id, quantity}` from the request, then:
  1. Load products in one query (`WHERE id IN (...)`).
  2. Validate each product exists, is active, and `stock >= quantity`.
  3. Snapshot authoritative `price` and `name` from the product row into the order line.
  4. Run stock decrement + order insert inside a single DB transaction (see initiative #3).
- Define a typed error per failure mode: `ErrProductUnavailable`, `ErrInsufficientStock`, `ErrPriceChanged` (optional, if FE wants to confirm). Return a 409 with a structured body listing the offending line items so the FE can reconcile the local cart.
- Update Swagger annotations on `POST /orders`.

### Phase 1B — Frontend
- Cart store in localStorage (key `goshop:cart:v1`), schema-versioned. One module owns add/remove/updateQty/clear.
- Cart page reads from store; on mount, fetch current product info for each line to show fresh price + availability badges (don't mutate the stored cart — show a diff).
- Checkout flow: POST `/orders` with `{items: [{product_id, quantity}]}`. On 409 reconciliation response, show a modal listing changed/unavailable lines and let the user confirm before retrying.
- Optional: if user is logged in, periodic background `PUT /me/cart-snapshot` (telemetry only, not authoritative).

### Phase 1C — Tests
- Unit: `OrderService.CreateOrder` table-driven — happy path, missing product, insufficient stock, price snapshotted correctly, transaction rollback on partial failure.
- Integration (testcontainers — Postgres): create products, fire `POST /orders` with mixed valid/invalid lines, assert 409 body and that no order or stock change persisted.
- FE unit: cart store reducers, schema migration from older `:v0` keys.
- FE e2e (Playwright/Cypress, optional): add to cart → reload → checkout → assert order created.

---

## 2. Stripe payment integration

Goal: charge the customer for an order via Stripe, using **Payment Intents + webhook confirmation** (do not trust client-side success). Order lifecycle: `pending_payment → paid → fulfilled` (or `payment_failed`, `cancelled`).

### Phase 2A — Backend
- New package `pkg/payment` with a `Provider` interface (`CreateIntent`, `RetrieveIntent`, `VerifyWebhook`). Concrete impl `pkg/payment/stripe` wraps `github.com/stripe/stripe-go/v76`.
- Config additions in `pkg/config`: `stripe_secret_key`, `stripe_webhook_secret`, `stripe_publishable_key` (the last one is exposed to FE via a `GET /config/public` endpoint).
- New domain `internal/payment/`:
  - `model.Payment` — `id`, `order_id`, `provider`, `provider_intent_id`, `amount`, `currency`, `status`, timestamps. Migration adds `payments` table.
  - `repository`, `service`, HTTP port.
  - `Order.Status` extended; add `PaymentID *string` on `Order`.
- Endpoints:
  - `POST /api/v1/orders/:id/payment-intent` → creates Stripe PaymentIntent, returns `client_secret`. Idempotent per order.
  - `POST /api/v1/webhooks/stripe` → verifies signature using `stripe_webhook_secret`, advances order/payment status on `payment_intent.succeeded` / `.payment_failed`. **No JWT middleware on this route.**
- On `payment_intent.succeeded`: inside a tx, mark payment `succeeded`, set order `paid`. Stock was already reserved at order creation (initiative #3), so no further inventory work here. On `.payment_failed` or webhook timeout (cron sweep): release reserved stock and mark order `cancelled`.
- Idempotency: store Stripe event IDs in a `stripe_events` table; reject duplicates.
- Never log card data, full PaymentIntent payloads, or `client_secret`.

### Phase 2B — Frontend
- Install `@stripe/stripe-js` + `@stripe/react-stripe-js`. Load publishable key from `/config/public`.
- After `POST /orders` returns the order, call `POST /orders/:id/payment-intent` to get `client_secret`, then mount `<PaymentElement>`.
- On `stripe.confirmPayment` success → redirect to order status page that polls `GET /orders/:id` until status flips to `paid` (webhook-driven).
- Show clear states: pending / processing / paid / failed (with retry).

### Phase 2C — Tests
- Unit: `payment/stripe` provider with mocked `stripe.Client`; webhook signature verification (valid, invalid, expired).
- Unit: order status transitions on each Stripe event type; double-event idempotency.
- Integration (testcontainers — Postgres + Stripe CLI mock or `stripe-mock` container): full flow — create order → create intent → simulate webhook → assert DB state. Use `stripe-mock` (`stripemock/stripe-mock`) container.
- FE unit: payment component states; mock Stripe.js.
- E2E with Stripe test cards (`4242…`) optional in a staging env, not CI.

---

## 3. Strict inventory management

Problem: under concurrency, two users can both pass a `stock >= qty` check and both succeed at order creation, oversold.

**Strategy: pessimistic-decrement-with-constraint.** A DB CHECK constraint (`stock >= 0`) plus a single conditional UPDATE per line is enough to make oversell impossible without explicit row locks. Reservations decouple "user committed to buy" from "payment cleared."

### Phase 3A — Backend
- Migration: add `CHECK (stock >= 0)` on `products.stock`. Add `reserved_stock INT NOT NULL DEFAULT 0` column. "Available stock" = `stock - reserved_stock`.
- Add table `stock_reservations` — `id`, `order_id`, `product_id`, `quantity`, `expires_at`, `status` (`active`/`released`/`committed`).
- `ProductRepository.ReserveStock(ctx, productID, qty) error`: single SQL `UPDATE products SET reserved_stock = reserved_stock + ? WHERE id = ? AND (stock - reserved_stock) >= ?` — check `RowsAffected == 1`; if 0, return `ErrInsufficientStock`. No SELECT…FOR UPDATE needed.
- `CommitReservation(ctx, reservationID)`: in one tx, decrement `stock` and `reserved_stock` by qty, mark reservation `committed`. The CHECK constraint is the safety net.
- `ReleaseReservation(ctx, reservationID)`: decrement `reserved_stock`, mark reservation `released`.
- Order creation flow:
  1. Begin tx → for each line, `ReserveStock` → insert order + lines + reservations → commit.
  2. Reservations have `expires_at = now + 15min`.
  3. Stripe webhook `succeeded` → `CommitReservation` for all lines.
  4. Stripe `failed` / order `cancelled` / TTL expiry → `ReleaseReservation`.
- Background sweeper (goroutine started from `cmd/api/main.go`, ticker every 60s): finds reservations with `expires_at < now AND status='active'`, releases them and cancels the parent order if still `pending_payment`. Must be context-cancellable on shutdown.
- Admin endpoint `POST /admin/products/:id/stock` with role check to add stock (audited).
- Emit `inventory.low_stock` event when `stock - reserved_stock < threshold` (initiative #4).

### Phase 3B — Frontend
- On product page, show "available" = backend `stock - reserved_stock`. Add badges: "Low stock", "Out of stock".
- During checkout, after order creation, show a 15-minute reservation timer. When timer expires, the order is auto-cancelled — show clear messaging.
- On 409 from order creation, surface per-line failure reasons.

### Phase 3C — Tests
- Unit: `ReserveStock` returns `ErrInsufficientStock` when row update affects 0 rows; happy path; commit/release math.
- Integration (testcontainers — Postgres): **concurrency tests** are the headline. Spawn N goroutines all reserving the last unit; assert exactly 1 succeeds and `stock - reserved_stock == 0`.
- Integration: TTL sweeper releases expired reservation and cancels order.
- Integration: full lifecycle — reserve → webhook commit → stock decremented, reservation committed.
- FE unit: stock-badge component thresholds; reservation countdown.

---

## 4. Notifications integration

Build on the existing `pkg/notification` interface (currently logger-only). Add real channels (email, optionally SMS/push) and an event bus so domains emit events without knowing about transports.

### Phase 4A — Backend
- Extend `pkg/notification`:
  - Channels: `EmailSender` (SMTP via `gomail` or transactional like SendGrid/Postmark), `SMSSender` (optional), `WebPushSender` (optional).
  - `Notifier` interface with `Notify(ctx, event Event) error` — fans out to subscribed channels per user preference.
- New domain `internal/notification/`:
  - `model.NotificationPreference` — per-user, per-event-type, per-channel toggle.
  - `model.NotificationLog` — audit trail of sent notifications for debugging/idempotency.
  - HTTP endpoints: `GET/PUT /api/v1/me/notification-preferences`.
- Event types (initial set): `OrderCreated`, `OrderPaid`, `OrderCancelled`, `OrderShipped`, `LowStock` (admin-only), `WelcomeEmail`.
- Internal event bus: lightweight in-process pub/sub (`pkg/eventbus`) with typed events. Domains publish; notification service subscribes. No external broker yet — keep the interface so we can swap to NATS/Kafka later.
- Templating: `text/template` files in `internal/notification/templates/` for email subject + HTML/text bodies.
- Async delivery: notification service consumes events on a worker goroutine pool with bounded queue; failures retry with backoff (max 3) then go to a `dead_letter_notifications` table.
- Config: `smtp_host`, `smtp_port`, `smtp_user`, `smtp_password`, `email_from`, optional `sendgrid_api_key`. Never log credentials or full email bodies.
- Wire-up: `OrderService.CreateOrder` publishes `OrderCreated`; Stripe webhook publishes `OrderPaid` / `OrderCancelled`; inventory sweeper publishes `LowStock`.

### Phase 4B — Frontend
- Notification preferences page under user settings — toggle per event-type × channel.
- In-app toast notifications via WebSocket or SSE endpoint `GET /api/v1/me/notifications/stream` (optional, later).
- Order detail page shows notification history for that order (from `NotificationLog`).

### Phase 4C — Tests
- Unit: each channel sender with a mocked transport (assert subject/body/recipient); template rendering with golden files; preference filtering.
- Unit: event bus delivers to all subscribers; one failing subscriber does not block others; retry/backoff logic.
- Integration (testcontainers — Postgres + [MailHog](https://github.com/mailhog/MailHog) container as fake SMTP): publish `OrderCreated` → assert MailHog received an email matching the template, and `notification_logs` row exists.
- Integration: preferences honored — disabling email skips delivery.
- FE unit: preferences form; stream/toast component.

---

## Cross-cutting test infrastructure

- Add `internal/testutil/` with helpers to spin up Postgres / Redis / MailHog / stripe-mock testcontainers and run migrations. Reuse across all integration suites.
- Integration tests live next to handler code as `*_integration_test.go` with `//go:build integration` build tag. CI runs them in a separate job.
- Each integration test gets an isolated DB schema (or a fresh container per package) — never share state across tests.

## Suggested execution order

1. Initiative #1 BE → #3 BE (they touch the same `CreateOrder` path; do together).
2. Initiative #2 BE.
3. Initiative #4 BE.
4. All FE work (#1, #3, #2, #4 in that UX order).
5. All test phases (unit alongside each BE phase; integration suites last, once contracts are stable).

## Commit Convention

```
feat(payment): add Stripe Checkout Session creation endpoint
fix(webhook): handle out-of-order payment_intent.succeeded after charge.refunded
perf(redis): batch idempotency checks
test(e2e): add full payment flow integration test
```

Types: `feat`, `fix`, `refactor`, `perf`, `test`, `docs`, `chore`, `ci`
Scope: `payment`, `webhook`, `blockchain`, `kafka`, `redis`, `mysql`, `api`, `config`, `bench`

**Author:** commit as configured `git config user.name` / `user.email`. Do **not** append a `Co-Authored-By: Claude …` trailer — repo history has none, keep it that way.
