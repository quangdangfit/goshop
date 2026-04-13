---
description: Key shared packages in pkg/ - database, redis, JWT, middleware, utils, paging
globs:
  - "pkg/**"
alwaysApply: false
---

# Key shared packages

- `pkg/dbs` — PostgreSQL wrapper around GORM. `Database` interface is what all repositories depend on. `dbs.WithQuery`, `dbs.WithPreload`, `dbs.WithLimit`, `dbs.WithOffset`, `dbs.WithOrder` are the option helpers passed to `Find`/`FindOne`/`Count`.
- `pkg/redis` — Redis wrapper. `Redis` interface. Product list/get responses are cached by request URI and invalidated with `RemovePattern("*product*")` on write.
- `pkg/jtoken` — JWT generation (`GenerateAccessToken`, `GenerateRefreshToken`) and `ValidateToken`. Tokens embed `id`, `email`, `role`, and `type` in a `payload` claim.
- `pkg/middleware` — `JWTAuth()` / `JWTRefresh()` for HTTP (extracts `userId` and `role` into Gin context); `AuthInterceptor` for gRPC (skips `/user.UserService/Login` and `/user.UserService/Register`).
- `pkg/utils` — `utils.Copy(dst, src)` uses JSON marshal/unmarshal to map between types (model↔DTO). Used everywhere instead of manual field mapping.
- `pkg/paging` — `paging.New(page, limit, total)` returns a `Pagination` struct; default and max page size is 20.
