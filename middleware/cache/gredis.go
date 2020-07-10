package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"gitlab.com/quangdangfit/gocommon/utils/logger"

	"goshop/config"
)

const (
	RedisExpiredTimes = 600
)

var ctx = context.Background()

type gredis struct {
	client *redis.Client
}

func NewRedis() *gredis {
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

func (g *gredis) Get(key string, data interface{}) error {
	val, err := g.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil
	}
	if err != nil {
		logger.Info("Cache fail to get: ", err)
		return nil
	}
	logger.Debugf("Get from redis %s - %s", key, val)

	err = json.Unmarshal(val, &data)
	if err != nil {
		return err
	}

	return nil
}

func (g *gredis) Set(key string, val []byte) error {
	err := g.client.Set(ctx, key, val, RedisExpiredTimes*time.Second).Err()
	if err != nil {
		logger.Error("Cache fail to set: ", err)
		return err
	}
	logger.Debugf("Set to redis %s - %s", key, val)

	return nil
}

func (g *gredis) Remove(key string) error {
	err := g.client.Del(ctx, key).Err()
	if err != nil {
		logger.Errorf("Cache fail to delete key %s: %s", key, err)
		return err
	}
	logger.Debug("Cache deleted key", key)

	return nil
}
