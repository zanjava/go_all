package main

import (
	handler "go/post/handler/gin"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main2() {
	Init()

	go ListenTermSignal() //监听终止信号

	// gin.SetMode(gin.ReleaseMode) //GIN线上发布模式
	// gin.DefaultWriter = io.Discard //禁用GIN日志
	engine := gin.Default()

	// 修改静态资源不需要重启GIN，刷新页面即可
	engine.Static("/js", "./views/js") //在url是访问目录/js相当于访问文件系统中的views/js目录
	engine.Static("/css", "./views/css")
	engine.StaticFile("/favicon.ico", "./views/img/dqq.png") //在url中访问文件/favicon.ico，相当于访问文件系统中的views/img/dqq.png文件
	engine.LoadHTMLGlob("./views/html/*")                    //使用这些.html文件时就不需要加路径了

	// 数据上报给Prometheus , grafana做监控和报警等报表展示
	// 系统监控
	engine.Use(handler.Metric)                      // 全局中间件，上报每一个接口的耗时和调用次数
	engine.GET("/metrics", func(ctx *gin.Context) { //Promethus要来访问这个接口
		promhttp.Handler().ServeHTTP(ctx.Writer, ctx.Request)
	})

	engine.GET("/login", func(ctx *gin.Context) { ctx.HTML(http.StatusOK, "login.html", nil) })
	engine.POST("/login", func(ctx *gin.Context) { ctx.HTML(http.StatusOK, "login.html", nil) }) // Auth失败重定向时可能会以post方式请求/login页面
	engine.GET("/regist", func(ctx *gin.Context) { ctx.HTML(http.StatusOK, "user_regist.html", nil) })
	engine.GET("/modify_pass", handler.Auth, func(ctx *gin.Context) { ctx.HTML(http.StatusOK, "update_pass.html", nil) })
	engine.POST("/login/submit", handler.Login)
	engine.POST("/regist/submit", handler.RegistUser)
	engine.POST("/modify_pass/submit", handler.Auth, handler.UpdatePassword)
	engine.GET("/user", handler.GetUserInfo)
	engine.GET("/logout", handler.Logout)

	group := engine.Group("/news")
	group.GET("", handler.NewsList)
	group.GET("/issue", func(ctx *gin.Context) { ctx.HTML(http.StatusOK, "news_issue.html", nil) })
	group.POST("/issue/submit", handler.Auth, handler.PostNews)
	group.GET("/belong", handler.NewsBelong)
	group.GET("/:id", handler.GetNewsById)
	group.GET("/delete/:id", handler.Auth, handler.DeleteNews)
	group.POST("/update", handler.Auth, handler.UpdateNews)

	engine.GET("", func(ctx *gin.Context) { ctx.Redirect(http.StatusMovedPermanently, "news") }) //新闻列表页是默认的首页

	if err := engine.Run("127.0.0.1:5678"); err != nil {
		panic(err)
	}

}
