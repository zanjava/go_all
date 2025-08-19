package database

import "go/frame/dependency_injection/database/model"

type INewsDB interface {
	GetNews() []*model.News
}
