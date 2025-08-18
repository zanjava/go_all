package main

import (
	"fmt"
	"strconv"
	"time"

	_ "github.com/gogf/gf/contrib/nosql/redis/v2" //注册redis adapter
	"github.com/gogf/gf/v2/database/gredis"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/os/gcfg"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/util/gconv"
)

// 本地缓存
func LocalCache() {
	ctx := gctx.New()       // 用context.TODO()或Background()也行
	cache := gcache.New(10) // 10是缓存的容量上限（采用LRU淘汰策略）
	defer cache.Close(ctx)  // 关闭缓存对象，让GC回收资源

	// 写缓存
	err := cache.Set(ctx, "k1", "v1", 0) //0表示永不过期
	if err != nil {
		panic(err)
	}

	// 读缓存
	value, err := cache.Get(ctx, "k1")
	if err != nil {
		panic(err)
	}
	fmt.Printf("value=%s\n", value)

	// 再添加20个元素
	for i := 0; i < 20; i++ {
		s := strconv.Itoa(i)
		cache.Set(ctx, s, s, time.Second) //1秒后过期
	}

	// 获取缓存大小
	size, err := cache.Size(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Printf("size=%d\n", size) //应该是10

	// 遍历缓存
	mp, err := cache.Data(ctx)
	if err != nil {
		panic(err)
	}
	for k, v := range mp {
		fmt.Println(k, v)
	}

	time.Sleep(time.Second)

	// 获取缓存大小
	size, err = cache.Size(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Printf("size=%d\n", size) //应该是0

}

// 分布式缓存
func DistributeCache() {
	// 指定配置文件
	g.Cfg().GetAdapter().(*gcfg.AdapterFile).SetFileName("conf/learn.yaml")

	ctx := gctx.New()
	addr, _ := g.Cfg().Get(ctx, "redis.addr")
	db, _ := g.Cfg().Get(ctx, "redis.db")
	// fmt.Println(addr, db)
	redisConfig := &gredis.Config{
		Address: gconv.String(addr), //gvar.Var转为具体类型（推荐方式）
		Db:      gconv.Int(db),      //或者先调用Val()转为any，再断言为具体类型也行
	}

	redis, err := gredis.New(redisConfig)
	if err != nil {
		panic(err)
	}
	cache := gcache.New()
	cache.SetAdapter(gcache.NewAdapterRedis(redis))

	key := "age"
	cache.Set(ctx, key, 18, time.Second) //写入，带超时
	value, _ := cache.Get(ctx, key)      //读取
	fmt.Printf("value=%d\n", gconv.Int(value))

	time.Sleep(time.Second)

	// 判断key是否存在
	exists, _ := cache.Contains(ctx, key)
	fmt.Printf("key存在吗？ %t\n", exists)
}
