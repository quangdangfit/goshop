# Go Shop

An example of gin contains many useful features for e-commerce websites

## How to run

### Required Environment

- Postgres
- Redis

You can see the docker compose file [here](https://github.com/quangdangfit/docker-compose-template/blob/master/base/docker-compose.yml) to setup required environment

### Config
- Copy config file: `cp config/config.sample.yaml config/config.yaml`
- You should modify `config/config.yaml`

```yaml
environment: development
database_uri: postgres://username:password@host:5432/database

redis:
  enable: true
  host: localhost
  port: 6397
  password:
  database: 0

cache:
  enable: true
  expiry_time: 3600
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
[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:	export GIN_MODE=release
 - using code:	gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /swagger/*any             --> github.com/swaggo/gin-swagger.CustomWrapHandler.func1 (1 handlers)
[GIN-debug] POST   /auth/auth/register       --> goshop/app/api.(*User).Register-fm (1 handlers)
[GIN-debug] POST   /auth/auth/login          --> goshop/app/api.(*User).Login-fm (1 handlers)
[GIN-debug] POST   /admin/roles              --> goshop/app/api.(*Role).CreateRole-fm (1 handlers)
[GIN-debug] GET    /api/v1/users/:uuid       --> goshop/app/api.(*User).GetUserByID-fm (1 handlers)
[GIN-debug] GET    /api/v1/products          --> goshop/app/api.(*Product).GetProducts-fm (1 handlers)
[GIN-debug] POST   /api/v1/products          --> goshop/app/api.(*Product).CreateProduct-fm (1 handlers)
[GIN-debug] GET    /api/v1/products/:uuid    --> goshop/app/api.(*Product).GetProductByID-fm (1 handlers)
[GIN-debug] PUT    /api/v1/products/:uuid    --> goshop/app/api.(*Product).UpdateProduct-fm (1 handlers)
[GIN-debug] GET    /api/v1/categories        --> goshop/app/api.(*Category).GetCategories-fm (1 handlers)
[GIN-debug] POST   /api/v1/categories        --> goshop/app/api.(*Category).CreateCategory-fm (1 handlers)
[GIN-debug] GET    /api/v1/categories/:uuid  --> goshop/app/api.(*Category).GetCategoryByID-fm (1 handlers)
[GIN-debug] GET    /api/v1/categories/:uuid/products --> goshop/app/api.(*Product).GetProductByCategoryID-fm (1 handlers)
[GIN-debug] PUT    /api/v1/categories/:uuid  --> goshop/app/api.(*Category).UpdateCategory-fm (1 handlers)
[GIN-debug] GET    /api/v1/warehouses        --> goshop/app/api.(*Warehouse).GetWarehouses-fm (1 handlers)
[GIN-debug] POST   /api/v1/warehouses        --> goshop/app/api.(*Warehouse).CreateWarehouse-fm (1 handlers)
[GIN-debug] GET    /api/v1/warehouses/:uuid  --> goshop/app/api.(*Warehouse).GetWarehouseByID-fm (1 handlers)
[GIN-debug] PUT    /api/v1/warehouses/:uuid  --> goshop/app/api.(*Warehouse).UpdateWarehouse-fm (1 handlers)
[GIN-debug] GET    /api/v1/quantities        --> goshop/app/api.(*Quantity).GetQuantities-fm (1 handlers)
[GIN-debug] POST   /api/v1/quantities        --> goshop/app/api.(*Quantity).CreateQuantity-fm (1 handlers)
[GIN-debug] GET    /api/v1/quantities/:uuid  --> goshop/app/api.(*Quantity).GetQuantityByID-fm (1 handlers)
[GIN-debug] PUT    /api/v1/quantities/:uuid  --> goshop/app/api.(*Quantity).UpdateQuantity-fm (1 handlers)
[GIN-debug] GET    /api/v1/orders            --> goshop/app/api.(*Order).GetOrders-fm (1 handlers)
[GIN-debug] POST   /api/v1/orders            --> goshop/app/api.(*Order).CreateOrder-fm (1 handlers)
[GIN-debug] GET    /api/v1/orders/:uuid      --> goshop/app/api.(*Order).GetOrderByID-fm (1 handlers)
[GIN-debug] PUT    /api/v1/orders/:uuid      --> goshop/app/api.(*Order).UpdateOrder-fm (1 handlers)
```

### Document
* API document at: `http://localhost:8888/swagger/index.html`

### Techstack
- RESTful API
- Gorm
- Swagger
- Logging
- Jwt-Go
- Gin
- Graceful restart or stop (fvbock/endless)
- Cron Job
- Redis
- Dig (Dependency Injection)