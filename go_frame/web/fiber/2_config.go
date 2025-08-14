package main

import (
	"os"
	"time"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v3"
	fiberlog "github.com/gofiber/fiber/v3/log"
)

func main2() {
	logger := fiberlog.DefaultLogger() //这个logger业务侧也可以用
	logger.SetLevel(fiberlog.LevelInfo)
	//日志控制
	logFile, _ := os.OpenFile("../data/log/fiber.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
	logger.SetOutput(logFile)

	logger.Info("start fiber")

	app := fiber.New(
		fiber.Config{
			AppName: "zgw",
			ErrorHandler: func(c fiber.Ctx, err error) error { //recover捕获运行时异常，或Handler返回error时，会触发ErrorHandler
				logger.Error(err)
				return nil
			},
			ServerHeader:      "Fiber",       //响应头里会带一个Server=Fiber
			BodyLimit:         4 << 20,       //请求体最大为多少
			StreamRequestBody: false,         //是否允许请求体为流（即多次发送请求体），如果有上传大文件的场景需要设置为true
			Concurrency:       256 << 10,     //最多同时开几个连接
			DisableKeepalive:  false,         //使用长连接，即server端返回response后，client还可以用这个连接再次发送request。http1.1默认使用的是长连接
			IdleTimeout:       time.Minute,   //启用Keepalive时，server端等待下一个request的超时时间
			ReadTimeout:       time.Minute,   //读取整个请求（包括请求体）的超时
			WriteTimeout:      time.Minute,   //返回整个响应的超时
			ReadBufferSize:    4 << 10,       //socket层一次read使用的缓冲大小，因为url和请求头是一次性全读出来的，所以如果你的url或请求头特别长（比如含有大cookie）请调大ReadBufferSize
			WriteBufferSize:   4 << 10,       //response数据攒够多大才发给对方（除非显式调用flush）
			JSONEncoder:       sonic.Marshal, //默认使用go标准库的json序列化，也可以任意指定第三方库，更高效
			JSONDecoder:       sonic.Unmarshal,
			// EnablePrintRoutes: true, //在日志里打印设置的路由信息，即：path、method、handler
		},
	)

	app.Hooks().OnRoute(func(r fiber.Route) error { //添加钩子，创建路由的时候执行该回调函数
		logger.Infof("route path %s method %s", r.Path, r.Method)
		return nil
	})

	app.Get("/home", func(c fiber.Ctx) error { //既允许GET也允许HEAD
		return c.SendString("这是一个安全的网站")
	})

	// 启动http server
	// if err := app.Listen("127.0.0.1:5678"); err != nil {
	// 	logger.Errorf("fiber app start failed:%s", err)
	// }

	// 启动https server。在浏览器的证书管理中，导致data/server.crt文件
	if err := app.Listen("127.0.0.1:5678", fiber.ListenConfig{CertFile: "../data/server.crt", CertKeyFile: "../data/rsa_private_key.pem"}); err != nil {
		logger.Errorf("fiber app start failed:%s", err)
	}
}
