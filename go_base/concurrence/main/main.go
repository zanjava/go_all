package main

import "go/base/concurrence"

func main() {
	// 调用 SimpleGoroutine 函数来演示协程的使用
	//concurrence.SimpleGoroutine()
	//concurrence.SubRoutine()
	//concurrence.WaitGroup()
	//concurrence.ServiceMain()

	//concurrence.CloseChannel()
	//concurrence.ChannelBlock()
	//concurrence.TraverseChannel()

	//concurrence.Block()

	//concurrence.Broadcast()
	//concurrence.CountDownLatch()

	//concurrence.DealMassFile("../../../data/biz_log")

	//concurrence.RoutineLimit()
	//concurrence.ListenMultiWay()
	concurrence.SelectBlock()

}
