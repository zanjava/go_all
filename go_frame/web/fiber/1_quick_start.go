package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
)

func homeHandler(ctx fiber.Ctx) error {
	fmt.Println("协议:", ctx.Protocol())       // HTTP/1.1
	fmt.Println("请求方法:", ctx.Method())       //GET
	fmt.Println("请求URL:", ctx.OriginalURL()) // /home
	fmt.Println("请求路径:", ctx.Path())         // /home
	fmt.Println("请求头:")
	for k, v := range ctx.GetReqHeaders() {
		fmt.Printf("%s=%s\n", k, v[0]) // User-Agent=Go-http-client/1.1
	}
	fmt.Println("请求体:", string(ctx.Request().Body()))
	fmt.Println("请求体:", string(ctx.Request().Body())) //可以多次使用请求体

	// 下面三项顺序没要求
	// 设置响应体
	// ctx.Response().BodyWriter().Write([]byte("welcome"))
	// ctx.Response().BodyWriter().Write([]byte(" to"))    //追加内容
	// ctx.Response().BodyWriter().Write([]byte(" China")) //追加内容
	// // 设置状态码
	// ctx.Response().SetStatusCode(http.StatusOK)
	// // 设置响应头
	// ctx.Response().Header.Add("language", "go")

	// fiber框架当然也封装了很多功能
	// 设置响应头
	ctx.Set("Set-Cookie", "language=go; Path=/; Domain=www.baidu.com; Max-Age=604800; HttpOnly")
	// 设置响应码（如果响应体为空，也会一并设置响应体）
	ctx.SendStatus(http.StatusOK)
	ctx.Status(http.StatusOK) //仅设置响应码
	// 设置响应体
	ctx.SendString("Welcome")      //设置响应体(注意：不是追加)
	ctx.WriteString(" to BeiJing") //向响应体里追加内容

	// ctx.JSON(map[string]any{"物理": 34, "化学": 70}) //设置响应体(注意：不是追加)
	// ctx.JSON(fiber.Map{"物理": 34, "化学": 70})      // type Map map[string]interface{}

	fmt.Println()
	return nil
}

func main1() {
	app := fiber.New()
	logFile, _ := os.OpenFile("../data/log/fiber.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
	app.Use(logger.New(
		logger.Config{
			TimeFormat:    "2006-01-02 - 15:04:05",
			TimeZone:      "Asia/Shanghai",
			DisableColors: true,
			Stream:        logFile, //日志输出到文件
		},
	))

	app.Use(recover.New())        //recover中间件，可以捕获nil pointer等运行时异常，触发ErrorHandler
	app.Get("/home", homeHandler) //既允许GET也允许HEAD
	app.Post("/home", homeHandler)
	app.Head("/ping", func(ctx fiber.Ctx) error {
		ctx.Status(200) //仅设置响应码
		return nil
	})
	if err := app.Listen("127.0.0.1:5678"); err != nil {
		slog.Error("fiber app start failed", "error", err)
	}
}

// go run ./web/fiber
