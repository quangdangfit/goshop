# Migrations

`AutoMigrate` in `cmd/api/main.go` handles development. **Production must apply the
numbered SQL files in this directory in order**; each is forward-only and idempotent
(`IF NOT EXISTS` / `IF EXISTS`).

## Conventions

- Filename: `NNNN_short_description.up.sql`. Optional `.down.sql` for emergency rollback.
- One concern per file going forward (the consolidated `0001_init.up.sql` covers the
  whole pre-prod batch — split future changes).
- Never drop columns or rename — add new columns/tables and deprecate. Long-running
  consumers may still read old shapes during a deploy.
- Add indexes in a separate migration from the column they back; index creation can
  acquire heavy locks on large tables.

## Files

| # | File | Purpose |
|---|------|---------|
| 0001 | `0001_init.up.sql` | Drop server-side cart, add Stripe `payments` + `provider_events`, add `stock_reservations` + `products.reserved_quantity` + CHECK guards, add `notification_preferences` and `dead_letter_notifications`, document new `orders.status` enum values. |

## Tooling

This project does not bundle a migration runner. Use `golang-migrate`:

```
migrate -path migrations -database "$DATABASE_URI" up
```
