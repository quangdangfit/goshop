---
description: Application configuration and environment variables
globs:
  - "pkg/config/**"
  - "config.yaml"
alwaysApply: false
---

# Configuration

Config is loaded from environment variables. For local development, place a `config.yaml` in `pkg/config/`. Required variables:

- `environment` — set to `production` to enable gin release mode
- `http_port`, `grpc_port`
- `auth_secret` — JWT signing secret
- `database_uri` — PostgreSQL DSN (e.g. `postgres://user:pass@localhost:5432/dbname`)
- `redis_uri`, `redis_password`, `redis_db`
