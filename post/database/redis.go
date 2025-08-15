package database

import (
	"context"
	"log/slog"
	"sync"

	"github.com/redis/go-redis/v9"
)

var (
	client    *redis.Client
	redisOnce sync.Once
)

func GetRedisClient() *redis.Client {
	redisOnce.Do(func() {
		client = redis.NewClient(&redis.Options{
			Addr: "127.0.0.1:6379",
			DB:   0, //redis默认会创建0-15号DB，这里使用默认的DB
		})
		//能ping成功才说明连接成功
		if err := client.Ping(context.Background()).Err(); err != nil {
			panic(err)
		} else {
			slog.Info("connect to redis")
		}
	})
	return client
}
