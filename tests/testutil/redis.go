package testutil

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"goshop/pkg/redis"
)

// StartRedis boots a throwaway redis:alpine container and returns a wired redis.Redis.
// Cleanup runs via t.Cleanup.
func StartRedis(ctx context.Context, t *testing.T) redis.Redis {
	t.Helper()
	r, terminate, err := startRedisContainer(ctx)
	if err != nil {
		t.Fatalf("start redis container: %v", err)
	}
	t.Cleanup(func() {
		shutdown, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		_ = terminate(shutdown)
	})
	return r
}

// StartRedisM is the TestMain-friendly variant. The returned cleanup must be invoked
// before os.Exit.
func StartRedisM(ctx context.Context) (redis.Redis, func(), error) {
	r, terminate, err := startRedisContainer(ctx)
	if err != nil {
		return nil, nil, err
	}
	cleanup := func() {
		shutdown, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		_ = terminate(shutdown)
	}
	return r, cleanup, nil
}

func startRedisContainer(ctx context.Context) (redis.Redis, func(context.Context, ...testcontainers.TerminateOption) error, error) {
	req := testcontainers.ContainerRequest{
		Image:        "redis:alpine",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForListeningPort("6379/tcp").WithStartupTimeout(30 * time.Second),
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("redis container: %w", err)
	}

	host, err := container.Host(ctx)
	if err != nil {
		_ = container.Terminate(ctx)
		return nil, nil, err
	}
	port, err := container.MappedPort(ctx, "6379/tcp")
	if err != nil {
		_ = container.Terminate(ctx)
		return nil, nil, err
	}

	r := redis.New(redis.Config{
		Address:  fmt.Sprintf("%s:%s", host, port.Port()),
		Password: "",
		Database: 0,
	})
	return r, container.Terminate, nil
}
