package database

import (
	"errors"
	"fmt"
	"go/post/database/model"
	"log/slog"
	"time"

	"gorm.io/gorm"
)

// 新闻发布
func PostNews(uid int, title, content string) (int, error) {
	now := time.Now()
	news := &model.News{
		UserId:     uid,
		Title:      title,
		Content:    content,
		PostTime:   &now,
		DeleteTime: nil, //写入数据库就是null
	}
	err := PostDB.Create(news).Error
	if err != nil {
		slog.Error("新闻发布失败", "title", title, "error", err)
		return 0, errors.New("新闻发布失败，请稍后重试")
	}
	return news.Id, nil
}

// 删除新闻
func DeleteNews(id int) error {
	tx := PostDB.Model(&model.News{}).Where("id=? and delete_time is null", id).Update("delete_time", time.Now())
	if tx.Error != nil {
		slog.Error("DeleteNews failed", "id", id, "error", tx.Error)
		return errors.New("新闻删除失败，请稍后重试")
	} else {
		if tx.RowsAffected <= 0 {
			return fmt.Errorf("新闻id[%d]不存在", id)
		} else {
			return nil
		}
	}
}

// 修改新闻
func UpdateNews(id int, title, content string) error {
	tx := PostDB.Model(&model.News{}).Where("id=? and delete_time is null", id).Updates(map[string]any{"title": title, "article": content})
	if tx.Error != nil {
		slog.Error("UpdateNews failed", "id", id, "error", tx.Error)
		return errors.New("新闻修改失败，请稍后重试")
	} else {
		if tx.RowsAffected <= 0 {
			return fmt.Errorf("新闻id[%d]不存在", id)
		} else {
			return nil
		}
	}
}

func GetNewsById(id int) *model.News {
	news := &model.News{Id: id}
	tx := PostDB.Select("*").Where("delete_time is null").First(news) //隐含的where条件是id。注意：Find不会返回ErrRecordNotFound
	if tx.Error != nil {
		if !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			slog.Error("GetNewsById failed", "id", id, "error", tx.Error)
		}
		return nil
	}
	news.ViewPostTime = news.PostTime.Format("2006-01-02 15:04:05")
	return news
}

func GetNewsByUid(uid int) []*model.News {
	var news []*model.News
	tx := PostDB.Select("*").Where("delete_time is null and user_id=?", uid).Find(&news)
	if tx.Error != nil {
		slog.Error("GetNewsByUid failed", "uid", uid, "error", tx.Error)
	}
	return news
}

// pageNo从1开始编号
func GetNewsByPage(pageNo, pageSize int) (int, []*model.News) {
	var total int64
	err := PostDB.Model(model.News{}).Where("delete_time is null").Count(&total).Error
	if err != nil {
		slog.Error("get news count failed", "error", err)
		return 0, nil
	}

	var news []*model.News
	// 新闻按发布时间降序排列
	tx := PostDB.Select("*").Where("delete_time is null").Order("create_time desc").Limit(pageSize).Offset(pageSize * (pageNo - 1)).Find(&news)
	if tx.Error != nil {
		slog.Error("GetNewsByPage failed", "pageNo", pageNo, "pageSize", pageSize, "error", tx.Error)
	}
	if len(news) > 0 {
		for _, ele := range news {
			ele.ViewPostTime = ele.PostTime.Format("2006-01-02 15:04:05")
		}
	}
	return int(total), news
}
