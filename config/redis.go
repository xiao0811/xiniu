package config

import (
	"github.com/go-redis/redis"
)

func GetRedis() *redis.Client {
	client := redis.NewClient(&redis.Options{})
	return client
}
