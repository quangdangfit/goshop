package http

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
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

func TestServer_GetEngine(t *testing.T) {
	mockDB := dbMocks.NewDatabase(t)
	mockRedis := redisMocks.NewRedis(t)

	server := NewServer(validation.New(), mockDB, mockRedis)
	assert.NotNil(t, server)

	engine := server.GetEngine()
	assert.NotNil(t, engine)
}

func TestServer_MapRoutes(t *testing.T) {
	mockDB := dbMocks.NewDatabase(t)
	mockRedis := redisMocks.NewRedis(t)

	server := NewServer(validation.New(), mockDB, mockRedis)
	assert.NotNil(t, server)

	err := server.MapRoutes()
	assert.Nil(t, err)
}

func TestServer_Run(t *testing.T) {
	mockDB := dbMocks.NewDatabase(t)
	mockRedis := redisMocks.NewRedis(t)

	// Find a free port
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	assert.NoError(t, err)
	port := ln.Addr().(*net.TCPAddr).Port
	_ = ln.Close()

	server := NewServer(validation.New(), mockDB, mockRedis)
	server.cfg.HttpPort = port
	server.cfg.Environment = config.ProductionEnv

	// Start server in background (goroutine intentionally leaks — acceptable in tests)
	go func() { _ = server.Run() }()

	// Wait for server to be ready
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		resp, err := http.Get(fmt.Sprintf("http://localhost:%d/health", port))
		if err == nil {
			_, _ = io.ReadAll(resp.Body)
			_ = resp.Body.Close()
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Error("server did not start in time")
}

func TestServer_Shutdown(t *testing.T) {
	mockDB := dbMocks.NewDatabase(t)
	mockRedis := redisMocks.NewRedis(t)

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	assert.NoError(t, err)
	port := ln.Addr().(*net.TCPAddr).Port
	_ = ln.Close()

	server := NewServer(validation.New(), mockDB, mockRedis)
	server.cfg.HttpPort = port

	runDone := make(chan struct{})
	go func() {
		_ = server.Run()
		close(runDone)
	}()

	// Wait until the server is serving.
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		resp, err := http.Get(fmt.Sprintf("http://localhost:%d/health", port))
		if err == nil {
			_, _ = io.ReadAll(resp.Body)
			_ = resp.Body.Close()
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	assert.NoError(t, server.Shutdown(ctx))

	select {
	case <-runDone:
	case <-time.After(2 * time.Second):
		t.Error("server Run did not return after Shutdown")
	}
}
