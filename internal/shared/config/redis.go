package config

import (
	"os"

	"github.com/redis/go-redis/v9"
)

var (
	RedisClient *redis.Client
)

func InitRedis() {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}
	password := os.Getenv("REDIS_PASSWORD")
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password, // add password from env
	})
}

func GetRedisClient() *redis.Client {
	return RedisClient
}

func CloseRedis() error {
	if RedisClient != nil {
		return RedisClient.Close()
	}
	return nil
} 