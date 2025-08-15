package handler

import (
	database "go/post/database/gorm"
	"go/post/handler/model"
	"go/post/util"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 新闻发布
func PostNews(ctx *gin.Context) {
	loginUid := ctx.Value(UID_IN_CTX).(int)
	var news model.News
	err := ctx.ShouldBind(&news)
	if err != nil {
		ctx.String(http.StatusBadRequest, util.BindErrMsg(err))
		return
	}

	id, err := database.PostNews(loginUid, news.Title, news.Content)
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"id": id})
}

// 新闻详情页
func GetNewsById(ctx *gin.Context) {
	idStr := ctx.Param("id")
	if id, err := strconv.Atoi(idStr); err != nil || id <= 0 {
		ctx.String(http.StatusBadRequest, "非法的新闻id")
		return
	} else {
		news := database.GetNewsById(id)
		if news == nil {
			ctx.Status(http.StatusNotFound)
			return
		}

		user := database.GetUserById(news.UserId)
		if user != nil {
			news.UserName = user.Name
		}
		ctx.HTML(http.StatusOK, "news.html", news)
		return
	}
}

// 删除新闻
func DeleteNews(ctx *gin.Context) {
	loginUid := ctx.Value(UID_IN_CTX).(int)
	idStr := ctx.Param("id")
	if id, err := strconv.Atoi(idStr); err != nil || id <= 0 {
		ctx.String(http.StatusBadRequest, "非法的新闻id")
		return
	} else {
		if !newsBelongUser(id, loginUid) {
			ctx.String(http.StatusForbidden, "无权限删除")
			return
		}
		err = database.DeleteNews(id)
		if err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}
		ctx.Status(http.StatusOK)
	}
}

// 修改新闻
func UpdateNews(ctx *gin.Context) {
	loginUid := ctx.Value(UID_IN_CTX).(int)
	var news model.News
	err := ctx.ShouldBind(&news)
	if err != nil {
		ctx.String(http.StatusBadRequest, util.BindErrMsg(err))
		return
	}
	if news.Id <= 0 {
		ctx.String(http.StatusBadRequest, "非法的新闻id")
		return
	}

	if !newsBelongUser(news.Id, loginUid) {
		ctx.String(http.StatusForbidden, "无权限修改")
		return
	}

	err = database.UpdateNews(news.Id, news.Title, news.Content)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.Status(http.StatusOK)
}

// 新闻列表页
func NewsList(ctx *gin.Context) {
	PageNo, err1 := strconv.Atoi(ctx.DefaultQuery("page_no", "1"))
	PageSize, err2 := strconv.Atoi(ctx.DefaultQuery("page_size", "3"))
	if err1 != nil || err2 != nil {
		ctx.String(http.StatusBadRequest, "page参数不合法")
		return
	}

	total, news := database.GetNewsByPage(PageNo, PageSize)
	if total < len(news) {
		ctx.String(http.StatusInternalServerError, "出错了")
		return
	}
	for _, ele := range news {
		user := database.GetUserById(ele.UserId)
		if user != nil {
			ele.UserName = user.Name
		} else {
			slog.Warn("could not get name of user", "uid", ele.UserId)
		}
	}
	ctx.HTML(http.StatusOK, "news_list.html", gin.H{"total": total, "data": news, "page": PageNo})
}

// 从jwt里解析出uid，判断news id是否属于uid
func NewsBelong(ctx *gin.Context) {
	newsId := ctx.Query("id") //新闻id
	nid, err := strconv.Atoi(newsId)
	if err != nil {
		ctx.String(http.StatusBadRequest, "invalid news id")
		return
	}

	loginUid := GetLoginUid(ctx)
	//loginUid := GetLoginUidFromSession(ctx)
	if newsBelongUser(nid, loginUid) {
		ctx.String(http.StatusOK, "true") //新闻的作者id就是当前登录者的uid
	} else {
		ctx.String(http.StatusOK, "false")
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
