package main

import (
	"fmt"
	_ "net/http/pprof" //在线pprof

	_ "github.com/go-sql-driver/mysql" //注册mysql驱动
)

func InitLogger() {
	fmt.Println("init logger")
	fmt.Println("main是否匹配正则表达式", Reg.Match([]byte("hello123")))

}

func main() {
	CheckReg()
	InitLogger()
	InitDatabase()

	fmt.Println("server start")
}

func InitDatabase() {
	fmt.Println("init database")
}

// go build -o ab.exe .\init\  生成可执行文件ab.exe
