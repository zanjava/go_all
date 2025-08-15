package handler

import (
	mysql "go/post/database/gorm"
	"go/post/handler/model"
	"go/post/util"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// 注册新用户。pass是md5之后的密码
func RegistUser(ctx *gin.Context) {
	var user model.User
	err := ctx.ShouldBind(&user)
	if err != nil {
		ctx.String(http.StatusBadRequest, util.BindErrMsg(err))
		return
	}

	err = mysql.RegistUser(user.Name, user.PassWord)
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	// ctx.Status(http.StatusOK) //这行代码可以不写，默认情况下就是返回200
}

const (
	SESSION_KEY_PREFIX = "session_"
	SESSION_LIFE       = 86400
	COOKIE_NAME1       = "session_id"
)

type userInfo struct {
	Name string
	Id   int
}

// 登录
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

	slog.Info("登录成功", "uid", user2.Id)
	// 2.jwt认证
	// 登录成功，返回cookie
	header := util.DefautHeader
	payload := util.JwtPayload{ //payload以明文形式编码在token中，server用自己的密钥可以校验该信息是否被篡改过
		Issue:       "news",
		IssueAt:     time.Now().Unix(),                                //因为每次的IssueAt不同，所以每次生成的token也不同
		Expiration:  time.Now().Add(COOKIE_LIFE * time.Second).Unix(), //7天后过期
		UserDefined: map[string]any{UID_IN_TOKEN: user2.Id},           //用户自定义字段。如果token里包含敏感信息，请结合https使用
	}
	if token, err := util.GenJWT(header, payload, KeyConfig.GetString("secret")); err != nil {
		slog.Error("生成token失败", "error", err)
		ctx.String(http.StatusInternalServerError, "token生成失败")
	} else {
		//response header里会有一条 Set-Cookie: jwt=xxx; other_key=other_value，浏览器后续请求会自动把同域名下的cookie再放到request header里来，即request header里会有一条Cookie: jwt=xxx; other_key=other_value
		ctx.SetCookie(
			COOKIE_NAME,
			token,       //注意：受cookie本身的限制，这里的token不能超过4K
			COOKIE_LIFE, //maxAge，cookie的有效时间，时间单位秒。如果不设置过期时间，默认情况下关闭浏览器后cookie被删除
			"/",         //path，cookie存放目录
			"localhost", //cookie从属的域名,不区分协议和端口。如果不指定domain则默认为本host(如b.a.com)，如果指定的domain是一级域名(如a.com)，则二级域名(b.a.com)下也可以访问。访问登录页面时必须用http://localhost:5678/login，而不能用http://127.0.0.1:5678/login，否则浏览器不会保存这个cookie
			false,       //是否只能通过https访问
			true,        //设为false,允许js修改这个cookie（把它设为过期）,js就可以实现logout。如果为true，则需要由后端来重置过期时间
		)
	}

	// 1.cookie认证
	// ctx.SetCookie(
	// 	"uid",
	// 	strconv.Itoa(user2.Id),
	// 	3600,
	// 	"/",
	// 	"localhost",
	// 	false,
	// 	true,
	// )

	//3.登录成功，生成session_id，以cookie的形式返回给client
	//session 认证
	// if GetSessionId(ctx) == "" {
	// 	sessionID := xid.New().String() //生成一个随机的字符串
	// 	ctx.SetCookie(
	// 		COOKIE_NAME1,
	// 		sessionID,    //注意：受cookie本身的限制，这里的sessionID不能超过4K
	// 		SESSION_LIFE, //maxAge，cookie的有效时间，时间单位秒。如果不设置过期时间，默认情况下关闭浏览器后cookie被删除
	// 		"/",          //path，cookie存放目录
	// 		"localhost",  //cookie从属的域名,不区分协议和端口。如果不指定domain则默认为本host(如b.a.com)，如果指定的domain是一级域名(如a.com)，则二级域名(b.a.com)下也可以访问。访问登录页面时必须用http://localhost:5678/login，而不能用http://127.0.0.1:5678/login，否则浏览器不会保存这个cookie
	// 		false,        //是否只能通过https访问
	// 		true,         //设为false,允许js修改这个cookie（把它设为过期）,js就可以实现logout。如果为true，则需要由后端来重置过期时间
	// 	)
	// 	log.Printf("种下cookie %s", COOKIE_NAME1)
	// 	//把用户身份信息存入Redis
	// 	info, _ := json.Marshal(userInfo{Name: user2.Name, Id: user2.Id})
	// 	//key是sessionID, value是sso的用户身份信息
	// 	database.GetRedisClient().Set(context.Background(), SESSION_KEY_PREFIX+sessionID, string(info), SESSION_LIFE*time.Second)
	// }

}

func GetSessionId(ctx *gin.Context) string {
	sessionID := ""
	for _, cookie := range ctx.Request.Cookies() {
		if cookie.Name == COOKIE_NAME1 {
			sessionID = cookie.Value
		}
	}
	return sessionID
}

// 退出登录
func Logout(ctx *gin.Context) {
	// if sessionID := GetSessionId(ctx); sessionID != "" {
	// 	database.GetRedisClient().Del(context.Background(), SESSION_KEY_PREFIX+sessionID)
	// }
	// ctx.SetCookie(COOKIE_NAME1, "", -1, "/", "localhost", false, true) //把Max-Age设为负数即要求浏览器删除该cookie
	ctx.SetCookie(COOKIE_NAME, "", -1, "/", "localhost", false, true) //把Max-Age设为负数即要求浏览器删除该cookie
}

// 修改密码
func UpdatePassword(ctx *gin.Context) {
	uid, ok := ctx.Value(UID_IN_CTX).(int)
	if !ok {
		ctx.String(http.StatusForbidden, "请先登录")
		return
	}

	var req model.ModifyPassRequest
	err := ctx.ShouldBind(&req)
	if err != nil {
		ctx.String(http.StatusBadRequest, util.BindErrMsg(err))
		return
	}

	//uid := GetUidFromCookie(ctx)
	// if uid <= 0 {
	// 	ctx.String(http.StatusForbidden, "请先登录")
	// 	return
	// }

	err = mysql.UpdatePassword(uid, req.NewPass, req.OldPass)
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
}

func GetUserInfo(ctx *gin.Context) {
	loginUid := GetLoginUid(ctx)
	//loginUid := GetLoginUidFromSession(ctx)
	if loginUid > 0 { //成功从cookie里拿到了登录者的user id
		user := mysql.GetUserById(loginUid)
		if user != nil {
			slog.Info("GetUserInfo", "uid", user.Id, "name", user.Name)
			ctx.JSON(http.StatusOK, user) //返回用户信息
			return
		}
	}
	ctx.JSON(http.StatusOK, model.User{})
}

// func GetUidFromCookie(ctx *gin.Context) int {
// 	uidStr, err := ctx.Cookie("uid")
// 	if err != nil {
// 		return 0
// 	}
// 	uid, err := strconv.Atoi(uidStr)
// 	if err != nil {
// 		return 0
// 	}
// 	return uid
// }
