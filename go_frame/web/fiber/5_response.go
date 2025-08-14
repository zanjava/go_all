package main

import (
	"log/slog"
	"net/http"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/static"
	"github.com/gofiber/template/html/v2"
)

func text(app *fiber.App) {
	app.Get("/user/text", func(ctx fiber.Ctx) error {
		return ctx.SendString("hi boy") //响应头 Content-Type=text/plain
	})
}

func json0(app *fiber.App) {
	app.Get("/user/json0", func(ctx fiber.Ctx) error {
		var stu struct { //匿名结构体
			Name    string `json:"name"`
			Address string `json:"addr"`
		}
		stu.Name = "zcy"
		stu.Address = "bj"
		s, _ := sonic.MarshalString(stu)
		ctx.Set("Content-Type", "application/json")
		ctx.SendStatus(http.StatusOK)
		return ctx.SendString(s)
	})
}

func json1(app *fiber.App) {
	app.Get("/user/json1", func(ctx fiber.Ctx) error {
		return ctx.JSON(fiber.Map{"name": "zcy", "addr": "bj"}) //type Map map[string]interface{}  //响应头 Content-Type=application/json
	})
}

func json2(app *fiber.App) {
	var stu struct { //匿名结构体
		Name    string `json:"name"`
		Address string `json:"addr"`
	}
	stu.Name = "zcy"
	stu.Address = "bj"
	app.Get("/user/json2", func(ctx fiber.Ctx) error {
		return ctx.JSON(stu) //响应头 Content-Type=application/json
	})
}

func jsonp(app *fiber.App) {
	var stu struct { //匿名结构体
		Name    string `json:"name"`
		Address string `json:"addr"`
	}
	stu.Name = "zcy"
	stu.Address = "bj"
	app.Get("/user/jsonp", func(ctx fiber.Ctx) error {
		return ctx.JSONP(stu, //Content-Type为application/javascript
			"customFunc", // 若不指定，默认为callback
		) //响应体：customFunc({"name":"zcy","addr":"bj"});
	})
}

func xml(app *fiber.App) {
	type Stu struct {
		Name    string `xml:"name"`
		Address string `xml:"addr"`
	}
	var stu Stu //xml不支持匿名结构体
	stu.Name = "zcy"
	stu.Address = "bj"
	app.Get("/user/xml", func(ctx fiber.Ctx) error {
		return ctx.XML(stu) //不能使用匿名结构体。响应头 Content-Type=application/xml
	})
}

func Html(app *fiber.App) {
	app.Get("/user/html", func(ctx fiber.Ctx) error {
		return ctx.Render("template", fiber.Map{
			"title": "用户信息",
			"name":  "zcy",
			"addr":  "bj",
		})
	})
}

func redirect(app *fiber.App) {
	app.Get("/user/old_page", func(ctx fiber.Ctx) error {
		return ctx.Redirect().Status(http.StatusMovedPermanently).To("/user/html")
	})
}

func main5() {
	//Fiber支持多种模板引擎，如：html（go标准库的html template）、django、slim等
	engine := html.New("../static", ".html") //将来ctx.Render()时使用的template name就是web/static/下的文件名去掉.html后缀
	// app := django.New("./views", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Use("/css", static.New("../static/css", static.Config{ //在url是访问目录/css相当于访问文件系统中的web/static/css目录
		Browse:   true,  //允许浏览目录
		Download: false, //请求文件时直接下载
	}))

	text(app)
	json0(app)
	json1(app)
	json2(app)
	jsonp(app)
	xml(app)
	Html(app)
	redirect(app)

	if err := app.Listen("127.0.0.1:5678"); err != nil {
		slog.Error("fiber app start failed", "error", err)
	}
}
