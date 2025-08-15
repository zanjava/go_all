package main

import (
	"context"
	"encoding/json"
	"go/post/database"
	mysql "go/post/database/gorm"
	"go/post/handler/model"
	"go/post/util"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
)

const (
	SESSION_KEY_PREFIX = "session_"
	SESSION_LIFE       = 86400
	COOKIE_NAME        = "sesion_id"
)

type userInfo struct {
	Name string
	Id   int
}

func Login(ctx *gin.Context) {
	var user model.User
	err := ctx.ShouldBind(&user)
	if err != nil {
		ctx.String(http.StatusBadRequest, util.BindErrMsg(err))
		return
	}
	user2 := mysql.GetUserByName(user.Name)
	if user2 == nil {
		ctx.String(http.StatusBadRequest, "用户名不存在")
		return
	}
	if user2.PassWord != user.PassWord {
		ctx.String(http.StatusBadRequest, "密码错误")
		return
	}

	//登录成功，生成session_id，以cookie的形式返回给client
	sessionID := xid.New().String() //生成一个随机的字符串
	ctx.SetCookie(
		COOKIE_NAME,
		sessionID,    //注意：受cookie本身的限制，这里的sessionID不能超过4K
		SESSION_LIFE, //maxAge，cookie的有效时间，时间单位秒。如果不设置过期时间，默认情况下关闭浏览器后cookie被删除
		"/",          //path，cookie存放目录
		"localhost",  //cookie从属的域名,不区分协议和端口。如果不指定domain则默认为本host(如b.a.com)，如果指定的domain是一级域名(如a.com)，则二级域名(b.a.com)下也可以访问。访问登录页面时必须用http://localhost:5678/login，而不能用http://127.0.0.1:5678/login，否则浏览器不会保存这个cookie
		false,        //是否只能通过https访问
		true,         //设为false,允许js修改这个cookie（把它设为过期）,js就可以实现logout。如果为true，则需要由后端来重置过期时间
	)
	log.Printf("种下cookie %s", COOKIE_NAME)
	//把用户身份信息存入Redis
	info, _ := json.Marshal(userInfo{Name: user2.Name, Id: user2.Id})
	//key是sessionID, value是sso的用户身份信息
	database.GetRedisClient().Set(context.Background(), SESSION_KEY_PREFIX+sessionID, string(info), SESSION_LIFE*time.Second)
}

func AuthMiddleWare(ctx *gin.Context) {
	if cookie, err := ctx.Request.Cookie(COOKIE_NAME); err == nil {
		sessionID := cookie.Value
		result := database.GetRedisClient().Get(context.Background(), SESSION_KEY_PREFIX+sessionID)
		if result.Err() == nil {
			info := result.Val()
			var user userInfo
			if err := json.Unmarshal([]byte(info), &user); err == nil {
				//把用户信息放入ctx
				ctx.Set("user_name", user.Name)
				ctx.Set("user_id", strconv.Itoa(user.Id))
				return
			} else {
				log.Printf("user info反序列化失败: %s", err)
			}
		} else {
			log.Printf("找不到session id %s: %s", sessionID, err)
		}
	} else {
		log.Printf("读不到cookie %s", COOKIE_NAME)
	}
	ctx.Redirect(http.StatusFound, "/login") //认证失败，重定向到登录页面
}

func main3() {
	mysql.ConnectPostDB("./post/conf", "db", util.YAML, "./log")
	router := gin.Default()

	router.Static("/js", "post/views/js") //在url是访问目录/js相当于访问文件系统中的views/js目录
	router.Static("/css", "post/views/css")
	router.StaticFile("/favicon.ico", "post/views/img/dqq.png") //在url中访问文件/favicon.ico，相当于访问文件系统中的views/img/dqq.png文件
	router.LoadHTMLGlob("post/views/html/*")                    //使用这些.html文件时就不需要加路径了

	// 登录页面
	router.GET("/login", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "session_login.html", nil)
	})
	// 登录
	router.POST("/login/submit", Login)
	//需要登录才能访问的页面1
	router.GET("/page1", AuthMiddleWare, func(ctx *gin.Context) { // 使用身份认证中间件
		// 从ctx里取出用户信息
		uname, _ := ctx.Value("user_name").(string)
		uid, _ := ctx.Value("user_id").(string)
		ctx.String(200, "这是page1, 欢迎 "+uname+"["+uid+"]")
	})
	//需要登录才能访问的页面2
	router.GET("/page2", AuthMiddleWare, func(ctx *gin.Context) { // 使用身份认证中间件
		// 从ctx里取出用户信息
		uname, _ := ctx.Value("user_name").(string)
		uid, _ := ctx.Value("user_id").(string)
		ctx.String(200, "这是page2, 欢迎 "+uname+"["+uid+"]")
	})

	if err := router.Run("127.0.0.1:5678"); err != nil {
		panic(err)
	}
}

// go run ./post/login/session
// 在浏览器里访问 http://localhost:5678/page1  注意一定要走域名localhost
