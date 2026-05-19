---
description: Database patterns - GORM, migrations, transactions, query optimization
globs:
  - "**/repository/**"
  - "**/model/**"
  - "pkg/dbs/**"
  - "cmd/api/main.go"
alwaysApply: false
---

# Database Patterns

## GORM models

- Every model must have `ID`, `CreatedAt`, `UpdatedAt`. Use `gorm.Model` or define explicitly.
- Use `BeforeCreate` hooks for UUID generation and code generation.
- Define database constraints in struct tags: `gorm:"uniqueIndex"`, `gorm:"not null"`, `gorm:"type:varchar(255)"`.
- Use pointer types (`*string`, `*time.Time`) for nullable fields. Non-pointer fields are NOT NULL by default in GORM.
- Add `TableName()` method to control table naming explicitly.

## Repository layer

- Repositories handle ONLY data access. No business logic.
- Accept `context.Context` as first parameter for cancellation and tracing.
- Use the `dbs.Database` interface with functional options (`WithQuery`, `WithPreload`, etc.) — never call GORM directly outside `pkg/dbs`.
- Return domain errors (`ErrNotFound`, `ErrDuplicate`) instead of raw GORM/SQL errors.

## Transactions

- Use transactions for any operation that modifies multiple tables or rows.
- Keep transactions short — do validation before starting the transaction.
- Always handle rollback: use `defer tx.Rollback()` with commit at the end.
- Never hold locks longer than necessary. Avoid SELECT...FOR UPDATE unless truly needed.

## Migrations

- Schema is managed by versioned SQL files in `migrations/` and applied with
  [golang-migrate](https://github.com/golang-migrate/migrate). **`AutoMigrate` is not
  called by the app** — neither in `cmd/api/main.go` nor in tests. Tests apply the same
  SQL files via `testutil.ApplyMigrations(db)`.
- Filenames: `NNNN_short_description.up.sql` + matching `.down.sql`. Scaffold the next
  pair with `make migrate-new name=<desc>`. One concern per file.
- Forward-only and idempotent (`IF NOT EXISTS` / `IF EXISTS`). Never drop columns or
  rename in a single step — add new, deprecate old. Long-running consumers may still
  read the old shape during a deploy.
- Add indexes in a separate migration from the column they back; index creation can
  acquire heavy locks on large tables.
- Local dev: `make migrate-up` / `migrate-down` / `migrate-status`. Production / k8s:
  apply via a separate Job or init container — the image bundles `migrations/` at
  `/app/migrations`.

## Query patterns

- Use `Preload` for eager loading associations. Use `Joins` for filtering by association.
- Always paginate list queries. Default limit: 20, max limit: 100.
- Use `Select()` to limit columns when you don't need the full model.
- For complex queries, use raw SQL via `Raw()` with parameterized placeholders — never string concatenation.
- Use `Count()` separately from `Find()` for paginated responses.
