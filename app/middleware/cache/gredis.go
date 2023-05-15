package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/quangdangfit/gocommon/logger"

	"goshop/config"
)

const (
	RedisExpiredTimes = 600
)

var ctx = context.Background()

type gredis struct {
	client     *redis.Client
	expiryTime int
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
	expiryTime := config.Config.Cache.ExpiryTime
	if expiryTime <= 0 {
		expiryTime = RedisExpiredTimes
	}

	return &gredis{client: rdb, expiryTime: expiryTime}
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
	err := g.client.Set(ctx, key, val, time.Duration(g.expiryTime)*time.Second).Err()
	if err != nil {
		logger.Error("Cache fail to set: ", err)
		return err
	}
	logger.Debugf("Set to redis %s - %s", key, val)

	return nil
}

func (g *gredis) Remove(keys ...string) error {
	err := g.client.Del(ctx, keys...).Err()
	if err != nil {
		logger.Errorf("Cache fail to delete key %s: %s", keys, err)
		return err
	}
	logger.Debug("Cache deleted key", keys)

	return nil
}

func (g *gredis) Keys(pattern string) ([]string, error) {
	keys, err := g.client.Keys(ctx, pattern).Result()
	if err != nil {
		return nil, err
	}

	return keys, nil
}
