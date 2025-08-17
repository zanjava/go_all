package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

func InitLogger() {
	logFile, _ := os.OpenFile("log/im_client.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	log.SetOutput(logFile)
}

func ListenTermSignal() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	sig := <-c
	log.Println("receive term signal " + sig.String() + ", going to exit")

	ReleaseWebSocket()
	GetRabbitMQ().Release()
	CloseShell()

	os.Exit(0)
}

func main() {
	InitLogger()
	go ListenTermSignal()
	InitGroupBelong()
	InitApp()
	app.Run()
}

// 如果是Windows系统，最好在wsl里运行
// go run ./im/im_based_rabbitmq/client
// dqq-im » login 1
// dqq-im » send -t g1 -m 我是u1我来了
// Ctrl+A
// dqq-im » exit
