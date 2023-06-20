package configs

import "github.com/go-redis/redis"

// Client is setting connection with redis
var RedisClient redis.Client

func InitRedis()  {
	RedisClient = *redis.NewClient(&redis.Options{
		Addr: AppConfig.RedisAddress,
		Password: AppConfig.RedisPassword,
		DB: 0, // AppConfig.RedisDatabase,
	})
}