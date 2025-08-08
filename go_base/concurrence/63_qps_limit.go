package concurrence

import (
	"sync"
	"time"
)

// 限流的通道  只有100个并发
var qps = make(chan struct{}, 100)

func handler() {
	qps <- struct{}{}
	defer func() {
		<-qps
	}()
	time.Sleep(3 * time.Second)
}

func QpsLimit() {
	const P = 1000 //模拟1000个请求
	wg := sync.WaitGroup{}
	wg.Add(P)
	for i := 0; i < P; i++ {
		go func() {
			defer wg.Done()
			handler()
		}()
	}
	wg.Wait()
}
