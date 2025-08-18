package main

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcfg"
)

// 打日志
func Logging() {
	// 指定logger配置文件在哪儿
	g.Cfg().GetAdapter().(*gcfg.AdapterFile).SetFileName("conf/learn.yaml")

	ctx := context.Background()
	g.Log().Debug(ctx, "this is debug log") // 指定的Level为INFO，Debug日志不会输出
	g.Log().Info(ctx, "this is info log")
	g.Log().Error(ctx, "this is error log") //error日志还会打印调用堆栈

	type Company struct {
		Name    string `json:"name"`
		Address string `json:"addr"`
	}
	g.Log().Info(ctx, "json is coming", Company{"Microsoft", "Redmond"}) //当传struct或map时，默认打印出来的就是json
}
