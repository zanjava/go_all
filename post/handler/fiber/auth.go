package handler

import (
	"go/post/util"
	"net/http"

	"github.com/gofiber/fiber/v3"
)

const (
	UID_IN_TOKEN = "uid"
	UID_IN_CTX   = "uid"
	COOKIE_NAME  = "jwt"
)

var (
	KeyConfig = util.InitViper("./conf", "jwt", util.YAML)
)

// 从cookie里取出jwt，从而取出uid
func GetLoginUid(ctx fiber.Ctx) int {
	//依靠浏览器自动回传的cookie，提取出jwt token
	token := ctx.Cookies(COOKIE_NAME)
	// slog.Info("take cookie", "name", COOKIE_NAME, "value", token)
	return GetUidFromJwt(token)
}

// 从jwt里取出uid
func GetUidFromJwt(token string) int {
	_, payload, err := util.VerifyJwt(token, KeyConfig.GetString("secret"))
	if err != nil {
		return 0
	}
	for k, v := range payload.UserDefined {
		if k == UID_IN_TOKEN {
			return int(v.(float64))
		}
	}
	return 0
}

// 身份认证中间件，先确保是登录状态
func Auth(ctx fiber.Ctx) error {
	loginUid := GetLoginUid(ctx)
	// slog.Info("auth", "uid", loginUid)
	if loginUid <= 0 {
		return ctx.Redirect().Status(http.StatusTemporaryRedirect).To("/login") //重定向到登录页面
	} else {
		ctx.Locals(UID_IN_CTX, loginUid) //把登录的uid放入ctx中
		return ctx.Next()
	}
}
