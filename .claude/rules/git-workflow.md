---
description: Git workflow and commit conventions
globs:
alwaysApply: true
---

# Git Workflow

## Commit messages

- Format: `type(scope): concise description`
- Types: `feat`, `fix`, `refactor`, `test`, `docs`, `chore`, `perf`, `ci`
- Scope: domain or package name (`product`, `user`, `middleware`, `dbs`)
- Body (optional): explain WHY, not WHAT. The diff shows what changed.
- Examples:
  - `feat(order): add status workflow with validation`
  - `fix(auth): prevent refresh token reuse after rotation`
  - `refactor(product): extract cache invalidation to helper`

## Branch strategy

- `master` — production-ready code. Protected.
- Feature branches: `feat/short-description` or `fix/short-description`.
- Keep branches short-lived. Rebase on master before merging.

## Before committing

- Run `go vet ./...` — catches common mistakes.
- Run `golangci-lint run` — enforces project lint rules.
- Run `go test ./...` — all tests must pass.
- Run `go mod tidy` — clean up unused dependencies.
- Check `git diff --staged` — review what you're actually committing.

## Before pushing

- Squash WIP commits into meaningful commits.
- Ensure CI will pass: lint + test + build.
- Never force-push to shared branches.
