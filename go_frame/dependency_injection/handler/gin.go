package handler

import (
	"go/frame/dependency_injection/service"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GinHandler struct {
	*gin.Engine
	recommender service.INewsRecommender // 把依赖注入到结构体里去
}

// 这种非自己主动初始化依赖，而通过外部来传入依赖的方式，称为依赖注入
func NewGinHandler(recommender service.INewsRecommender) (*GinHandler, func()) {
	h := &GinHandler{
		gin.Default(),
		recommender,
	}
	return h, h.Close
}

func (h *GinHandler) Rec(ctx *gin.Context) {
	news := h.recommender.RecNews()
	ctx.HTML(http.StatusOK, "news_list.html", gin.H{"data": news})
}

func (h *GinHandler) Route() {
	h.Static("/js", "dependency_injection/views/js")                       //在url是访问目录/js相当于访问文件系统中的views/js目录
	h.StaticFile("/favicon.ico", "dependency_injection/views/img/dqq.png") //在url中访问文件/favicon.ico，相当于访问文件系统中的views/img/dqq.png文件
	h.LoadHTMLGlob("dependency_injection/views/html/*")                    //使用这些.html文件时就不需要加路径了
	h.GET("/news", h.Rec)
}

func (h *GinHandler) Close() {
	log.Println("close gin")
}
