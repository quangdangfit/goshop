package http

import (
	"testing"

	"github.com/quangdangfit/gocommon/validation"
	"github.com/stretchr/testify/assert"

	dbMocks "goshop/pkg/dbs/mocks"
	redisMocks "goshop/pkg/redis/mocks"
)

func TestNewServer(t *testing.T) {
	mockDB := dbMocks.NewIDatabase(t)
	mockRedis := redisMocks.NewIRedis(t)

	server := NewServer(validation.New(), mockDB, mockRedis)
	assert.NotNil(t, server)
}

func TestServer_GetEngine(t *testing.T) {
	mockDB := dbMocks.NewIDatabase(t)
	mockRedis := redisMocks.NewIRedis(t)

	server := NewServer(validation.New(), mockDB, mockRedis)
	assert.NotNil(t, server)

	engine := server.GetEngine()
	assert.NotNil(t, engine)
}

func TestServer_MapRoutes(t *testing.T) {
	mockDB := dbMocks.NewIDatabase(t)
	mockRedis := redisMocks.NewIRedis(t)

	server := NewServer(validation.New(), mockDB, mockRedis)
	assert.NotNil(t, server)

	err := server.MapRoutes()
	assert.Nil(t, err)
}
