package main

import (
	"log"
	"net/http"
	"strconv"

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
	go ListenEvent() //后端go监听区块链事件
}

func Update(ctx *gin.Context) {
	x := ctx.PostForm("x")
	if i, err := strconv.ParseInt(x, 10, 64); err != nil {
		ctx.String(http.StatusBadRequest, "x must be integer")
		return
	} else {
		if err := SetNumber(bcClient, i); err != nil {
			ctx.String(http.StatusInternalServerError, "write x failed")
			return
		} else {
			//默认返回200
		}
	}
}

func Get(ctx *gin.Context) {
	if x, err := GetNumber(bcClient); err != nil {
		ctx.String(http.StatusInternalServerError, "read x failed")
		return
	} else {
		ctx.HTML(http.StatusOK, "index.html", gin.H{"x": x})
		return
	}
}

func main() {
	Init()

	router := gin.Default()
	router.LoadHTMLFiles("index.html", "home.html")

	router.GET("/", Get)
	router.POST("/update", Update)

	//不走go后端，直接用web3.js
	router.GET("/web3", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "home.html", nil)
	})

	router.Run("127.0.0.1:5678")
}
