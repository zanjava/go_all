package database

import (
	"errors"
	"fmt"
	"go/post/database/model"
	"log/slog"
)

// 新闻发布
func PostNews(uid int, title, content string) (int, error) {
	news := model.News{
		UserId:  uid,
		Title:   title,
		Content: content,
	}
	_, err := PostDB.Insert(&news)
	if err != nil {
		slog.Error("新闻发布失败", "title", title, "error", err)
		return 0, errors.New("新闻发布失败，请稍后重试")
	}
	return news.Id, nil
}

// 删除新闻
func DeleteNews(id int) error {
	affected, err := PostDB.ID(int64(id)).Delete(model.News{})
	if err != nil {
		slog.Error("DeleteNews failed", "id", id, "error", err)
		return errors.New("新闻删除失败，请稍后重试")
	}
	if affected <= 0 {
		return fmt.Errorf("新闻id[%d]不存在", id)
	}
	return nil
}

// 修改新闻
func UpdateNews(id int, title, content string) error {
	affected, err := PostDB.Where("id=?", id).Update(model.News{Title: title, Content: content})
	if err != nil {
		slog.Error("UpdateNews failed", "id", id, "error", err)
		return errors.New("新闻修改失败，请稍后重试")
	}
	if affected <= 0 {
		return fmt.Errorf("新闻id[%d]不存在", id)
	}
	return nil
}

func GetNewsById(id int) *model.News {
	news := &model.News{Id: id}
	ok, err := PostDB.Get(news) //隐含的where条件是id
	if err != nil {
		slog.Error("GetNewsById failed", "id", id, "error", err)
		return nil
	}
	if !ok {
		return nil
	}
	news.ViewPostTime = news.PostTime.Format("2006-01-02 15:04:05")
	return news
}

func GetNewsByUid(uid int) []*model.News {
	var news []*model.News
	err := PostDB.Where("user_id=?", uid).Find(&news)
	if err != nil {
		slog.Error("GetNewsByUid failed", "uid", uid, "error", err)
	}
	return news
}

// pageNo从1开始编号
func GetNewsByPage(pageNo, pageSize int) (int, []*model.News) {
	var total int64
	var err error
	total, err = PostDB.Count(model.News{})
	if err != nil {
		slog.Error("get news count failed", "error", err)
		return 0, nil
	}

	var news []*model.News
	// 新闻按发布时间降序排列
	err = PostDB.Desc("create_time").Limit(pageSize, pageSize*(pageNo-1)).Find(&news)
	if err != nil {
		slog.Error("GetNewsByPage failed", "pageNo", pageNo, "pageSize", pageSize, "error", err)
	}
	if len(news) > 0 {
		for _, ele := range news {
			ele.ViewPostTime = ele.PostTime.Format("2006-01-02 15:04:05")
		}
	}
	return int(total), news
}
