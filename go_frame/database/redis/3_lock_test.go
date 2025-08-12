package distributed_test

import (
	"fmt"
	distributed "go/frame/database/redis"
	"sync"
	"testing"
	"time"
)

func TestLock(t *testing.T) {
	LockName := "schedule"
	const C = 5
	wg := sync.WaitGroup{}
	wg.Add(C)
	for i := 0; i < C; i++ {
		go func(i int) { //用单机的多协程来模拟分布式环境下的多进程
			defer wg.Done()
			if distributed.TryLock(client, LockName, 10*time.Minute) {
				fmt.Printf("第%d个协程获得锁\n", i)
			}
		}(i)
	}
	wg.Wait()
	distributed.ReleaseLock(client, LockName)
}

// go test -v ./database/redis -run=^TestLock$ -count=1
