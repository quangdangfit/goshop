package gredis

import (
	"context"
	"encoding/json"
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
	Get(key string) []byte
	Set(key string, val interface{}) error
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

func (g *gredis) Get(key string) []byte {
	val, err := g.client.Get(ctx, key).Bytes()
	if err != nil {
		logger.Info("Failed to get from redis: ", err)
		return nil
	}
	logger.Debugf("Get from redis %s - %s", key, val)

	return val
}

func (g *gredis) Set(key string, val interface{}) error {
	data, _ := json.Marshal(val)
	err := g.client.Set(ctx, key, data, RedisExpiredTimes*time.Second).Err()
	if err != nil {
		logger.Error("Failed to set to redis: ", err)
		return err
	}
	logger.Debugf("Set to redis %s - %s", key, val)

	return nil
}

var GRedis = NewRedis()
