---
description: Go language conventions - error handling, naming, idioms for senior backend engineers
globs:
  - "**/*.go"
alwaysApply: false
---

# Go Conventions

## Error handling

- Always handle errors explicitly. Never use `_ = someFunc()` for functions that return errors unless there is a documented reason.
- Wrap errors with context using `fmt.Errorf("operation failed: %w", err)` — preserve the original error for `errors.Is`/`errors.As`.
- Return errors to callers; only log at the top-level handler (HTTP/gRPC handler). Do NOT log-and-return in service/repository layers.
- Use sentinel errors (`var ErrNotFound = errors.New(...)`) for expected conditions. Use custom error types for errors that carry structured data.
- In HTTP handlers, map domain errors to proper HTTP status codes (e.g. `ErrNotFound` → 404, `ErrUnauthorized` → 401).

## Naming

- Follow Go naming conventions: `MixedCaps`/`mixedCaps`, not `snake_case` or `SCREAMING_SNAKE`.
- Interfaces: use `-er` suffix for single-method interfaces (`Reader`, `Writer`). Multi-method interfaces describe capability (`Repository`, `Service`).
- Receivers: short, 1-2 letter, consistent across methods (e.g. `s` for service, `r` for repository, `h` for handler).
- Package names: short, lowercase, singular (`user` not `users`, `model` not `models`).
- Exported functions/types must have doc comments. Unexported ones only need comments if the logic is non-obvious.

## Struct design

- Prefer composition over inheritance. Embed structs for shared behavior.
- Use pointer receivers for methods that modify state; value receivers for methods that don't.
- Zero values should be useful — design structs so the zero value is a valid, safe default.

## Imports

- Group imports in 3 blocks: stdlib, external packages, internal packages. Separated by blank lines.
- Never use dot imports (`. "package"`) or blank imports (`_ "package"`) except for driver registration (e.g. `_ "github.com/lib/pq"`).

## Functions

- Accept interfaces, return structs.
- Context (`context.Context`) is always the first parameter.
- Keep functions short (< 50 lines). Extract complex logic into well-named helpers.
- Prefer early returns to reduce nesting.

## Testing

- Table-driven tests are the default pattern for multiple test cases.
- Test function names: `TestFunctionName_Scenario` or use testify suites.
- Use `t.Helper()` in test helper functions.
- Never test unexported functions directly — test through the public API.