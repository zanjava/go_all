package main

import (
	"context"
	"fmt"
	"my_goframe/internal/model/do"

	_ "github.com/gogf/gf/contrib/drivers/mysql/v2" //注册mysql driver
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcfg"
	"github.com/gogf/gf/v2/os/gctx"
)

func Mysql() {
	// 指定mysql配置文件
	g.Cfg().GetAdapter().(*gcfg.AdapterFile).SetFileName("conf/learn.yaml")
	ctx := gctx.New()
	// Create(ctx)
	// Read(ctx)
	// Update(ctx)
	// Delete(ctx)

	Raw(ctx)
}

// 插入insert
func Create(ctx context.Context) {
	inst := do.User{
		Name:     "jane",
		Password: "12345678901234567890123456789012",
	}
	lastInsertId, err := g.
		DB("post").    //库名
		Model("user"). //表名
		Data(inst).    //插入的数据
		InsertAndGetId()
	if err != nil {
		fmt.Printf("insert failed: %s\n", err)
	} else {
		fmt.Printf("id %d\n", lastInsertId)
	}
}

// 读取select
func Read(ctx context.Context) {
	//查询1个
	var user do.User
	err := g.DB("post").Model("user").
		Where("name=?", "jane"). //where条件
		Scan(&user)              //查询结果赋给一个结构体
	if err != nil {
		fmt.Printf("select failed: %s\n", err)
	} else {
		fmt.Printf("read user %+v\n", user)
	}
	fmt.Println()

	//查询多个
	var users []do.User
	err = g.DB("post").Model("user").
		Where("id>?", 0).         //where条件
		OrderDesc("create_time"). //降序排列
		Limit(4).                 // limit offset
		Offset(2).
		Scan(&users) //查询结果赋给一个切片
	if err != nil {
		fmt.Printf("select failed: %s\n", err)
	} else {
		for _, user := range users {
			fmt.Printf("read user %+v\n", user)
		}
	}
}

// 更新update
func Update(ctx context.Context) {
	newPass := "01234567890123456789012345678901"
	result, err := g.DB("post").Model("user").
		Where("name=?", "jane"). //where条件
		// Data(do.User{Password: newPass}). //结构体里的非空字段都会被更新
		Data(g.Map{"password": newPass}). //或者使用map指定更新后的内容
		Update()
	if err != nil {
		fmt.Printf("update failed: %s\n", err)
	} else {
		n, _ := result.RowsAffected()
		fmt.Printf("update %d rows\n", n)
	}
}

// 删除 delete
func Delete(ctx context.Context) {
	result, err := g.DB("post").Model("user").Where("name=?", "jane").Delete()
	if err != nil {
		fmt.Printf("delete failed: %s\n", err)
	} else {
		n, _ := result.RowsAffected()
		fmt.Printf("delete %d rows\n", n)
	}
}

// 原生SQL语句
func Raw(ctx context.Context) {
	result, err := g.DB("test").Query(ctx, "select city,avg(id) as mid from user where id>0 group by city having mid>0")
	if err != nil {
		fmt.Printf("raw query failed: %s\n", err)
	} else {
		fmt.Println("city\tmid")
		for _, record := range result {
			fmt.Println(record["city"], record["mid"])
		}
	}
}
