---
description: Logging and observability patterns for Go backend services
globs:
  - "**/*.go"
  - "pkg/middleware/**"
alwaysApply: false
---

# Logging & Observability

## Structured logging

- Use structured logging (key-value pairs), not `fmt.Sprintf` for log messages.
- Log levels: `Debug` (development only), `Info` (normal operations), `Warn` (recoverable issues), `Error` (failures requiring attention).
- Always include contextual fields: `request_id`, `user_id`, `method`, `path`, `duration`.
- Log at service boundaries (handler entry/exit), not inside utility functions.

## What to log

- All incoming requests (method, path, status code, duration).
- Authentication failures (without exposing credentials).
- Database errors and slow queries (> 200ms).
- External service calls (URL, status, duration).
- Business logic decisions that affect flow (order state transitions, payment results).

## What NOT to log

- Passwords, tokens, API keys, or any credentials.
- Full request/response bodies in production (too verbose, potential PII).
- Expected conditions (cache miss, 404 for unknown ID) — these are noise.
- Health check endpoints (floods logs with no value).

## Request tracing

- Generate a unique `request_id` per request. Pass it via `context.Context`.
- Include `request_id` in all log entries and error responses.
- For gRPC, use metadata to propagate trace context.

## Metrics to track

- Request rate, error rate, and latency (RED method) per endpoint.
- Database connection pool utilization.
- Cache hit/miss ratio.
- Goroutine count and memory usage.
