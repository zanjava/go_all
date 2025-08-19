package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/golobby/container/v3"
)

func init() {
	// 告诉IoC容器如何创建接口的具体实现
	err := container.Singleton(func() IHandler {
		var h GinHandler
		err := container.Fill(&h) // 通过IoC容器创建一个bean（只IoC容器只负责给带container Tag的Field赋值）
		if err != nil {
			panic(err)
		}
		h.Engine = gin.Default() //自行给其他Field赋值
		return &h
	})
	if err != nil {
		panic(err)
	}
}
