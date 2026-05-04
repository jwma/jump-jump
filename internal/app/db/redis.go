package db

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

func GetRedisClient() *redis.Client {
	return redisClient
}

func InitRedis() error {
	dbIdx, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		return fmt.Errorf("REDIS_DB env is required and must be an integer")
	}

	redisClient = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       dbIdx,
	})

	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		return fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return nil
}

func CloseRedis() {
	if redisClient != nil {
		redisClient.Close()
	}
}
