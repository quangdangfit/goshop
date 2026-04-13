---
description: API design standards for REST (Gin) and gRPC endpoints
globs:
  - "**/port/http/**"
  - "**/port/grpc/**"
  - "proto/**"
alwaysApply: false
---

# API Design

## REST conventions

- Use RESTful resource naming: plural nouns (`/products`, `/orders`), not verbs.
- HTTP methods: `GET` (read), `POST` (create), `PUT` (full update), `PATCH` (partial update), `DELETE` (remove).
- Nested resources for belongs-to: `/products/:id/reviews`, not `/reviews?product_id=:id`.
- Response format must be consistent:
  - Success: `{"data": ..., "pagination": ...}` or `{"data": ...}`
  - Error: `{"error": {"code": "VALIDATION_ERROR", "message": "..."}}`
- Use proper HTTP status codes:
  - 200 OK, 201 Created, 204 No Content
  - 400 Bad Request, 401 Unauthorized, 403 Forbidden, 404 Not Found, 409 Conflict, 422 Unprocessable Entity
  - 500 Internal Server Error

## Request handling (Gin)

- Bind request body with `ShouldBindJSON` (returns error) not `BindJSON` (aborts on error).
- Validate DTOs with struct validation tags. Return 400 with specific field errors.
- Extract path params with `c.Param()`, query params with `c.Query()` or `c.DefaultQuery()`.
- Get authenticated user from context: `c.GetString("userId")`, `c.GetString("role")`.

## gRPC conventions

- Proto messages: use `snake_case` for fields, `PascalCase` for message/service names.
- Use `google.protobuf.Timestamp` for datetime fields, not string.
- Define proper error codes: `codes.NotFound`, `codes.InvalidArgument`, `codes.PermissionDenied`.
- Use interceptors for cross-cutting concerns (auth, logging, recovery).
- Keep proto files backward-compatible: never remove or renumber fields — mark them `reserved`.

## Pagination

- List endpoints accept `page` and `limit` query parameters.
- Default page: 1, default limit: 20, max limit: 100.
- Return pagination metadata: `{"total": N, "page": P, "limit": L, "total_pages": T}`.

## Versioning

- Prefix all routes with `/api/v1/`. New breaking changes go to `/api/v2/`.
- gRPC: version in package name (`package user.v1`).

## Documentation

- Annotate all HTTP handlers with Swagger comments for `swag init`.
- Include `@Summary`, `@Tags`, `@Accept`, `@Produce`, `@Param`, `@Success`, `@Failure`, `@Router`.
