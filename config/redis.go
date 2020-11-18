package config

import (
	"github.com/go-redis/redis"
)

// GetRedis 获取redis连接
func GetRedis() *redis.Client {
	client := redis.NewClient(&redis.Options{})
	return client
}
