package basic

import (
	"fmt"
	"os"
	"runtime/debug"
	"runtime/trace"
	"sync"
	"time"
)

var arr [1000]int                 //数组长度固定，通常在栈上
var slc []int = make([]int, 1000) //动态创建的结构体、切片、map在堆上。引用类型的全局变量内存分配在堆上，值类型的全局变量分配在栈上

func stack_heap() {
	var brr [1000]byte //函数的入参、出参、局部变量一般在栈上

	var crr [128 << 10]byte
	var drr [128<<10 + 1]byte //数组超过128K就会分配到堆上(moved to heap)

	err := make([]byte, 64<<10)   //函数的入参、出参、局部变量一般在栈上
	frr := make([]byte, 64<<10+1) //切片超过64K就会分配到堆上(escapes to heap)

	_ = arr
	_ = brr
	_ = crr
	_ = drr
	_ = err
	_ = frr
	_ = slc
}

const (
	NumWorkers    = 4     // Number of workers.
	NumTasks      = 500   // Number of tasks.
	MemoryIntense = 10000 // Size of memory-intensive task (number of elements).
)

// 所有的垃圾回收都是针对堆的
func Gc() {
	// Write to the trace file.
	f, err := os.Create("../data/trace.out")
	if err != nil {
		fmt.Println("Failed to create trace file:", err)
		return
	}
	trace.Start(f)
	defer trace.Stop()

	// Set the target percentage for the garbage collector. Default is 100%.
	debug.SetGCPercent(100) //触顶时再执行GC

	// Task queue and result queue.
	taskQueue := make(chan int, NumTasks)
	resultQueue := make(chan int, NumTasks)

	// Start workers.
	var wg sync.WaitGroup
	wg.Add(NumWorkers)
	for i := 0; i < NumWorkers; i++ {
		go worker(taskQueue, resultQueue, &wg)
	}

	// Send tasks to the queue.
	for i := 0; i < NumTasks; i++ {
		taskQueue <- i
	}
	close(taskQueue)

	// Retrieve results from the queue.
	go func() {
		wg.Wait()
		close(resultQueue)
	}()

	// Process the results.
	for result := range resultQueue {
		fmt.Println("Result:", result)
	}

	fmt.Println("Done!")
}

// Worker function.
func worker(tasks <-chan int, results chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()

	for task := range tasks {
		result := performMemoryIntensiveTask(task)
		results <- result
	}
}

// performMemoryIntensiveTask is a memory-intensive function.
func performMemoryIntensiveTask(task int) int {
	// Create a large-sized slice.
	data := make([]int, MemoryIntense) //申请一块大内存，虽是局部变量，但会逃逸到堆上，长时间累积会触发GC
	for i := 0; i < MemoryIntense; i++ {
		data[i] = i + task
	}

	// Latency imitation.
	time.Sleep(10 * time.Millisecond)

	// Calculate the result.
	result := 0
	for _, value := range data {
		result += value
	}
	return result
}

func main333() {
	stack_heap()
	// gc() // 程序运行完之后生成一个文件data/trace.out, 然后执行 go tool trace data/trace.out
}

// go run -gcflags=-m ./reflection/gc.go
// go run ./reflection/gc.go
// go tool trace data/trace.out

// go run -gcflags=-m
// go tool trace ../data/trace.out
