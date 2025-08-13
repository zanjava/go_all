package gorm

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 构建结构何体，跟表结构进行对应
type Login struct {
	Id       int `gorm:"primaryKey;column:id"` // 显式指定主键
	Username string
	Password string
}

// 指定结构体对应的表名是什么
func (Login) TableName() string {
	return "login"
}

func GormQuickStart() {
	//连接数据库
	host := "localhost"
	port := 3306
	dbname := "test"
	user := "tester"
	pass := "123456"
	// data source name
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, pass, host, port, dbname)
	db, err := gorm.Open(mysql.Open(dsn), nil)
	if err != nil {
		panic(err)
	}

	// 所谓"orm"即不需要关心sql语句，把注意力集中到结构体上面来
	// 写入
	instance1 := Login{Username: "zgw1", Password: "123456"}
	db.Create(&instance1)
	// 查询
	var instance2 Login
	db.Find(&instance2)
	fmt.Printf("%#v\n", instance2)
}
