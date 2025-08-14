package main

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func homeHandler(ctx *gin.Context) {
	fmt.Println("请求头:")
	for k, v := range ctx.Request.Header {
		fmt.Printf("%s=%s\n", k, v[0])
	}
	fmt.Println("请求体:")
	io.Copy(os.Stdout, ctx.Request.Body)

	// 下面三项顺序没要求
	// 设置响应体
	ctx.Writer.WriteString("welcome")
	// 设置响应码
	ctx.Writer.WriteHeader(http.StatusOK)
	// 设置响应头
	ctx.Writer.Header().Add("language", "go")
	ctx.Header("Set-Cookie", "language=go; Path=/; Domain=www.baidu.com; Max-Age=604800; HttpOnly")

	// gin框架当然也封装了很多功能
	ctx.String(http.StatusOK, " to BeiJing ")
	ctx.JSON(http.StatusOK, map[string]any{"物理": 34, "化学": 70})
	ctx.JSON(http.StatusOK, gin.H{"物理": 34, "化学": 70}) //type H map[string]any
}

func main1() {
	engine := gin.Default() // Default 使用 Logger 和 Recovery 中间件
	// engine := gin.New() //不使用默认的中间件
	// engine.Use(gin.Logger())
	// engine.Use(gin.Recovery())
	engine.GET("/home", homeHandler)
	if err := engine.Run("127.0.0.1:5678"); err != nil {
		slog.Error("gin server start failed", "error", err)
	}
}

// go run ./web/gin
