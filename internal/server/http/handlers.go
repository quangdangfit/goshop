package http

import (
	"github.com/gin-gonic/gin"
	"github.com/quangdangfit/gocommon/redis"
	"github.com/quangdangfit/gocommon/validation"

	productHttp "goshop/internal/product/port/http"
	userHttp "goshop/internal/user/port/http"
)

func (s Server) MapRoutes(e *gin.Engine, validator validation.Validation, cache redis.IRedis) error {
	v1 := e.Group("/api/v1")
	userHttp.Routes(v1, validator)
	productHttp.Routes(v1, validator, cache)
	return nil
}
