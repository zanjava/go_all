package xorm

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"xorm.io/xorm"
)

// 构建结构何体，跟表结构进行对应
type Login struct {
	Username string
	Password string
}

func XormQuickStart() {
	//连接数据库
	host := "localhost"
	port := 3306
	dbname := "test"
	user := "tester"
	pass := "123456"
	// data source name
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, pass, host, port, dbname)
	engine, err := xorm.NewEngine("mysql", dsn)
	if err != nil {
		panic(err)
	}
	defer engine.Close() //关闭Engine

	//创建完成engine之后，并没有立即连接数据库，此时可以通过Ping()来测试是否可以连接到数据库
	if err = engine.Ping(); err != nil {
		panic(err)
	}

	// 所谓"orm"即不需要关心sql语句，把注意力集中到结构体上面来
	// 写入
	instance1 := Login{Username: "dqq", Password: "123456"}
	engine.Insert(&instance1)

	// 查询
	var instance2 Login
	engine.Get(&instance2)
	fmt.Printf("%#v\n", instance2)
}
