---
description: Security best practices for Go backend services
globs:
  - "**/*.go"
  - "pkg/middleware/**"
  - "pkg/jtoken/**"
alwaysApply: false
---

# Security

## Input validation

- Validate ALL user input at the handler layer before passing to service. Use struct validation tags.
- Sanitize string inputs — trim whitespace, limit length, escape HTML where applicable.
- Never trust client-provided IDs for authorization. Always verify ownership from the JWT/session.
- Use parameterized queries only. GORM handles this by default — never use `Raw()` with string concatenation.

## Authentication & Authorization

- JWT tokens must have short expiry. Access tokens: 15-30 min. Refresh tokens: 7 days max.
- Always validate token signature and expiry. Check `type` claim to prevent access/refresh token misuse.
- Role-based access: check `role` from JWT context, not from request body.
- Passwords: use bcrypt with cost >= 10. Never log or return passwords in responses.

## Secrets management

- Never hardcode secrets, API keys, or credentials in source code.
- Use environment variables or config files excluded from git (`.gitignore`).
- Never log sensitive data (tokens, passwords, PII). Use structured logging and redact sensitive fields.

## HTTP security

- CORS: restrict origins to known domains in production. Dev can use `*`.
- Rate limiting: apply to auth endpoints (login, register, password reset) at minimum.
- Set security headers: `X-Content-Type-Options: nosniff`, `X-Frame-Options: DENY`.
- Use HTTPS in production. Redirect HTTP to HTTPS.

## Dependencies

- Run `go mod tidy` to remove unused dependencies.
- Audit dependencies with `govulncheck ./...` periodically.
- Pin dependency versions in `go.mod`. Review breaking changes before upgrading.

## Error responses

- Never expose internal error details (stack traces, SQL errors) to clients.
- Return generic messages for 500 errors. Log the full error server-side.
- Use consistent error response format: `{"error": {"code": "...", "message": "..."}}`.
