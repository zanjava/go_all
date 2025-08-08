package concurrence

import (
	"fmt"
	"sync"
	"time"
)

// 5种阻塞
// 1. 等待时间到  time.Sleep(time.Second)
// 2. 等待条件满足  sync.Cond
// 3. 等待channel数据  <-ch
// 4. 等待Mutex解锁  mu.Lock()
// 5. 等待select分支就绪
func Block() {
	// 1. 等待时间到
	time.Sleep(time.Second) //时间到了，会自动解除阻塞
	fmt.Println("sleep over")

	// 2.waiting condition
	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
	fmt.Println("wait over")

	// 3. 等待channel数据
	ch := make(chan bool, 10)
	<-ch
	fmt.Println("receive channel over")
	// for ele := range ch {
	// 	fmt.Println(ele)
	// }
	// fmt.Println("traverse channel over")

	// 4. 等待Mutex解锁
	mu := sync.Mutex{}
	mu.Lock()
	mu.Lock()
	fmt.Println("got lock")

	// 5. 等待select分支就绪
	select {}

}
