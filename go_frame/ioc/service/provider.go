package service

import (
	"log"

	"github.com/golobby/container/v3"
)

func init() {
	// 指定如何创建接口的具体实现。打算同时使用多种具体实现时，可以通过Name进行区分
	err := container.NamedTransient("sort_time", func() INewsRecommender { //Transient每次需要时就创建一个全新的实例
		log.Printf("创建RecByTime实例")
		var recByTime RecByTime
		if err := container.Fill(&recByTime); err != nil { //让IoC容器提供一个bean(容器只负责给带container Tag的Field赋值)
			panic(err)
		} else {
			recByTime.Useless = 4 //手动给不带container Tag的Field赋值
			return &recByTime     // 在执行init()时这个函数会被调用。执行Resolve()时会再次被调用
		}
	})
	if err != nil {
		panic(err)
	}

	err = container.NamedSingleton("sort_poster", func() INewsRecommender { //Singleton全局只创建一个实例
		log.Printf("创建RecByPoster实例")
		var recByPoster RecByTime
		err := container.Fill(&recByPoster)
		if err != nil {
			panic(err)
		}
		recByPoster.Useless = 4
		return &recByPoster // 在执行init()时这个函数会被调用。执行Resolve()时不再被调用（因为是Singleton模式）
	})
	if err != nil {
		panic(err)
	}
}
