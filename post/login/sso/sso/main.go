package main

import (
	database "go/post/database/gorm"
	"go/post/handler/model"
	"go/post/login/sso"
	"go/post/util"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	SECRET = "8j9huumu"
)

func Login(ctx *gin.Context) {
	var user model.User
	err := ctx.ShouldBind(&user)
	if err != nil {
		ctx.String(http.StatusBadRequest, util.BindErrMsg(err))
		return
	}
	user2 := database.GetUserByName(user.Name)
	if user2 == nil {
		ctx.String(http.StatusBadRequest, "用户名不存在")
		return
	}
	if user2.PassWord != user.PassWord {
		ctx.String(http.StatusBadRequest, "密码错误")
		return
	}

	var token string
	//登录成功，生成JWT token，返回cookie
	header := util.DefautHeader
	payload := util.JwtPayload{
		Issue:       "sso",
		IssueAt:     time.Now().Unix(),
		Expiration:  time.Now().Add(sso.SSO_TOKEN_LIFE * time.Second).Unix(),
		UserDefined: map[string]any{sso.KEY_OF_UID: strconv.Itoa(user2.Id), sso.KEY_OF_NAME: user2.Name},
	}
	if token, err = util.GenJWT(header, payload, SECRET); err != nil {
		slog.Error("生成token失败", "error", err)
	}
	ctx.String(200, token)
	return
}

// 验证一个token是否合法。如果不合法，返回code=403；如果合法返回json，包含uid和name
func Identify(ctx *gin.Context) {
	token := ctx.Query(sso.SSO_TOKEN_QUERY_NAME)
	_, payload, err := util.VerifyJwt(token, SECRET)
	if err != nil {
		ctx.String(http.StatusForbidden, "token验证失败")
		return
	} else {
		var uid string
		var name string
		if v, exists := payload.UserDefined[sso.KEY_OF_UID]; exists {
			uid = v.(string)
		}
		if v, exists := payload.UserDefined[sso.KEY_OF_NAME]; exists {
			name = v.(string)
		}
		if uid == "" || name == "" {
			ctx.String(http.StatusForbidden, "token验证失败")
			return
		} else {
			ctx.JSON(http.StatusOK, gin.H{sso.KEY_OF_UID: uid, sso.KEY_OF_NAME: name})
		}
	}
}

func main9() {
	database.ConnectPostDB("./conf", "db", util.YAML, "./log")

	engine := gin.Default()
	engine.Static("/js", "./views/js") //在url是访问目录/js相当于访问文件系统中的views/js目录
	engine.Static("/css", "./views/css")
	engine.StaticFile("/favicon.ico", "./views/img/dqq.png") //在url中访问文件/favicon.ico，相当于访问文件系统中的views/img/dqq.png文件
	engine.LoadHTMLGlob("./views/html/*")                    //使用这些.html文件时就不需要加路径了

	// 登录页面
	engine.GET("/login", func(ctx *gin.Context) {
		service := ctx.Query("service")
		ctx.HTML(http.StatusOK, "oss_login.html", gin.H{"service": service})
	})
	// 登录
	engine.POST("/login/submit", Login)
	// 验证一个token是否合法
	engine.GET("/identify", Identify)
	if err := engine.Run(sso.SSO_URL); err != nil {
		panic(err)
	}
}

// go run ./post/login/sso/sso
