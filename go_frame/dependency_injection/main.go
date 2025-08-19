package main

import (
	"context"
	"go/frame/dependency_injection/database"
	gd "go/frame/dependency_injection/database/gorm"

	//xd "go/frame/dependency_injection/database/xorm"
	"go/frame/dependency_injection/handler"
	"go/frame/dependency_injection/service"
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

func ListenTermSignal(cleanups ...func()) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	sig := <-c
	slog.Info("receive term signal " + sig.String() + ", going to exit")

	//逆序执行各个cleanup函数
	for i := len(cleanups) - 1; i >= 0; i-- {
		cleanups[i]()
	}

	// 等Web Server完全终止
	if server != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		server.Shutdown(ctx) //Shutdown会结束Go进程
	}
}

func main() {
	// 指定每一层的依赖
	var nd database.INewsDB //dao层
	var cleanup1 func()
	db1, cleanup1 := gd.GetMysqlDB()
	nd = gd.NewGormNews(db1)
	// db2, cleanup1 := xd.GetMysqlDB()
	// nd = xd.NewXormNews(db2)

	var s service.INewsRecommender // service层
	s = service.NewRecByPoster(nd)
	s = service.NewRecByTime(nd)

	var h handler.IHandler // handler层
	var cleanup2 func()
	h, cleanup2 = handler.NewGinHandler(s)

	go ListenTermSignal(cleanup1, cleanup2)

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

// go run ./dependency_injection
// 在浏览器里访问 http://localhost:5678/
