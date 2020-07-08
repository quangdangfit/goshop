package repositories

import "goshop/caches/gredis"

type Repository struct{}

func (r *Repository) GetCache(key string) []byte {
	return gredis.GRedis.Get(key)
}

func (r *Repository) SetCache(key string, value interface{}) error {
	return gredis.GRedis.Set(key, value)
}
