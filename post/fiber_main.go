package main

import (
	database "go/post/database/xorm"
	handler "go/post/handler/fiber"
	"go/post/util"
	"net/http"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/static"
	"github.com/gofiber/template/html/v2"
	"github.com/robfig/cron/v3"
)

func InitXorm() {
	util.InitSlog("./log/post.log")
	database.ConnectPostDB("./conf", "db", util.YAML, "./log")

	crontab := cron.New()
	crontab.AddFunc("*/30 * * * *", database.PingPostDB) // 分，时，日，月，星期。每隔30分钟ping一次数据库
	crontab.Start()
}

func main() {
	InitXorm()
	go ListenTermSignal() //监听信号

	engine := html.New("./views/html", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})
	app.Use("/css", static.New("./views/css")) //在url是访问目录/js相当于访问文件系统中的views/js目录
	app.Use("/js", static.New("./views/js"))
	app.Use("/favicon.ico", static.New("./views/img/dqq.png"))

	app.Get("/login", func(ctx fiber.Ctx) error { return ctx.Status(http.StatusOK).Render("login", nil) })
	app.Post("/login", func(ctx fiber.Ctx) error { return ctx.Status(http.StatusOK).Render("login", nil) }) // Auth失败重定向时可能会以post方式请求/login页面
	app.Get("/regist", func(ctx fiber.Ctx) error { return ctx.Status(http.StatusOK).Render("user_regist", nil) })
	app.Get("/modify_pass", handler.Auth, func(ctx fiber.Ctx) error { return ctx.Status(http.StatusOK).Render("update_pass", nil) })
	app.Post("/login/submit", handler.Login)
	app.Post("/regist/submit", handler.RegistUser)
	app.Post("/modify_pass/submit", handler.Auth, handler.UpdatePassword)
	app.Get("/user", handler.GetUserInfo)
	app.Get("/logout", handler.Logout)

	group := app.Group("/news")
	group.Get("", handler.NewsList)
	group.Get("/issue", func(ctx fiber.Ctx) error { return ctx.Status(http.StatusOK).Render("news_issue", nil) })
	group.Post("/issue/submit", handler.Auth, handler.PostNews)
	group.Get("/belong", handler.NewsBelong)
	group.Get("/:id", handler.GetNewsById)
	group.Get("/delete/:id", handler.Auth, handler.DeleteNews)
	group.Post("/update", handler.Auth, handler.UpdateNews)

	app.Get("", func(ctx fiber.Ctx) error { return ctx.Status(http.StatusMovedPermanently).Redirect().To("/news") }) //新闻列表页是默认的首页

	if err := app.Listen("127.0.0.1:5678"); err != nil {
		panic(err)
	}
}

// go run ./Post
// 在浏览器里访问 http://localhost:5678/
