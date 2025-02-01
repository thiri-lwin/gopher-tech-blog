package redis

import (
	"time"

	"github.com/go-redis/redis/v8"
)

var RDB *redis.Client

// Initialize Redis connection
func InitRedis(addr string, user string, pass string) {
	RDB = redis.NewClient(&redis.Options{
		Addr:        addr,
		Username:    user,
		Password:    pass,
		DialTimeout: 30 * time.Second,
	})
}
