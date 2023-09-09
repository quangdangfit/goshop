package grpc

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
