package main

type cache interface {
	Get(key any) (value any, ok bool)
	Set(key any, value any)
}

type RedisCache struct {
	cmd *redis.Client
}
