package service

import "go/frame/ioc/database/model"

type INewsRecommender interface {
	RecNews() []*model.News
}
