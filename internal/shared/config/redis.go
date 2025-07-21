package config

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

func InitRedis() {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}
	password := os.Getenv("REDIS_PASSWORD")
	redisClient = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
	})
}

func SetOrderCache(ctx context.Context, orderID string, order interface{}, expiration time.Duration) error {
	data, err := json.Marshal(order)
	if err != nil {
		return err
	}
	return redisClient.Set(ctx, "order:"+orderID, data, expiration).Err()
}

func GetOrderCache(ctx context.Context, orderID string, dest interface{}) (bool, error) {
	val, err := redisClient.Get(ctx, "order:"+orderID).Result()
	if err == redis.Nil {
		return false, nil
	} else if err != nil {
		return false, err
	}
	if err := json.Unmarshal([]byte(val), dest); err != nil {
		return false, err
	}
	return true, nil
}

func DeleteOrderCache(ctx context.Context, orderID string) error {
	return redisClient.Del(ctx, "order:"+orderID).Err()
}
