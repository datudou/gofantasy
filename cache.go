package gofantasy

import (
	"context"
	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/redis/go-redis/v9"
)

type ICache interface {
	Get(ctx context.Context, key string) (value any, ok bool)
	Set(ctx context.Context, key string, value any)
}

type RedisCache struct {
	cache redis.Cmdable
}

func NewRedisCache(cmd redis.Cmdable) *RedisCache {
	return &RedisCache{
		cache: cmd,
	}
}

func (r *RedisCache) Get(ctx context.Context, key string) (value any, ok bool) {
	panic("implement me")
}

func (r *RedisCache) Set(key any, value any) {
	panic("implement me")
}

type LocalCache struct {
	cache *lru.Cache[string, any]
}

func (l *LocalCache) Get(ctx context.Context, key string) (value any, ok bool) {
	return l.cache.Get(key)
}

func (l *LocalCache) Set(ctx context.Context, key string, value any) {
	l.cache.Add(key, value)
}

func NewLocalCache(size int) *LocalCache {
	l, _ := lru.New[string, any](size)
	return &LocalCache{
		cache: l,
	}
}
