package handler

import (
	"go/frame/ioc/service"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GinHandler struct {
	Engine    *gin.Engine
	sort_time service.INewsRecommender `container:"name"`
}

func (h *GinHandler) Rec(ctx *gin.Context) {
	news := h.sort_time.RecNews()
	ctx.HTML(http.StatusOK, "news_list.html", gin.H{"data": news})
}

func (h *GinHandler) Route() {
	h.Engine.Static("/js", "dependency_injection/views/js")                       //在url是访问目录/js相当于访问文件系统中的views/js目录
	h.Engine.StaticFile("/favicon.ico", "dependency_injection/views/img/dqq.png") //在url中访问文件/favicon.ico，相当于访问文件系统中的views/img/dqq.png文件
	h.Engine.LoadHTMLGlob("dependency_injection/views/html/*")                    //使用这些.html文件时就不需要加路径了
	h.Engine.GET("/news", h.Rec)
}

func (h *GinHandler) Close() {
	log.Println("close gin")
}

func (h *GinHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.Engine.ServeHTTP(w, r)
}
