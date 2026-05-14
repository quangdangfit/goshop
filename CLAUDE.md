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
| [configuration](/.claude/rules/configuration.md) | `pkg/config/**`, `config.yaml` | Env vars + root-level `config.yaml` (override via `CONFIG_FILE`) |
| [grpc-proto](/.claude/rules/grpc-proto.md) | `proto/**` | Proto generation, plugin versions |
| [shared-packages](/.claude/rules/shared-packages.md) | `pkg/**` | dbs, redis, jtoken, middleware, utils |

---

# Completed Initiatives

The four roadmap initiatives below are **shipped** — BE, FE, and tests (incl.
testcontainers integration suites) all live in the tree. This section is kept
as a map of where each lives so future work can extend the right package.

## 1. Client-side cart ✅
- Server-side cart removed. No `internal/cart/` package; no cart proto service.
- `POST /api/v1/orders` accepts `{items: [{product_id, quantity}]}`; server re-validates products, prices, and stock and snapshots authoritative `price`/`name` onto order lines.
- FE cart lives in `web/src` (localStorage-backed) via `CartPage.tsx` + `CheckoutPage.tsx`.

## 2. Stripe payments ✅
- `pkg/payment/stripe` implements the `Provider` interface; `internal/payment/` owns the domain (model, repo, service, HTTP port).
- Endpoints: `POST /api/v1/orders/:id/payment-intent`, `POST /api/v1/webhooks/stripe` (signature-verified, no JWT), `GET /api/v1/config/public`.
- `payments` and `stripe_events` (provider event dedupe) tables auto-migrated. Status flow: `pending_payment → paid` on webhook success; `payment_failed`/`cancelled` releases stock.
- FE: `PaymentPage.tsx` with `@stripe/react-stripe-js`.

## 3. Strict inventory ✅
- `products.reserved_quantity` column + `stock_reservations` table.
- `ProductRepository.ReserveStock` uses a single conditional UPDATE (no `SELECT … FOR UPDATE`); CHECK constraint on stock is the safety net.
- Reservations expire after 15 min; background sweeper (`runReservationSweeper` in `cmd/api/main.go`) releases expired ones and cancels still-`pending_payment` orders.
- Concurrency-tested via `reservation_integration_test.go` (testcontainers Postgres).

## 4. Notifications ✅
- `pkg/notification` channels (logger, email/SMTP); `pkg/eventbus` in-process pub/sub.
- `internal/notification/` domain with `Preference` and `DeadLetterNotification` models, worker pool with retry/backoff, MailHog integration test in `pkg/notification/email_integration_test.go`.
- Endpoints: `GET/PUT /api/v1/me/notification-preferences`.
- Order/payment services publish `OrderCreated`/`OrderPaid`/`OrderCancelled` events.

## Cross-cutting test infrastructure ✅
- `internal/testutil/postgres.go` spins up Postgres testcontainers with migrations.
- Integration tests are tagged `*_integration_test.go` and live alongside the package under test.

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
