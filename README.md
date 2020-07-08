# Go Shop

An example of gin contains many useful features for e-commerce websites

## Installation
```
$ go get gitlab.com/quangdangfit/goshop
```

## How to run

### Required

- Postgres
- Redis

### Config

You should modify `config/config.yaml`

```yaml
database:
  host: localhost
  port: 5432
  name: goshop
  env: development
  user: postgres
  password: 1234
  sslmode: disable
```

### Run
```shell script
$ cd $GOPATH/src/goshop
$ go run main.go 
```

Project information and existing API

```
[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:	export GIN_MODE=release
 - using code:	gin.SetMode(gin.ReleaseMode)

[GIN-debug] POST   /auth/register            --> goshop/objects/user.Service.Register-fm (3 handlers)
[GIN-debug] POST   /auth/login               --> goshop/objects/user.Service.Login-fm (3 handlers)
[GIN-debug] POST   /admin/roles              --> goshop/objects/role.Service.CreateRole-fm (4 handlers)
[GIN-debug] GET    /api/v1/users/:uuid       --> goshop/objects/user.Service.GetUserByID-fm (4 handlers)
[GIN-debug] GET    /api/v1/products          --> goshop/objects/product.Service.GetProducts-fm (4 handlers)
[GIN-debug] POST   /api/v1/products          --> goshop/objects/product.Service.CreateProduct-fm (4 handlers)
[GIN-debug] GET    /api/v1/products/:uuid    --> goshop/objects/product.Service.GetProductByID-fm (4 handlers)
[GIN-debug] PUT    /api/v1/products/:uuid    --> goshop/objects/product.Service.UpdateProduct-fm (4 handlers)
[GIN-debug] GET    /api/v1/categories        --> goshop/objects/category.Service.GetCategories-fm (4 handlers)
[GIN-debug] POST   /api/v1/categories        --> goshop/objects/category.Service.CreateCategory-fm (4 handlers)
[GIN-debug] GET    /api/v1/categories/:uuid  --> goshop/objects/category.Service.GetCategoryByID-fm (4 handlers)
[GIN-debug] GET    /api/v1/categories/:uuid/products --> goshop/objects/product.Service.GetProductByCategory-fm (4 handlers)
[GIN-debug] PUT    /api/v1/categories/:uuid  --> goshop/objects/category.Service.UpdateCategory-fm (4 handlers)

Listening port: 8888
```

### Contains
- RESTful API
- Gorm
- Swagger
- Logging
- Jwt-go
- Gin
- Graceful restart or stop (fvbock/endless)
- Cron
- Redis