package distributed

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// TryLock 尝试获得分布式锁，成功返回true，失败返回false
func TryLock(rc *redis.Client, key string, expire time.Duration) bool {
	cmd := rc.SetNX(context.Background(), key, "value随意", expire) //SetNX如果key不存在则返回true，写入key，并设置过期时间
	if cmd.Err() != nil {
		return false
	} else {
		return cmd.Val()
	}
}

// ReleaseLock 释放分布式锁
func ReleaseLock(rc *redis.Client, key string) {
	rc.Del(context.Background(), key)
}
