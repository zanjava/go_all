package main

import (
	_ "my_goframe/internal/packed"

	"github.com/gogf/gf/v2/os/gctx"

	"my_goframe/internal/cmd"
)

func main() {
	cmd.Main.Run(gctx.GetInitCtx())
}
