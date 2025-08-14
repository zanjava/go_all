package main

import (
	"log/slog"

	"github.com/gin-gonic/gin"
)

func M1(ctx *gin.Context) {
	ctx.String(200, "M1 Begin\n")
	ctx.Next() // 先把后面的函数执行完，再进入下一行
	ctx.String(200, "M1 End\n")
}

func M2(ctx *gin.Context) {
	ctx.String(200, "Here is M2\n")
	// 即使没调ctx.Next()，本函数执行完后，也会进入下一个函数
}

func M3(ctx *gin.Context) {
	ctx.String(200, "Here is M3\n")
}

func M4(ctx *gin.Context) {
	ctx.String(200, "Here is M4\n")
	ctx.Abort() //通过Abort()使中间件后面的handler不再执行，但是本handler还是会继续执行。所以下一行代码需要显式return
	ctx.String(200, "M4 Second Line\n")
}

func M5(ctx *gin.Context) {
	ctx.String(200, "Here is M5\n")
}

func M6(ctx *gin.Context) {
	slog.Info("visit", "path", ctx.Request.URL)
}

func main2() {
	engine := gin.Default()
	engine.Use(M6)                          //全局中间件
	engine.GET("/", M1, M2, M3, M2, M4, M5) //中间件可以有多个
	engine.GET("/2", M2, M3, M2, M4, M5)
	engine.Run("127.0.0.1:5678")
}
