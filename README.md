# Go Shop
[![Master](https://github.com/quangdangfit/goshop/workflows/master/badge.svg)](https://github.com/quangdangfit/goshop/actions)
[![Codecov](https://img.shields.io/codecov/c/github/quangdangfit/goshop?style=flat-square)](https://codecov.io/gh/quangdangfit/goshop)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/quangdangfit/goshop?style=flat-square)
[![GitHub](https://img.shields.io/github/license/jrapoport/gothic?style=flat-square)](https://github.com/quangdangfit/goshop/blob/master/LICENSE)

An example of gin contains many useful features for e-commerce websites

## How to run

### Required Environment

- Postgres
- Redis

You can see the docker compose file [here](https://github.com/quangdangfit/docker-compose-template/blob/master/base/docker-compose.yml) to set up required environment

### Config
- Copy config file: `cp pkg/config/config.sample.yaml pkg/config/config.yaml`
- You should modify `pkg/config/config.yaml`

```yaml
environment: production
http_port: 8888
grpc_port: 8889
auth_secret: ######
database_uri: postgres://username:password@host:5432/database
redis_uri: localhost:6379
redis_password:
redis_db: 0
```

### Run
```shell script
$ go run cmd/api/main.go 
```

### Test
```shell script
$ go test
```

### Test with Coverage
```shell script
go test -timeout 9000s -a -v -coverprofile=coverage.out -coverpkg=./... ./...
```

**or**

```shell script
make unittest
```

Project information and existing API

```
[GIN-debug] POST   /api/v1/auth/register     --> goshop/internal/user/port/http.(*UserHandler).Register-fm (3 handlers)
[GIN-debug] POST   /api/v1/auth/login        --> goshop/internal/user/port/http.(*UserHandler).Login-fm (3 handlers)
[GIN-debug] POST   /api/v1/auth/refresh      --> goshop/internal/user/port/http.(*UserHandler).RefreshToken-fm (4 handlers)
[GIN-debug] GET    /api/v1/auth/me           --> goshop/internal/user/port/http.(*UserHandler).GetMe-fm (4 handlers)
[GIN-debug] PUT    /api/v1/auth/change-password --> goshop/internal/user/port/http.(*UserHandler).ChangePassword-fm (4 handlers)
[GIN-debug] GET    /api/v1/products          --> goshop/internal/product/port/http.(*ProductHandler).ListProducts-fm (3 handlers)
[GIN-debug] POST   /api/v1/products          --> goshop/internal/product/port/http.(*ProductHandler).CreateProduct-fm (4 handlers)
[GIN-debug] PUT    /api/v1/products/:id      --> goshop/internal/product/port/http.(*ProductHandler).UpdateProduct-fm (4 handlers)
[GIN-debug] GET    /api/v1/products/:id      --> goshop/internal/product/port/http.(*ProductHandler).GetProductByID-fm (3 handlers)
[GIN-debug] POST   /api/v1/orders            --> goshop/internal/order/port/http.(*OrderHandler).PlaceOrder-fm (4 handlers)
[GIN-debug] GET    /api/v1/orders/:id        --> goshop/internal/order/port/http.(*OrderHandler).GetOrderByID-fm (4 handlers)
[GIN-debug] GET    /api/v1/orders            --> goshop/internal/order/port/http.(*OrderHandler).GetOrders-fm (4 handlers)
[GIN-debug] PUT    /api/v1/orders/:id/cancel --> goshop/internal/order/port/http.(*OrderHandler).CancelOrder-fm (4 handlers)
[GIN-debug] GET    /swagger/*any             --> github.com/swaggo/gin-swagger.CustomWrapHandler.func1 (3 handlers)
[GIN-debug] GET    /health                   --> goshop/internal/server/http.Server.Run.func1 (3 handlers)
2023-08-20T13:50:57.175+0700    INFO    http/server.go:53       Server is listening on PORT: 8888
2023-09-08T21:03:00.950+0700    INFO    grpc/server.go:48       GRPC server is listening on PORT: 8889
[GIN-debug] Listening and serving HTTP on :8888
```

### Document
* API document at: `http://localhost:8888/swagger/index.html`

### Tech stack
- Restful API
- GRPC
- DDD
- Gorm
- Swagger
- Logging
- Jwt-Go
- Gin-gonic
- Redis

### What's next?
- gRPC functions for products and orders
- Push message to notify place order successfully
- Put database object into interface
- Unittest for repositories (mock database)
- Define error response wrapper
