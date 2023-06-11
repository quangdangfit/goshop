# Go Shop
[![Master](https://github.com/quangdangfit/goshop/workflows/CI%20GoShop/badge.svg)](https://github.com/quangdangfit/goshop/actions)

An example of gin contains many useful features for e-commerce websites

## How to run

### Required Environment

- Postgres
- Redis

You can see the docker compose file [here](https://github.com/quangdangfit/docker-compose-template/blob/master/base/docker-compose.yml) to set up required environment

### Config
- Copy config file: `cp config/config.sample.yaml config/config.yaml`
- You should modify `config/config.yaml`

```yaml
environment: production
port: 8888
auth_secret: ######
database_uri: postgres://username:password@host:5432/database
redis_uri: localhost:6379
redis_password:
redis_db: 0
```

### Run
```shell script
$ go run main.go 
```

### Test
```shell script
$ go test
```

### Test with Coverage
```shell script
go test -timeout 9000s -a -v -coverpkg=./... ./test
```

Project information and existing API

```
[GIN-debug] GET    /swagger/*any             --> github.com/swaggo/gin-swagger.CustomWrapHandler.func1 (3 handlers)
[GIN-debug] POST   /auth/register            --> goshop/app/api.(*UserAPI).Register-fm (3 handlers)
[GIN-debug] POST   /auth/login               --> goshop/app/api.(*UserAPI).Login-fm (3 handlers)
[GIN-debug] POST   /auth/refresh             --> goshop/app/api.(*UserAPI).RefreshToken-fm (4 handlers)
[GIN-debug] GET    /auth/me                  --> goshop/app/api.(*UserAPI).GetMe-fm (4 handlers)
[GIN-debug] PUT    /auth/change-password     --> goshop/app/api.(*UserAPI).ChangePassword-fm (4 handlers)
[GIN-debug] GET    /api/v1/products          --> goshop/app/api.(*ProductAPI).ListProducts-fm (3 handlers)
[GIN-debug] POST   /api/v1/products          --> goshop/app/api.(*ProductAPI).CreateProduct-fm (4 handlers)
[GIN-debug] PUT    /api/v1/products/:id      --> goshop/app/api.(*ProductAPI).UpdateProduct-fm (4 handlers)
[GIN-debug] GET    /api/v1/products/:id      --> goshop/app/api.(*ProductAPI).GetProductByID-fm (3 handlers)
[GIN-debug] POST   /api/v1/orders            --> goshop/app/api.(*OrderAPI).PlaceOrder-fm (4 handlers)
[GIN-debug] GET    /api/v1/orders/:id        --> goshop/app/api.(*OrderAPI).GetOrderByID-fm (4 handlers)
[GIN-debug] GET    /api/v1/orders            --> goshop/app/api.(*OrderAPI).GetOrders-fm (4 handlers)
[GIN-debug] PUT    /api/v1/orders/:id/cancel --> goshop/app/api.(*OrderAPI).CancelOrder-fm (4 handlers)
2023-06-11T13:31:47.587+0700    INFO    goshop/main.go:34       Listen at: 8888
```

### Document
* API document at: `http://localhost:8888/swagger/index.html`

### Tech stack
- Restful API
- Gorm
- Swagger
- Logging
- Jwt-Go
- Gin-gonic
- Redis
- Dig (Dependency Injection)