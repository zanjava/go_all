package main

import "github.com/gin-gonic/gin"

func main3() {
	engine := gin.Default()

	{ // 圈定变量作用域
		g1 := engine.Group("/v1") //本组path的公共前缀
		g1.Use(M6)                //本组的公共中间件
		g1.GET("/a", func(ctx *gin.Context) {
			ctx.String(200, "name=gw")
		})
		g1.GET("/b", func(ctx *gin.Context) {
			ctx.String(200, "age=18")
		})
	}

	{
		g2 := engine.Group("/v2")
		g2.GET("/a", func(ctx *gin.Context) {
			ctx.JSON(200, gin.H{"name": "gw"})
		})
		g2.GET("/b", func(ctx *gin.Context) {
			ctx.JSON(200, gin.H{"age": 18})
		})
	}

	engine.Run("127.0.0.1:5678")
}
