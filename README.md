# GoShop

[![CI](https://github.com/quangdangfit/goshop/workflows/master/badge.svg)](https://github.com/quangdangfit/goshop/actions)
[![codecov](https://codecov.io/gh/quangdangfit/goshop/graph/badge.svg?token=78BO8FQDB0)](https://codecov.io/gh/quangdangfit/goshop)
![Go Version](https://img.shields.io/github/go-mod/go-version/quangdangfit/goshop?style=flat-square)
[![License](https://img.shields.io/github/license/jrapoport/gothic?style=flat-square)](https://github.com/quangdangfit/goshop/blob/master/LICENSE)

A production-ready e-commerce backend built with Go, featuring a dual-server architecture that exposes both a REST API and a gRPC API from a single service.

## Architecture

The application runs two servers concurrently:

- **HTTP (REST)** — Gin framework, port `8888`
- **gRPC** — port `8889`, with JWT auth interceptor

Each domain (`user`, `product`, `order`, `cart`) follows a ports-and-adapters layout:

```
internal/{domain}/
├── model/       # GORM models
├── dto/         # Request/response structs with validation tags
├── repository/  # Database access (depends on dbs.Database interface)
├── service/     # Business logic (depends on repository interfaces)
└── port/
    ├── http/    # Gin handlers and route registration
    └── grpc/    # gRPC handlers and server registration
```

| Domain | HTTP | gRPC |
|--------|------|------|
| user | ✓ | ✓ |
| product | ✓ | ✓ |
| order | ✓ | ✓ |
| cart | ✓ | ✓ |

## Tech Stack

| Concern | Library |
|---------|---------|
| HTTP framework | [Gin](https://github.com/gin-gonic/gin) |
| gRPC | [grpc-go](https://github.com/grpc/grpc-go) |
| ORM | [GORM](https://gorm.io) + PostgreSQL |
| Cache | [go-redis](https://github.com/go-redis/redis) |
| Auth | JWT ([golang-jwt](https://github.com/golang-jwt/jwt)) |
| Validation | [gocommon/validation](https://github.com/quangdangfit/gocommon) |
| API Docs | [Swagger](https://github.com/swaggo/swag) |
| Testing | [testify](https://github.com/stretchr/testify) + [mockery](https://github.com/vektra/mockery) |
| Proto codegen | [buf](https://buf.build) |

## Prerequisites

- Go 1.24+
- PostgreSQL
- Redis

Docker Compose for local dependencies: [docker-compose-template](https://github.com/quangdangfit/docker-compose-template/blob/master/base/docker-compose.yml)

## Getting Started

**1. Clone and configure**

```bash
git clone https://github.com/quangdangfit/goshop.git
cd goshop
cp pkg/config/config.sample.yaml pkg/config/config.yaml
```

Edit `pkg/config/config.yaml`:

```yaml
environment: production
http_port: 8888
grpc_port: 8889
auth_secret: your-secret-key
database_uri: postgres://username:password@localhost:5432/goshop
redis_uri: localhost:6379
redis_password:
redis_db: 0
```

**2. Run**

```bash
go run cmd/api/main.go
```

```
INFO    HTTP server is listening on PORT: 8888
INFO    GRPC server is listening on PORT: 8889
```

**3. Browse the API**

Swagger UI: [http://localhost:8888/swagger/index.html](http://localhost:8888/swagger/index.html)

## API Reference

### Auth
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/auth/register` | Register |
| POST | `/api/v1/auth/login` | Login |
| POST | `/api/v1/auth/refresh` | Refresh access token |

### Users
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/users/me` | Get current user |
| PUT | `/api/v1/users/change-password` | Change password |

### Products
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/products` | List products (cached) |
| GET | `/api/v1/products/:id` | Get product (cached) |
| POST | `/api/v1/products` | Create product |
| PUT | `/api/v1/products/:id` | Update product |

### Orders
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/orders` | Place order |
| GET | `/api/v1/orders` | List my orders |
| GET | `/api/v1/orders/:id` | Get order details |
| PUT | `/api/v1/orders/:id/cancel` | Cancel order |

### Cart
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/cart` | Get my cart |
| POST | `/api/v1/cart` | Add product to cart |
| DELETE | `/api/v1/cart/:productId` | Remove product from cart |

## Development

**Run all tests with coverage**

```bash
make unittest
```

**Run a single test suite**

```bash
go test ./internal/product/service/... -v -run TestProductServiceTestSuite
```

**Run a single test case**

```bash
go test ./internal/product/service/... -v -run TestProductServiceTestSuite/TestCreateSuccess
```

**Regenerate mocks**

```bash
make mock
```

**Regenerate Swagger docs**

```bash
make doc
```

**Regenerate proto (Uses https://buf.build)**

```bash
cd proto && make build
```
