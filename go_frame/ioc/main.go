package main

import (
	"context"
	"go/frame/ioc/handler"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/golobby/container/v3"
	"gorm.io/gorm"
)

var (
	h      handler.IHandler
	server *http.Server
)

func main() {
	err := container.Resolve(&h) //让IoC容器提供一个接口的具体实现
	if err != nil {
		panic(err)
	}

	go ListenTermSignal()

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

func ListenTermSignal() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	sig := <-c
	slog.Info("receive term signal " + sig.String() + ", going to exit")

	h.Close()
	var db *gorm.DB
	if err := container.Resolve(&db); err == nil {
		sqlDB, _ := db.DB()
		sqlDB.Close()
		log.Println("关闭gorm数据库连接")
	}

	// 等Web Server完全终止
	if server != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		server.Shutdown(ctx) //Shutdown会结束Go进程
	}
}

// go run ./ioc
// 在浏览器里访问 http://localhost:5678/
