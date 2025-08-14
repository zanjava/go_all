package main

import (
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/limiter"
	"github.com/gofiber/fiber/v3/middleware/logger"
)

func bizHandler(ctx fiber.Ctx) error {
	ctx.WriteString("welcome\n")
	return nil
}

func M1(ctx fiber.Ctx) error {
	ctx.WriteString("M1 Begin\n")
	ctx.Next() // 先把后面的函数执行完，再进入下一行
	ctx.WriteString("M1 End\n")
	return nil
}

func M2(ctx fiber.Ctx) error {
	ctx.WriteString("Here is M2\n")
	return ctx.Next() //进入下一个中间件，必须显式调用Next()
}

func M3(ctx fiber.Ctx) error {
	ctx.WriteString("Here is M3\n")
	return ctx.Next()
}

func M4(ctx fiber.Ctx) error {
	ctx.WriteString("Here is M4\n")
	return nil
}

func M5(ctx fiber.Ctx) error {
	ctx.WriteString("Here is M5\n")
	return ctx.Next()
}

func M6(ctx fiber.Ctx) error {
	slog.Info("visit", "path", ctx.Path) //注意：ctx.Route().Path中注册路由(注册中间件)时使用的path，不是用户真实请求的path
	return ctx.Next()
}

func main3() {
	app := fiber.New()
	app.Use(logger.New())
	app.Use(M6)                           //全局中间件。访问根目录里会执行M6
	app.Use("/2", M2)                     //根据路径前缀添加中间件
	app.Use([]string{"/3", "/4"}, M2, M4) //指定多个路径前缀,指定多个中间件。同时也会定义路由，即访问/3或/4时，会执行M2, M4

	app.Get("/1", bizHandler, M1, M2, M3, M2, M4, M5) //中间件可以有多个，中间件都执行完之后才会执行bizHandler。如果某个中间件没有执行ctx.Next()，则之后的中间件和bizHandler都不会执行
	app.Get("/2", bizHandler, M3, M2, M4, M5)

	//固定窗口限流器
	limit := limiter.New(
		limiter.Config{
			// Next指明哪些情况下可跳过限流中间件
			Next: func(c fiber.Ctx) bool {
				return c.IP() == "127.0.0.2" //对方IP
			},
			//时间窗口
			Expiration: time.Minute,
			//时间窗口内最多允许请求几次
			Max: 5,
			LimitReached: func(c fiber.Ctx) error {
				return c.SendString("请求过于频繁")
			},
			// LimiterMiddleware:limiter.FixedWindow{},//固定窗口
			LimiterMiddleware: limiter.SlidingWindow{}, //滑动窗口
		},
	)

	app.Get("/lt", bizHandler, limit)

	if err := app.Listen("127.0.0.1:5678"); err != nil {
		slog.Error("fiber app start failed", "error", err)
	}
}
