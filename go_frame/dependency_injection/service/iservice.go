package service

import "go/frame/dependency_injection/database/model"

type INewsRecommender interface {
	RecNews() []*model.News
}
