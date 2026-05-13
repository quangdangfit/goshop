---
description: Application configuration and environment variables
globs:
  - "pkg/config/**"
  - "config.yaml"
  - "config.sample.yaml"
alwaysApply: false
---

# Configuration

Config is loaded from environment variables. For local development, copy `config.sample.yaml` → `config.yaml` at the **repo root** (loaded from the working directory; override with `CONFIG_FILE=/path/to/config.yaml`).

## Required

- `environment` — set to `production` to enable gin release mode
- `http_port`, `grpc_port`
- `auth_secret` — JWT signing secret
- `database_uri` — PostgreSQL DSN (e.g. `postgres://user:pass@localhost:5432/dbname`)
- `redis_uri`, `redis_password`, `redis_db`

## Optional / feature-specific

- `cors_allowed_origins` (default `*`), `rate_limit_requests` (100), `rate_limit_window_seconds` (60)
- Stripe (payments): `stripe_secret_key`, `stripe_webhook_secret`, `stripe_publishable_key`, `stripe_api_base` (test override pointing at stripe-mock)
- SMTP (notifications): `smtp_host`, `smtp_port` (default 25), `smtp_user`, `smtp_password`, `email_from`. In dev, point at MailHog (`localhost:1025`).
