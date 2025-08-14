package main

import (
	"log/slog"
	"sync"

	"github.com/gin-gonic/gin"
)

var (
	counter sync.Map
)

func init() {
	counter.Store("ye38ry4928", 0)
	counter.Store("689034132t", 0)
}

// 记录每个token调用接口的次数
func CountMD(ctx *gin.Context) {
	token, err := ctx.Cookie("token") //没在cookie里传token
	if err != nil {
		ctx.Abort()
	}
	if v, exists := counter.Load(token); !exists {
		ctx.Abort() //传的token不合法
	} else {
		c, _ := v.(int)
		counter.Store(token, c+1)
		slog.Info("visit counter", "token", token, "count", c+1)
	}
}

func SetCookie(ctx *gin.Context) {
	name := "language"
	value := "go"
	maxAge := 86400 * 7       //cookie的有效时间，时间单位秒。如果不设置过期时间，默认情况下关闭浏览器后cookie被删除
	path := "/"               //cookie存放目录
	domain := "www.baidu.com" //cookie从属的域名,不区分协议和端口。如果不指定domain则默认为本host(如b.a.com)，如果指定的domain是一级域名(如a.com)，则二级域名(b.a.com)下也可以访问
	secure := false           //是否只能通过https访问
	httpOnly := true          //是否允许别人通过js获取(或修改)该cookie，设为false防止XSS攻击
	//SetCookie只能执行一次,第二次SetCookie无效
	ctx.SetCookie(name, value, maxAge, path, domain, secure, httpOnly) //对应的response header key是"Set-Cookie"
}

func main7() {
	engine := gin.Default()
	engine.Use(CountMD)
	engine.GET("/ck", SetCookie)
	engine.Run("127.0.0.1:5678")
}
