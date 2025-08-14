package main

import (
	"log/slog"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
)

func main6() {
	app := fiber.New()
	g0 := app.Group("/", M6) //全局中间件

	g1 := g0.Group("/v1", logger.New()) //组1的公共中间件
	g1.Get("/a", func(ctx fiber.Ctx) error {
		return ctx.SendString("name=dqq")
	})
	g1.Get("/b", func(ctx fiber.Ctx) error {
		return ctx.SendString("age=18")
	})

	g2 := g0.Group("/v2", logger.New()) //组2的公共中间件
	g2.Get("/a", func(ctx fiber.Ctx) error {
		return ctx.JSON(fiber.Map{"name": "dqq"})
	})
	g2.Get("/b", func(ctx fiber.Ctx) error {
		return ctx.JSON(fiber.Map{"age": 18})
	})

	if err := app.Listen("127.0.0.1:5678"); err != nil {
		slog.Error("fiber app start failed", "error", err)
	}
}
