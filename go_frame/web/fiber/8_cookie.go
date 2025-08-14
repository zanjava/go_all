package main

import (
	"log/slog"
	"net/http"
	"sync"

	"github.com/gofiber/fiber/v3"
)

var (
	counter sync.Map
)

func init() {
	counter.Store("ye38ry4928", 0)
	counter.Store("689034132t", 0)
}

// 记录每个token调用接口的次数
func CountMD(ctx fiber.Ctx) error {
	// 请求头 Cookie: token=32jr34r; city=bj; pass=432432
	token := ctx.Cookies("token") //获取Cookie
	if len(token) == 0 {          //没在cookie里传token
		return fiber.NewError(http.StatusForbidden, "must indicate token")
	}
	if v, exists := counter.Load(token); !exists {
		return fiber.NewError(http.StatusForbidden, "invalid token")
	} else {
		c, _ := v.(int)
		counter.Store(token, c+1)
		slog.Info("visit counter", "token", token, "count", c+1)
		return ctx.Next()
	}
}

func SetCookie(ctx fiber.Ctx) error {
	// ctx.ClearCookie()                    //使客户端所有cookie过期
	//ctx.ClearCookie("language", "token") //指定cookie过期

	//SetCookie只能执行一次,第二次SetCookie无效
	ctx.Cookie( //设置Cookie，对应的response header key是"Set-Cookie"
		&fiber.Cookie{
			Name:     "language",
			Value:    "go",
			MaxAge:   86400 * 7,       //cookie的有效时间，时间单位秒。如果不设置过期时间，默认情况下关闭浏览器后cookie被删除
			Path:     "/",             //cookie存放目录
			Domain:   "www.baidu.com", //cookie从属的域名,不区分协议和端口。如果不指定domain则默认为本host(如b.a.com)，如果指定的domain是一级域名(如a.com)，则二级域名(b.a.com)下也可以访问
			Secure:   false,           //是否只能通过https访问
			HTTPOnly: true,            //是否允许别人通过js获取自己的cookie，设为false防止XSS攻击
			SameSite: "Strict",        //同源Cookie，防止CSRF攻击。子域名不是同源，所以用户体验不好，且目前仅新版的Chrome和Firefox支持
		},
	)
	return ctx.SendString("请保管好Cookie")
}

func main() {
	app := fiber.New()
	app.Use(CountMD)
	app.Get("ck", SetCookie)
	app.Post("ckp", SetCookie)
	if err := app.Listen("127.0.0.1:5678"); err != nil {
		slog.Error("fiber app start failed", "error", err)
	}
}
