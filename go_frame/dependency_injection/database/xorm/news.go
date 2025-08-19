package database

import (
	"go/frame/dependency_injection/database/model"
	"log/slog"

	"xorm.io/xorm"
)

type XormNews struct {
	db *xorm.Engine // 把依赖注入到结构体里去
}

// 这种非自己主动初始化依赖，而通过外部来传入依赖的方式，称为依赖注入
func NewXormNews(db *xorm.Engine) *XormNews {
	return &XormNews{
		db: db,
	}
}

func (n *XormNews) GetNews() []*model.News {
	var news []*model.News
	// db := GetMysqlDB()   // 临场抓起一个依赖，如果要换依赖，每一个“临场”的代码都要改
	db := n.db
	err := db.Limit(100).Find(&news)
	if err != nil {
		slog.Error("GetNews failed", "error", err)
	}
	if len(news) > 0 {
		for _, ele := range news {
			ele.ViewPostTime = ele.PostTime.Format("2006-01-02 15:04:05")
		}
	}
	return news
}
