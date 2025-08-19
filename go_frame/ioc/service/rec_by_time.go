package service

import (
	"go/frame/ioc/database"
	"go/frame/ioc/database/model"
	"sort"
)

type RecByTime struct {
	Useless int               // 没加container Tag，表示不需要由IoC容器提供依赖
	Db      database.GormNews `container:"type"`
}

func (r *RecByTime) RecNews() []*model.News {
	news := r.Db.GetNews()
	sort.Slice(news, func(i, j int) bool {
		return news[i].PostTime.After(*news[j].PostTime)
	})
	return news
}
