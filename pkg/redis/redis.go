package redis

import (
	"context"
	"encoding/json"
	"time"

	goredis "github.com/go-redis/redis/v8"

	"github.com/quangdangfit/gocommon/logger"
)

const (
	Timeout = 1
)

// IRedis interface
//
//go:generate mockery --name=IRedis
type IRedis interface {
	IsConnected() bool
	Get(key string, value interface{}) error
	Set(key string, value interface{}) error
	SetWithExpiration(key string, value interface{}, expiration time.Duration) error
	Remove(keys ...string) error
	Keys(pattern string) ([]string, error)
	RemovePattern(pattern string) error
}

// Config redis
type Config struct {
	Address  string
	Password string
	Database int
}

type redis struct {
	cmd goredis.Cmdable
}

// New Redis interface with config
func New(config Config) IRedis {
	ctx, cancel := context.WithTimeout(context.Background(), Timeout*time.Second)
	defer cancel()

	rdb := goredis.NewClient(&goredis.Options{
		Addr:     config.Address,
		Password: config.Password,
		DB:       config.Database,
	})

	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		logger.Fatal(pong, err)
		return nil
	}

	return &redis{
		cmd: rdb,
	}
}

func (r *redis) IsConnected() bool {
	ctx, cancel := context.WithTimeout(context.Background(), Timeout*time.Second)
	defer cancel()

	if r.cmd == nil {
		return false
	}

	_, err := r.cmd.Ping(ctx).Result()
	if err != nil {
		return false
	}
	return true
}

func (r *redis) Get(key string, value interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), Timeout*time.Second)
	defer cancel()

	strValue, err := r.cmd.Get(ctx, key).Result()
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(strValue), value)
	if err != nil {
		return err
	}

	return nil
}

func (r *redis) SetWithExpiration(key string, value interface{}, expiration time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), Timeout*time.Second)
	defer cancel()

	bData, _ := json.Marshal(value)
	err := r.cmd.Set(ctx, key, bData, expiration).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *redis) Set(key string, value interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), Timeout*time.Second)
	defer cancel()

	bData, _ := json.Marshal(value)
	err := r.cmd.Set(ctx, key, bData, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *redis) Remove(keys ...string) error {
	ctx, cancel := context.WithTimeout(context.Background(), Timeout*time.Second)
	defer cancel()

	err := r.cmd.Del(ctx, keys...).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *redis) Keys(pattern string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), Timeout*time.Second)
	defer cancel()

	keys, err := r.cmd.Keys(ctx, pattern).Result()
	if err != nil {
		return nil, err
	}

	return keys, nil
}

func (r *redis) RemovePattern(pattern string) error {
	keys, err := r.Keys(pattern)
	if err != nil {
		return err
	}

	if len(keys) == 0 {
		return nil
	}

	err = r.Remove(keys...)
	if err != nil {
		return err
	}

	return nil
}
