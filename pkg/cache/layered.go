package cache

import "github.com/redis/go-redis/v9"



type LayeredCache struct {
	redis *redis.Client
}


