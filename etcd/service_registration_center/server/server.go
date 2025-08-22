package main

import (
	"context"
	"flag"
	"fmt"
	"go/etcd/idl"
	"go/etcd/service_registration_center"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"google.golang.org/grpc"
)

var (
	servicePort = flag.Int("port", 0, "grpc service port") // server监听本地的端口
)

type MyServer struct {
	idl.UnimplementedHelloServiceServer
}

func (*MyServer) Login(ctx context.Context, request *idl.LoginRequest) (*idl.LoginResponse, error) {
	return nil, nil
}

func (*MyServer) SayHello(ctx context.Context, request *idl.HelloRequest) (*idl.HelloResponse, error) {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println(err)
		}
	}()
	time.Sleep(100 * time.Millisecond)
	resp := &idl.HelloResponse{Greeting: "hello"}
	return resp, nil
}

func main() {
	flag.Parse() //获取命令行参数（即servicePort）
	lis, err := net.Listen("tcp", "127.0.0.1:"+strconv.Itoa(*servicePort))
	if err != nil {
		panic(err)
	}
	server := grpc.NewServer()
	// 绑定服务的实现
	idl.RegisterHelloServiceServer(server, new(MyServer))

	// 向注册中心注册自己
	selfLocalIp, err := service_registration_center.GetLocalIP()
	if err != nil {
		panic(err)
	}
	var heartBeat int64 = 3
	hub := GetServiceHub(service_registration_center.ETCD_CLUSTER, heartBeat)
	leaseId, err := hub.Regist(service_registration_center.HELLO_SERVICE, selfLocalIp+":"+strconv.Itoa(*servicePort), 0)
	if err != nil {
		panic(err)
	}
	//周期性地注册自己（上报心跳）
	go func() {
		for {
			hub.Regist(service_registration_center.HELLO_SERVICE, selfLocalIp+":"+strconv.Itoa(*servicePort), leaseId)
			time.Sleep(time.Duration(heartBeat)*time.Second - 100*time.Millisecond)
		}
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM) //注册信号
		sig := <-c                                        //阻塞，直到信号的到来
		log.Printf("接收到信号 %s\n", sig.String())
		hub.UnRegist(service_registration_center.HELLO_SERVICE, selfLocalIp+":"+strconv.Itoa(*servicePort)) //主动注销服务。即使没有主动注销服务，最晚HeartBeat秒之后Hub也会发现这个endpoint挂掉了
		os.Exit(0)                                                                                          //进程退出
	}()

	// 启动服务
	err = server.Serve(lis)
	if err != nil {
		hub.UnRegist(service_registration_center.HELLO_SERVICE, selfLocalIp+":"+strconv.Itoa(*servicePort)) //如果服务启动失败，则注销自己
		panic(err)
	}
}

// go run .\etcd\service_registration_center\server\ -port=5678
// go run .\etcd\service_registration_center\server\ -port=5679
// go run .\etcd\service_registration_center\server\ -port=5680
