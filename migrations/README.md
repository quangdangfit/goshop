# Migrations

Schema is managed entirely by versioned SQL files in this directory and applied
with [golang-migrate](https://github.com/golang-migrate/migrate). The app **does
not** call `AutoMigrate` on startup — production deploys must run `migrate up`
against the target database before (or as part of) rolling out a new image.

## Conventions

- Filename: `NNNN_short_description.up.sql` and matching `NNNN_short_description.down.sql`.
- `NNNN` is a strictly increasing 4-digit sequence. Use `make migrate-new name=add_foo`
  to scaffold a pair with the next number.
- One concern per file. Don't bundle unrelated changes.
- Forward-only and idempotent (`IF NOT EXISTS` / `IF EXISTS`) so re-runs are safe.
- Never drop columns or rename in a single step — add new, deprecate old. Long-running
  consumers may still read the old shape during a deploy.
- Add indexes in a separate migration from the column they back; index creation can
  acquire heavy locks on large tables.

## Files

| # | File | Purpose |
|---|------|---------|
| 0001 | `0001_init_schema.up.sql` | Full base schema: 14 tables (users/addresses/wishlists/categories/products/reviews/coupons/orders/order_lines/stock_reservations/payments/provider_events/preferences/dead_letter_notifications) plus PKs, indexes, FKs, the `chk_products_reserved_lte_stock` safety CHECK, and the partial `idx_stock_reservations_expires_at WHERE status='active'` for the sweeper. |

## Local development

Install the CLI once:

```bash
# macOS
brew install golang-migrate
# Go install
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

Then from the repo root:

```bash
# Apply everything
make migrate-up

# Roll back one step
make migrate-down

# Current version
make migrate-status

# Scaffold the next pair
make migrate-new name=add_orders_warehouse_id
```

The Makefile reads `DATABASE_URI` from the environment (defaults to
`postgres://postgres:test@localhost:5432/goshop?sslmode=disable` — match your
local Postgres).

## Production / Kubernetes

Run `migrate up` from a **separate job or init container**, not from the app
process. The image bundles the migrations directory at `/app/migrations`, so a
typical sidecar uses `migrate/migrate:v4` with that path:

```yaml
# k8s init container example
- name: db-migrate
  image: migrate/migrate:v4
  args:
    - "-path=/migrations"
    - "-database=$(DATABASE_URI)"
    - "up"
  envFrom: [...]
  volumeMounts:
    - { name: migrations, mountPath: /migrations }
```

## Tests

Integration tests (`tests/integration/...`) apply these migrations to a fresh
testcontainers Postgres via `testutil.ApplyMigrations(db)` — giving CI parity
with the schema production will see.
