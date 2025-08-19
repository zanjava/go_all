package service

import (
	"go/frame/dependency_injection/database"
	"go/frame/dependency_injection/database/model"
	"sort"
)

type RecByPoster struct {
	Useless int              `wire:"-"` //未使用
	Db      database.INewsDB //把依赖注入到结构体里去。 如果不使用构造函数，则成员必须是可导出的
}

// 这种非自己主动初始化依赖，而通过外部来传入依赖的方式，称为依赖注入
func NewRecByPoster(db database.INewsDB) *RecByPoster {
	return &RecByPoster{
		Db: db,
	}
}

func (r *RecByPoster) RecNews() []*model.News {
	news := r.Db.GetNews()
	sort.Slice(news, func(i, j int) bool {
		return news[i].UserId > news[j].UserId
	})
	return news
}
