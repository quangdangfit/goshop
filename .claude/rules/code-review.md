---
description: Code review standards and quality checklist for Go backend code
globs:
alwaysApply: true
---

# Code Review Standards

## Before writing code

- Understand the full requirement before implementing. Read related code first.
- Check if similar functionality already exists — reuse before creating.
- Consider backward compatibility and migration paths.

## Code quality checklist

- [ ] All errors are handled, not silently ignored
- [ ] No hardcoded values — use constants or config
- [ ] Functions do one thing and are < 50 lines
- [ ] No unused imports, variables, or dead code
- [ ] Race conditions considered for concurrent code
- [ ] SQL injection is impossible (parameterized queries only)
- [ ] Sensitive data is not logged or exposed in responses
- [ ] Tests cover the happy path AND error paths
- [ ] Edge cases handled: nil pointers, empty slices, zero values

## When modifying existing code

- Do not change function signatures without updating all callers.
- Run existing tests before and after changes. All must pass.
- If adding a new interface method, update all implementations AND mocks.
- When renaming, use project-wide find-and-replace. Check imports, proto files, and generated code.

## PR discipline

- One PR = one concern. Don't mix refactoring with feature work.
- PR title format: `type(scope): description` — e.g. `feat(product): add category filtering`.
- Types: `feat`, `fix`, `refactor`, `test`, `docs`, `chore`, `perf`.
- Every PR must include tests for new/changed behavior.
- Update Swagger docs if HTTP API surface changed.

## Go-specific review points

- `defer` placement: immediately after resource acquisition, not buried in logic.
- Channel and goroutine lifecycle: every goroutine must have a clear exit path.
- Interface satisfaction: verify with `var _ Interface = (*Struct)(nil)` compile-time checks.
- Avoid `interface{}` / `any` unless absolutely necessary. Use generics for type-safe containers.
