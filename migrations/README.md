# Migrations

`AutoMigrate` in `cmd/api/main.go` handles development. **Production must apply the
numbered SQL files in this directory in order**; each is forward-only and idempotent
(`IF NOT EXISTS`/`ADD COLUMN IF NOT EXISTS`).

## Conventions

- Filename: `NNNN_short_description.up.sql`. Optional `.down.sql` for emergency rollback.
- One concern per file. No mixing schema with data backfills.
- Never drop columns or rename — add new columns/tables and deprecate. Long-running
  consumers may still read old shapes during a deploy.
- Add indexes in a separate migration from the column they back; index creation can
  acquire heavy locks on large tables.

## Files

| # | File | Purpose |
|---|------|---------|
| 0001 | `0001_drop_carts.up.sql` | Drop server-side `carts` table — cart now lives client-side. |
| 0002 | `0002_payments.up.sql` | `payments` + `provider_events` tables for Stripe integration. |
| 0003 | `0003_stock_reservations.up.sql` | `stock_reservations` + `products.reserved_quantity` + CHECK. |
| 0004 | `0004_notification_preferences.up.sql` | `notification_preferences` table (per-user opt toggles). |
| 0005 | `0005_dead_letter_notifications.up.sql` | `dead_letter_notifications` for retry exhaustion audit. |
| 0006 | `0006_order_status_pending_payment.up.sql` | No-op: status column was already varchar; documents the new enum values. |

## Tooling

This project does not bundle a migration runner. Use `golang-migrate`:

```
migrate -path migrations -database "$DATABASE_URI" up
```
