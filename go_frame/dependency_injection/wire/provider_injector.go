//go:build wireinject

// 在本go文件的同目录下执行wire命令，会生成wire_gen.go文件
// go get github.com/google/wire@latest
// go install github.com/google/wire/cmd/wire@latest
package main

import (
	"go/frame/dependency_injection/database"
	"go/frame/dependency_injection/handler"
	"go/frame/dependency_injection/service"

	"github.com/google/wire"
)

// 定义Provider
// var (
// 	//gd.GetMysqlDB返回的是结构体，直接写上就可以了
// 	DbSet = wire.NewSet(gd.GetMysqlDB, gd.NewGormNews, wire.Bind(new(database.INewsDB), new(*gd.GormNews)))
// 	//service.RecByTime依赖的INewsDB是一个接口，需要通过wire.Bind()指定INewsDB接口的实现是谁
// 	// SvrSet = wire.NewSet(service.NewRecByTime, wire.Bind(new(service.INewsRecommender), new(*service.RecByTime)))

// 	//不使用构造函数，直接使用结构体创建实例
// 	SvrSet = wire.NewSet(
// 		// wire.Struct(new(service.RecByTime), "Db"), //指定结构体的哪个字段需要注入
// 		wire.Struct(new(service.RecByTime), "*"), //*表示结构体的所有字段都需要注入，对于不需要注入的字段可以加个Tag: `wire:"-"`
// 		wire.Bind(new(service.INewsRecommender), new(*service.RecByTime)))

// 	HdlSet = wire.NewSet(handler.NewGinHandler)

// 	// ProviderSet = wire.NewSet(DbSet, SvrSet, HdlSet) //参数顺序无所谓
// )

// 定义Injector（一定要和Provider分开写，不能把Provider的定义也放到这个函数里）
func InitHandler() (h *handler.GinHandler, cp func(), err error) { //不用自己去构造这3个返回值。第2个返回值是cleanup函数，用于释放资源，所以依赖的cleanup都会【有序】地合并到这个函数里。第2个、第3个返回值如果没有可以不写
	// wire.Build(ProviderSet)
	// wire.Build(DbSet, SvrSet, HdlSet) // 等价于上面一行
	wire.Build(service.SvrSet, handler.HdlSet, database.DbSet) //参数顺序无所谓

	return
}
