package cache

import (
	"context"
	"os"

	"github.com/redis/go-redis/v9"
)

var (
	RedisClient *redis.Client
	ctx         = context.Background()
)


func InitRedis() *redis.Client{
	RedisClient = redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_URL"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB: 0,
	})

	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		panic("Falha ao conectar com o Redis: " + err.Error())
	}

	return RedisClient
}
