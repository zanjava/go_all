package database

import (
	gd "go/frame/dependency_injection/database/gorm"

	"github.com/google/wire"
)

var (
	//gd.GetMysqlDB返回的是结构体，直接写上就可以了
	DbSet = wire.NewSet(gd.GetMysqlDB, gd.NewGormNews, wire.Bind(new(INewsDB), new(*gd.GormNews)))
)
