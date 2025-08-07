package basic

import (
	"fmt"
	"strings"
)

// channel是go语言里实现goroutine间通信的管道，可以在不同的goroutine间传递数据。channel是引用类型，声明时默认值为nil，必须初始化后才能使用。channel可以指定为只读或只写。
// channel的底层实现是一个环形队列，可以指定容量，容量满时再写入数据会阻塞，容量空时再读取数据也会阻塞。

func rangeChannel() {
	var ch chan int //声明
	if ch == nil {
		fmt.Printf("ch is nil,ch len %d cap %d\n", len(ch), cap(ch))
	}
	// ch <- 2 //不能向nil chan里发送数据
	if len(ch) == 0 { //引用类型未初始化时都是nil，可以对它们执行len()函数，返回0
		fmt.Println("ch length is 0")
	}
	ch = make(chan int, 8) //初始化，环形队列里可容纳8个int
	ch <- 1                //往管道里写入(send)数据
	ch <- 2
	ch <- 3
	ch <- 4
	ch <- 5
	fmt.Printf("ch len %d cap %d\n", len(ch), cap(ch))
	v := <-ch //从管道里取走(recv)数据
	fmt.Println(v)
	v = <-ch
	fmt.Println(v)
	fmt.Println()

	close(ch)
	//遍历并取走（receive）管道里的元素。当管道里已无剩余元素且没有close管道时，receive操作会一直阻塞，最终报deadlock。当管道为空且被close后，for循环退出。
	for ele := range ch {
		fmt.Println(ele)
	}

	c := make(chan int, 10)
	send(c)
	recv(c)
}

// 只能向channel里写数据
func send(c chan<- int) {
	c <- 1
}

// 只能取channel中的数据
func recv(c <-chan int) {
	v := <-c
	fmt.Printf("take %d from read-only channel\n", v)
}

// slice、map、channel是go语言里的三大引用类型，如果只是想改变它们引用的底层数据，不需要传指针，因为传引用类型本质上传的就是底层数据的指针

func changeArray(arr [3]int) {
	arr[0]++
}

func changeSlice(slc []int) {
	if len(slc) > 0 {
		slc[0]++
	}
}

func changeMap(mp map[int]bool) {
	mp[0] = true
}

func changeChan(ch chan bool) {
	if cap(ch) > len(ch) {
		ch <- true
	}
}

func main18() {
	rangeChannel()
	fmt.Println(strings.Repeat("*", 50)) //星号重复50次

	arr := [3]int{}
	changeArray(arr)
	fmt.Println(arr[0])

	slc := []int{1, 2, 3}
	changeSlice(slc)
	fmt.Println(slc[0])

	mp := map[int]bool{0: false}
	changeMap(mp)
	fmt.Println(mp[0])

	ch := make(chan bool, 3)
	changeChan(ch)
	fmt.Println(<-ch)
}
