package handler

import (
	database "go/post/database/xorm"
	"go/post/handler/model"
	"go/post/util"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/fatih/structs"
	"github.com/gofiber/fiber/v3"
)

// 新闻发布
func PostNews(ctx fiber.Ctx) error {
	loginUid, _ := ctx.Locals(UID_IN_CTX).(int)
	if loginUid <= 0 {
		return ctx.Status(http.StatusForbidden).SendString("请先登录")
	}

	var news model.News
	err := ctx.Bind().Form(&news)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).SendString(util.BindErrMsg(err))
	}

	id, err := database.PostNews(loginUid, news.Title, news.Content)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).SendString(err.Error())
	}
	return ctx.Status(http.StatusOK).JSON(fiber.Map{"id": id})
}

// 新闻详情页
func GetNewsById(ctx fiber.Ctx) error {
	idStr := ctx.Params("id")
	if id, err := strconv.Atoi(idStr); err != nil || id <= 0 {
		return ctx.Status(http.StatusBadRequest).SendString("非法的新闻id")
	} else {
		news := database.GetNewsById(id)
		if news == nil {
			return ctx.SendStatus(http.StatusNotFound)
		}

		user := database.GetUserById(news.UserId)
		if user != nil {
			news.UserName = user.Name
		}
		return ctx.Status(http.StatusOK).Render("news", structs.Map(news)) // 渲染时不支持传结构体，所以只好借助于第三方库把结构体转为map
	}
}

// 删除新闻
func DeleteNews(ctx fiber.Ctx) error {
	loginUid, _ := ctx.Locals(UID_IN_CTX).(int)
	if loginUid <= 0 {
		return ctx.Status(http.StatusForbidden).SendString("请先登录")
	}

	idStr := ctx.Params("id")
	if id, err := strconv.Atoi(idStr); err != nil || id <= 0 {
		return ctx.Status(http.StatusBadRequest).SendString("非法的新闻id")
	} else {
		if !newsBelongUser(id, loginUid) {
			return ctx.Status(http.StatusForbidden).SendString("无权限删除")
		}
		err = database.DeleteNews(id)
		if err != nil {
			return ctx.Status(http.StatusInternalServerError).SendString(err.Error())
		}
		return ctx.SendStatus(http.StatusOK)
	}
}

// 修改新闻
func UpdateNews(ctx fiber.Ctx) error {
	loginUid, _ := ctx.Locals(UID_IN_CTX).(int)
	if loginUid <= 0 {
		return ctx.Status(http.StatusForbidden).SendString("请先登录")
	}

	var news model.News
	err := ctx.Bind().Form(&news)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).SendString(util.BindErrMsg(err))
	}
	if news.Id <= 0 {
		return ctx.Status(http.StatusBadRequest).SendString("非法的新闻id")
	}

	if !newsBelongUser(news.Id, loginUid) {
		return ctx.Status(http.StatusForbidden).SendString("无权限删除")
	}

	err = database.UpdateNews(news.Id, news.Title, news.Content)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).SendString(err.Error())
	}
	return ctx.SendStatus(http.StatusOK)
}

// 新闻列表页
func NewsList(ctx fiber.Ctx) error {
	PageNo, err1 := strconv.Atoi(ctx.Query("page_no", "1"))
	PageSize, err2 := strconv.Atoi(ctx.Query("page_size", "3"))
	if err1 != nil || err2 != nil {
		return ctx.Status(http.StatusBadRequest).SendString("page参数不合法")
	}

	total, news := database.GetNewsByPage(PageNo, PageSize)
	if total < len(news) {
		return ctx.Status(http.StatusInternalServerError).SendString("出错了")
	}
	for _, ele := range news {
		user := database.GetUserById(ele.UserId)
		if user != nil {
			ele.UserName = user.Name
		} else {
			slog.Warn("could not get name of user", "uid", ele.UserId)
		}
	}
	return ctx.Status(http.StatusOK).Render("news_list", fiber.Map{"total": total, "data": news, "page": PageNo})
}

// 从jwt里解析出uid，判断news id是否属于uid
func NewsBelong(ctx fiber.Ctx) error {
	newsId := ctx.Query("id") //新闻id
	nid, err := strconv.Atoi(newsId)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).SendString("非法的新闻id")
	}

	loginUid := GetLoginUid(ctx)
	if newsBelongUser(nid, loginUid) {
		return ctx.Status(http.StatusOK).SendString("true") //新闻的作者id就是当前登录者的uid
	} else {
		return ctx.Status(http.StatusOK).SendString("false")
	}
}

// 判断新闻nid是不是用户uid发布的
func newsBelongUser(nid, uid int) bool {
	news := database.GetNewsById(nid)
	if news != nil {
		if news.UserId == uid {
			return true
		}
	}
	return false
}
