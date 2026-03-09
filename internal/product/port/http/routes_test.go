package http

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/quangdangfit/gocommon/validation"

	dbMocks "goshop/pkg/dbs/mocks"
	redisMocks "goshop/pkg/redis/mocks"
)

func TestRoutes(t *testing.T) {
	mockDB := dbMocks.NewDatabase(t)
	mockRedis := redisMocks.NewRedis(t)
	Routes(gin.New().Group("/"), mockDB, validation.New(), mockRedis)
}
