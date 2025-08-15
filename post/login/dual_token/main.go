package main

import (
	"context"
	"go/post/database"
	mysql "go/post/database/gorm"
	"go/post/handler/model"
	"go/post/util"
	"log"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
)

const (
	REFRESH_KEY_PREFIX  = "session_"
	REFRESH_TOKEN_LIFE  = 60 // refresh token的有效期（单位：秒），refresh token失效后用户就得重新登录
	REFRESH_COOKIE_NAME = "refresh"
	ACCESS_COOKIE_NAME  = "access"

	SECRET = "f4398"
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

	//登录成功，生成双token，均以cookie的形式返回给client
	refreshToken := xid.New().String() //生成一个随机的字符串
	ctx.SetCookie(
		REFRESH_COOKIE_NAME,
		refreshToken,
		REFRESH_TOKEN_LIFE,
		"/",
		"localhost",
		false,
		true,
	)
	header := util.DefautHeader
	payload := util.JwtPayload{
		Issue:       "dual_token",
		IssueAt:     time.Now().Unix(),
		Expiration:  0, //永不过期
		UserDefined: map[string]any{"user_id": strconv.Itoa(user2.Id), "user_name": user2.Name},
	}
	if accessToken, err := util.GenJWT(header, payload, SECRET); err != nil {
		slog.Error("生成access token失败", "error", err)
	} else {
		ctx.SetCookie(
			ACCESS_COOKIE_NAME,
			accessToken,
			0, //不设置cookie的过期时间，默认关闭浏览器后cookie失效
			"/",
			"localhost",
			false,
			true,
		)
		// 把<refreshToken, accessToken>存入Redis
		database.GetRedisClient().Set(context.Background(), REFRESH_KEY_PREFIX+refreshToken, accessToken, REFRESH_TOKEN_LIFE*time.Second)
	}
	return
}

func AuthMiddleWare(ctx *gin.Context) {
	if cookie, err := ctx.Request.Cookie(ACCESS_COOKIE_NAME); err == nil {
		accessToken := cookie.Value
		_, payload, err := util.VerifyJwt(accessToken, SECRET)
		if err == nil {
			log.Println("直接根据access token拿到了用户的身份信息")
			//把用户信息放入ctx
			ctx.Set("user_name", payload.UserDefined["user_name"])
			ctx.Set("user_id", payload.UserDefined["user_id"])
			return
		}
	} else {
		log.Println("cookie里没有access token")
	}
	// 如果access token已过期（比如浏览器重新打开了），则使用refresh token
	if cookie, err := ctx.Request.Cookie(REFRESH_COOKIE_NAME); err == nil {
		refreshToken := cookie.Value
		result := database.GetRedisClient().Get(context.Background(), REFRESH_KEY_PREFIX+refreshToken)
		if err == nil {
			log.Println("根据cookie里的refresh token重新获得了access token")
			accessToken := result.Val()
			_, payload, err := util.VerifyJwt(accessToken, SECRET)
			if err == nil {
				//把accessToken以cookie的形式传给client
				ctx.SetCookie(
					ACCESS_COOKIE_NAME,
					accessToken,
					0, //不设置cookie的过期时间，默认关闭浏览器后cookie失效
					"/",
					"localhost",
					false,
					true,
				)
				log.Println("把access token种到了浏览器的cookie里")
				//把用户信息放入ctx
				ctx.Set("user_name", payload.UserDefined["user_name"])
				ctx.Set("user_id", payload.UserDefined["user_id"])
				return
			}
		}
	}
	ctx.Redirect(http.StatusFound, "/login") //认证失败，重定向到登录页面
}

func main5() {
	util.InitSlog("./log/post.log")
	mysql.ConnectPostDB("./conf", "db", util.YAML, "./log")
	router := gin.Default()

	router.Static("/js", "./views/js") //在url是访问目录/js相当于访问文件系统中的views/js目录
	router.Static("/css", "./views/css")
	router.StaticFile("/favicon.ico", "./views/img/dqq.png") //在url中访问文件/favicon.ico，相当于访问文件系统中的views/img/dqq.png文件
	router.LoadHTMLGlob("./views/html/*")                    //使用这些.html文件时就不需要加路径了

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

// go run ./post/login/dual_token
// 在浏览器里访问 http://localhost:5678/page1  注意一定要走域名localhost
