package config

import "github.com/go-redis/redis"

// GetRedis 获取redis连接
func GetRedis() *redis.Client {
	conf := Conf.RedisConfig
	client := redis.NewClient(&redis.Options{
		Addr:     conf.Host + ":" + conf.Port,
		Password: conf.Password,
		DB:       conf.Database,
	})
	return client
}
