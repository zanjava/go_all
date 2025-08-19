package database

import (
	"github.com/golobby/container/v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func init() {
	//分为Singleton和Transient两种模式
	err := container.Singleton(func() (*gorm.DB, error) { //告诉IoC容器如何创建bean
		DataSourceName := "tester:123456@tcp(localhost:3306)/post?charset=utf8mb4&parseTime=True&loc=Local"
		return gorm.Open(mysql.Open(DataSourceName))
	})
	if err != nil {
		panic(err)
	}

	err = container.Transient(func() (GormNews, error) { //告诉IoC容器如何创建bean
		var gn GormNews
		e := container.Fill(&gn)
		return gn, e
	})
	if err != nil {
		panic(err)
	}
}
