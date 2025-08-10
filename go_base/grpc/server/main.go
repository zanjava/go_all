package main

import (
	"context"
	"errors"
	"fmt"
	grpc_service "go/base/grpc/idl/service"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

func timer(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	begin := time.Now()
	resp, err = handler(ctx, req)
	fmt.Printf("use time %d ms\n", time.Since(begin).Milliseconds())
	return
}

func counter(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	fmt.Printf("%s add once\n", info.FullMethod)
	resp, err = handler(ctx, req)
	return
}

func devKey(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	meta, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("获取不到dev_key")
	}
	if value, exists := meta["dev_key"]; !exists {
		return nil, errors.New("获取不到dev_key")
	} else {
		if value[0] != "123456" { //对开发者key的合法性进行检查
			return nil, errors.New("dev_key不正确")
		} else {
			resp, err = handler(ctx, req)
			return
		}
	}
}

func main() {
	// 监听本地的5678端口
	lis, err := net.Listen("tcp", "127.0.0.1:5678")
	if err != nil {
		panic(err)
	}

	creds, err := credentials.NewServerTLSFromFile("../data/server.crt", "../data/rsa_private_key.pem")
	if err != nil {
		panic(err)
	}

	server := grpc.NewServer(
		grpc.Creds(creds),
		// grpc.UnaryInterceptor(timer),
		grpc.ChainUnaryInterceptor(timer, counter, devKey), // 链式拦截器
	)
	// 注册服务的具体实现，可以注册多个服务
	grpc_service.RegisterStudentServer(server, Student{})
	// 启动server
	err = server.Serve(lis)
	if err != nil {
		panic(err)
	}
}
