package redis

import (
	"gofit-api/configs"
	"time"
)

func SetOTPtoRedis(key string, code int) error {
	expirationTime := 10 * time.Minute
	err := configs.RedisClient.Set(key, code, expirationTime).Err()
	return err
}

func GetOTPfromRedis(key string) (string, error) {
	value, err := configs.RedisClient.Get(key).Result()
	if err != nil {
		return "", err
	}
	return value, nil
}
