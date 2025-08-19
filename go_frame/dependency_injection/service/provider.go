package service

import (
	"github.com/google/wire"
)

var (
	//service.RecByTime依赖的INewsDB是一个接口，需要通过wire.Bind()指定INewsDB接口的实现是谁
	SvrSet = wire.NewSet(NewRecByTime, wire.Bind(new(INewsRecommender), new(*RecByTime)))
	// SvrSet = wire.NewSet(
	// 	// wire.Struct(new(RecByTime), "Db"), //指定结构体的哪个字段需要注入
	// 	wire.Struct(new(RecByTime), "*"), //*表示结构体的所有字段都需要注入，对于不需要注入的字段可以加个Tag: `wire:"-"`
	// 	wire.Bind(new(INewsRecommender), new(*RecByTime)))
)
