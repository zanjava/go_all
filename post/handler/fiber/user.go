package handler

import (
	database "go/post/database/xorm"
	"go/post/handler/model"
	"go/post/util"
	"log/slog"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v3"
)

const (
	COOKIE_LIFE = 7 * 86400
)

// 注册新用户。pass是md5之后的密码
func RegistUser(ctx fiber.Ctx) error {
	var user model.User
	err := ctx.Bind().Form(&user)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).SendString(util.BindErrMsg(err))
	}

	err = database.RegistUser(user.Name, user.PassWord)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).SendString(err.Error())
	}
	return nil
}

// 修改密码
func UpdatePassword(ctx fiber.Ctx) error {
	loginUid, _ := ctx.Locals(UID_IN_CTX).(int)
	if loginUid <= 0 {
		return ctx.Status(http.StatusForbidden).SendString("请先登录")
	}

	var req model.ModifyPassRequest
	err := ctx.Bind().Form(&req)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).SendString(util.BindErrMsg(err))
	}

	err = database.UpdatePassword(loginUid, req.NewPass, req.OldPass)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).SendString(err.Error())
	}
	return nil
}

// 登录
func Login(ctx fiber.Ctx) error {
	var user model.User
	err := ctx.Bind().Form(&user)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).SendString(util.BindErrMsg(err))
	}

	user2 := database.GetUserByName(user.Name)
	if user2 == nil {
		return ctx.Status(http.StatusBadRequest).SendString("用户名不存在")
	}
	if user2.PassWord != user.PassWord {
		return ctx.Status(http.StatusBadRequest).SendString("密码错误")
	}

	//登录成功，返回cookie
	header := util.DefautHeader
	payload := util.JwtPayload{ //payload以明文形式编码在token中，server用自己的密钥可以校验该信息是否被篡改过
		Issue:       "news",
		IssueAt:     time.Now().Unix(),                                //因为每次的IssueAt不同，所以每次生成的token也不同
		Expiration:  time.Now().Add(COOKIE_LIFE * time.Second).Unix(), //7天后过期
		UserDefined: map[string]any{UID_IN_TOKEN: user2.Id},           //用户自定义字段。如果token里包含敏感信息，请结合https使用
	}
	if token, err := util.GenJWT(header, payload, KeyConfig.GetString("secret")); err != nil {
		slog.Error("生成token失败", "error", err)
		return ctx.Status(http.StatusInternalServerError).SendString("token生成失败")
	} else {
		//response header里会有一条 Set-Cookie: jwt=xxx; other_key=other_value，浏览器后续请求会自动把同域名下的cookie再放到request header里来，即request header里会有一条Cookie: jwt=xxx; other_key=other_value
		ctx.Cookie( //设置Cookie，对应的response header key是"Set-Cookie"
			&fiber.Cookie{
				Name:     COOKIE_NAME,
				Value:    token,       //注意：受cookie本身的限制，这里的token不能超过4K
				MaxAge:   COOKIE_LIFE, //maxAge，cookie的有效时间，时间单位秒。如果不设置过期时间，默认情况下关闭浏览器后cookie被删除
				Path:     "/",         //path，cookie存放目录
				Domain:   "localhost", //cookie从属的域名,不区分协议和端口。如果不指定domain则默认为本host(如b.a.com)，如果指定的domain是一级域名(如a.com)，则二级域名(b.a.com)下也可以访问。访问登录页面时必须用http://localhost:5678/login，而不能用http://127.0.0.1:5678/login，否则浏览器不会保存这个cookie
				Secure:   false,       //是否只能通过https访问
				HTTPOnly: true,        //设为false,允许js修改这个cookie（把它设为过期）,js就可以实现logout。如果为true，则需要由后端来重置过期时间
				SameSite: "Strict",    //同源Cookie，防止CSRF攻击。子域名不是同源，所以用户体验不好，且目前仅新版的Chrome和Firefox支持
			},
		)
	}
	return nil
}

// 退出登录
func Logout(ctx fiber.Ctx) error {
	ctx.Cookie( //设置Cookie，对应的response header key是"Set-Cookie"
		&fiber.Cookie{
			Name:     COOKIE_NAME,
			Value:    "",
			MaxAge:   -1, //把Max-Age设为负数即要求浏览器删除该cookie
			Path:     "/",
			Domain:   "localhost",
			Secure:   false,
			HTTPOnly: true,
			SameSite: "Strict",
		},
	)
	return nil
}

func GetUserInfo(ctx fiber.Ctx) error {
	loginUid := GetLoginUid(ctx)
	if loginUid > 0 { //成功从cookie里拿到了登录者的user id
		user := database.GetUserById(loginUid)
		if user != nil {
			return ctx.Status(http.StatusOK).JSON(user) //返回用户信息
		}
	}
	return ctx.Status(http.StatusOK).JSON(model.User{})
}
