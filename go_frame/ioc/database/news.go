package database

import (
	"go/frame/ioc/database/model"
	"log/slog"

	"gorm.io/gorm"
)

type GormNews struct {
	db *gorm.DB `container:"type"` // container Tag表示由IoC容器提供该依赖
}

func (n *GormNews) GetNews() []*model.News {
	var news []*model.News
	// db, _ := GetMysqlDB() // 临场抓起一个依赖，如果要换依赖，每一个“临场”的代码都要改
	db := n.db
	tx := db.Select("*").Where("delete_time is null").Limit(100).Find(&news)
	if tx.Error != nil {
		slog.Error("GetNews failed", "error", tx.Error)
	}
	if len(news) > 0 {
		for _, ele := range news {
			ele.ViewPostTime = ele.PostTime.Format("2006-01-02 15:04:05")
		}
	}
	return news
}
