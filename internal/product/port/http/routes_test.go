package http

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/quangdangfit/gocommon/validation"
	"github.com/stretchr/testify/mock"

	dbMocks "goshop/pkg/dbs/mocks"
	redisMocks "goshop/pkg/redis/mocks"
)

func TestRoutes(t *testing.T) {
	mockDB := dbMocks.NewIDatabase(t)
	mockDB.On("AutoMigrate", mock.Anything).Return(nil).Times(1)
	mockRedis := redisMocks.NewIRedis(t)
	Routes(gin.New().Group("/"), mockDB, validation.New(), mockRedis)
}
