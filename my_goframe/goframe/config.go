package main

import (
	"fmt"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcfg"
	"github.com/gogf/gf/v2/os/gctx"
)

// 解析、读取配置文件
func ReadConfigFile() {
	// 指定配置文件。或者在运行程序时加个flag: --gf.gcfg.file=./conf/learn.yaml
	g.Cfg().GetAdapter().(*gcfg.AdapterFile).SetFileName("./conf/learn.yaml")

	ctx := gctx.New()
	secret, err := g.Cfg().Get(ctx, "key")
	if err != nil {
		panic(err)
	}
	fmt.Printf("secret=%s\n", secret)

	masterDSN, _ := g.Cfg().Get(ctx, "database.default.0.link")
	fmt.Printf("master data source name=%s\n", masterDSN)

	// GetWithEnv会优先读配置文件，获取不到时再去读环境变量。
	// 把环境变量中的大写全部转为小写，_转为.
	JavaHome, _ := g.Cfg().GetWithEnv(ctx, "java.home")
	fmt.Printf("JAVA_HOME=%s\n", JavaHome)
}
