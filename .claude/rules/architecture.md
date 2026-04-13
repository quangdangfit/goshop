---
description: Overall application architecture - dual server setup and domain structure
globs:
  - "cmd/**"
  - "internal/**"
alwaysApply: false
---

# Architecture

The app runs **two servers concurrently** from `cmd/api/main.go`: an HTTP server (Gin, REST) and a gRPC server. Both share the same `Database` and `Redis` instances.

## Domain structure

Each domain (`user`, `product`, `order`, `cart`) follows this layout:

```
internal/{domain}/
  model/      — GORM models; BeforeCreate hooks generate UUID IDs and codes
  dto/        — request/response structs (validated via gocommon/validation tags)
  repository/ — DB access only; depends on dbs.Database
  service/    — business logic; depends on repository interfaces
  port/
    http/     — Gin handlers + route registration
    grpc/     — gRPC handlers + server registration
```

## Which domains expose which transport

- Both HTTP and gRPC: `user`, `product`, `order`
- gRPC only: `cart`

## Interface naming convention

Service and repository interfaces use plain names without an `I` prefix (e.g. `ProductService`, `ProductRepository`). Concrete implementation structs are unexported camelCase (e.g. `productSvc`, `productRepo`). Generated mock struct names match the interface name (e.g. `mocks.ProductService`, constructor `mocks.NewProductService`).
