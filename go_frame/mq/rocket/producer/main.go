package main

func main() {
	//Sync()
	//Async()
	Delay()
	//Fifo()
	// Transaction()

	// 终止Producer
	producer.GracefulStop()
}

// go run ./mq/rocket/producer
