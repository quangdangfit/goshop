package gredis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"gitlab.com/quangdangfit/gocommon/utils/logger"
	"goshop/config"
	"time"
)

const (
	RedisExpiredTimes = 600
)

var ctx = context.Background()

type Redis interface {
	IsConnected() bool
	Get(key string) []byte
	Set(key string, val []byte) error
}

type gredis struct {
	client *redis.Client
}

// Setup Initialize the Redis instance
func NewRedis() Redis {
	redisConfig := config.Config.Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", redisConfig.Host, redisConfig.Port),
		Password: redisConfig.Password,
		DB:       redisConfig.Database,
	})

	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		logger.Error(pong, err)
		return nil
	}

	return &gredis{client: rdb}
}

func (g *gredis) IsConnected() bool {
	if g.client == nil {
		return false
	}

	_, err := g.client.Ping(ctx).Result()
	if err != nil {
		return false
	}
	return true
}

func (g *gredis) Get(key string) []byte {
	val, err := g.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil
	}
	if err != nil {
		logger.Info("Redis fail to get: ", err)
		return nil
	}
	logger.Debugf("Get from redis %s - %s", key, val)

	return val
}

func (g *gredis) Set(key string, val []byte) error {
	err := g.client.Set(ctx, key, val, RedisExpiredTimes*time.Second).Err()
	if err != nil {
		logger.Error("Redis fail to set: ", err)
		return err
	}
	logger.Debugf("Set to redis %s - %s", key, val)

	return nil
}

func (g *gredis) Remove(key string) error {
	return nil
}

var GRedis = NewRedis()
