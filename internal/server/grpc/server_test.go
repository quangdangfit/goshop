package grpc

import (
	"testing"
	"time"

	"github.com/quangdangfit/gocommon/logger"
	"github.com/quangdangfit/gocommon/validation"
	"github.com/stretchr/testify/assert"

	"goshop/pkg/config"
	dbMocks "goshop/pkg/dbs/mocks"
	redisMocks "goshop/pkg/redis/mocks"
)

func init() {
	logger.Initialize(config.ProductionEnv)
}

func TestNewServer(t *testing.T) {
	mockDB := dbMocks.NewDatabase(t)
	mockRedis := redisMocks.NewRedis(t)

	server := NewServer(validation.New(), mockDB, mockRedis)
	assert.NotNil(t, server)
}

func TestServer_Run(t *testing.T) {
	mockDB := dbMocks.NewDatabase(t)
	mockRedis := redisMocks.NewRedis(t)

	server := NewServer(validation.New(), mockDB, mockRedis)

	// Use port 0 so OS picks a free port
	server.cfg.GrpcPort = 0

	errCh := make(chan error, 1)
	go func() {
		errCh <- server.Run()
	}()

	// Let the server start
	time.Sleep(20 * time.Millisecond)

	// Stop the gRPC server
	server.engine.Stop()

	select {
	case err := <-errCh:
		assert.NoError(t, err)
	case <-time.After(time.Second):
		t.Error("server did not stop in time")
	}
}
