package web_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/middlewares/server/recovery"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/protocol"
	"github.com/gin-gonic/gin"
	gojson "github.com/goccy/go-json"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/recover"
)

var (
	name     = "language"
	value    = "go"
	maxAge   = 86400 * 7
	path     = "/"
	domain   = "www.baidu.com"
	secure   = false
	httpOnly = true
)

func InitGin() *gin.Engine {
	engine := gin.New()
	engine.Use(gin.Recovery())
	return engine
}

func InitHertz(port int) *server.Hertz {
	engine := server.New(server.WithHostPorts("127.0.0.1:" + strconv.Itoa(port)))
	engine.Use(recovery.Recovery())
	return engine
}

func InitFiber() *fiber.App {
	app := fiber.New(fiber.Config{
		JSONEncoder: gojson.Marshal,
		JSONDecoder: gojson.Unmarshal,
	})
	app.Use(recover.New())
	return app
}

type Arg struct {
	Name string `form:"name" binding:"required,gt=4" validate:"required,gt=4" vd:"len($)>4"`
	Age  int    `form:"age" binding:"required,gt=0" validate:"required,gt=0" vd:"$>0"`
}

func GinHandler(ctx *gin.Context) {
	var arg Arg
	ctx.ShouldBind(&arg)         //绑定参数并校验
	ctx.Header("company", "dqq") //响应头
	ctx.SetCookie(name, value, maxAge, path, domain,
		secure, httpOnly) //返回Cookie
	ctx.JSON(200, arg) //响应体
}

func HertzHandler(c context.Context, ctx *app.RequestContext) {
	var arg Arg
	ctx.BindAndValidate(&arg)    //绑定参数并校验
	ctx.Header("company", "dqq") //响应头
	ctx.SetCookie(name, value, maxAge, path, domain,
		protocol.CookieSameSiteDefaultMode, secure, httpOnly) //返回Cookie
	ctx.JSON(200, arg) //响应体
}

func FiberHander(ctx fiber.Ctx) error {
	var arg Arg
	ctx.Bind().Form(&arg)     //绑定参数并校验
	ctx.Set("company", "dqq") //响应头

	ctx.Cookie(&fiber.Cookie{Name: name, Value: value, Path: path,
		MaxAge: maxAge, Domain: domain, Secure: secure, HTTPOnly: httpOnly}) //返回Cookie

	ctx.JSON(arg) //响应体
	return nil
}

func GinServe(port int) {
	engine := InitGin()
	engine.POST("", GinHandler)
	engine.Run("127.0.0.1:" + strconv.Itoa(port))
}

func HertzServe(port int) {
	engine := InitHertz(port)
	engine.POST("", HertzHandler)
	engine.Spin()
}

func FiberServe(port int) {
	app := InitFiber()
	app.Post("", FiberHander)
	app.Listen("127.0.0.1:" + strconv.Itoa(port))
}

func Post(uri string) {
	resp, _ := http.PostForm(uri, url.Values{"name": {"张三"}, "age": {"18"}})
	for k, v := range resp.Header {
		fmt.Printf("%s=%s\n", k, v[0])
	}
	io.Copy(os.Stdout, resp.Body)
	resp.Body.Close()
}

func TestGIN(t *testing.T) {
	port := 5678
	go GinServe(port)
	time.Sleep(3 * time.Second)
	Post("http://127.0.0.1:" + strconv.Itoa(port))
}

func TestHertz(t *testing.T) {
	port := 5680
	go HertzServe(port)
	time.Sleep(3 * time.Second)
	Post("http://127.0.0.1:" + strconv.Itoa(port))
}

func TestFiber(t *testing.T) {
	port := 5679
	go FiberServe(port)
	time.Sleep(3 * time.Second)
	Post("http://127.0.0.1:" + strconv.Itoa(port))
}

func BenchmarkGIN(b *testing.B) {
	port := 5678
	go GinServe(port) //启动Server
	time.Sleep(3 * time.Second)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		http.PostForm("http://127.0.0.1:"+strconv.Itoa(port), url.Values{"name": {"张三"}, "age": {"18"}})
	}
}

func BenchmarkHertz(b *testing.B) {
	port := 5680
	go HertzServe(port) //启动Server
	time.Sleep(3 * time.Second)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		http.PostForm("http://127.0.0.1:"+strconv.Itoa(port), url.Values{"name": {"张三"}, "age": {"18"}})
	}
}

func BenchmarkFiber(b *testing.B) {
	port := 5679
	go FiberServe(port) //启动Server
	time.Sleep(3 * time.Second)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		http.PostForm("http://127.0.0.1:"+strconv.Itoa(port), url.Values{"name": {"张三"}, "age": {"18"}})
	}
}

// go test -v ./web -run=^TestGIN$ -count=1 --ldflags="-checklinkname=0"
// go test -v ./web -run=^TestFiber$ -count=1 --ldflags="-checklinkname=0"
// go test -v ./web -run=^TestHertz$ -count=1 --ldflags="-checklinkname=0"
// go test ./web -bench=^BenchmarkGIN$ -run=^$ -count=1 -ldflags="-checklinkname=0"
// go test ./web -bench=^BenchmarkFiber$ -run=^$ -count=1 --ldflags="-checklinkname=0"
// go test ./web -bench=^BenchmarkHertz$ -run=^$ -count=1 --ldflags="-checklinkname=0"

/**
goos: windows
goarch: amd64
pkg: dqq/go/frame/web
cpu: 11th Gen Intel(R) Core(TM) i7-1165G7 @ 2.80GHz

都使用标准库的json序列化
BenchmarkGIN-8              3194            373073 ns/op
BenchmarkFiber-8            2698            424316 ns/op

Fiber使用go-json序列化
BenchmarkGIN-8              3012            366533 ns/op
BenchmarkFiber-8            3132            403631 ns/op

BenchmarkFiber-8            3058            379229 ns/op
BenchmarkGIN-8              3115            385951 ns/op
BenchmarkHertz-8            3086            389009 ns/op
*/
