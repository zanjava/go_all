package service

import (
	"go/frame/ioc/database"
	"go/frame/ioc/database/model"
	"sort"
)

type RecByPoster struct {
	Useless int               // 没加container Tag，表示不需要由IoC容器提供依赖
	Db      database.GormNews `container:"type"`
}

func (r *RecByPoster) RecNews() []*model.News {
	news := r.Db.GetNews()
	sort.Slice(news, func(i, j int) bool {
		return news[i].UserId > news[j].UserId
	})
	return news
}
