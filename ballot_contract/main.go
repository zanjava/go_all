package main

import (
	"log"
	"net/http"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gin-gonic/gin"
)

var (
	bcClient *ethclient.Client
)

func Init() {
	var err error
	bcClient, err = ethclient.Dial("http://127.0.0.1:7545") //本地Ganache测试链
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	Init()

	router := gin.Default()
	router.Static("/js", "views/js")    //在url是访问目录/js相当于访问文件系统中的views/js目录
	router.Static("/img", "views/img")  //在url是访问目录/js相当于访问文件系统中的views/js目录
	router.LoadHTMLGlob("views/*.html") //使用这些.html文件时就不需要加路径了

	//不走go后端，直接用web3.js
	router.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "ballot.html", nil)
	})

	router.Run("127.0.0.1:5678")
}
