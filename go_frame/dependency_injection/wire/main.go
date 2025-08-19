package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	server *http.Server
)

func ListenTermSignal(cleanup func()) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	sig := <-c
	slog.Info("receive term signal " + sig.String() + ", going to exit")

	//再执行cleanup函数
	cleanup()

	// 等Web Server完全终止
	if server != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		server.Shutdown(ctx) //Shutdown会结束Go进程
	}
}

func main() {
	h, cleanup, _ := InitHandler()
	go ListenTermSignal(cleanup)

	// 启动http server
	server = &http.Server{
		Addr:    "localhost:5678",
		Handler: h,
	}
	h.Route()
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		panic(err)
	}

}

// go run ./dependency_injection/wire
// 在浏览器里访问 http://localhost:5678/
