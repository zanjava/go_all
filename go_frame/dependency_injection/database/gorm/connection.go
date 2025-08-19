package database

import (
	"log"
	"log/slog"
	"sync"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	mysqlDB     *gorm.DB
	mysqlDBOnce sync.Once
)

func GetMysqlDB() (*gorm.DB, func()) {
	mysqlDBOnce.Do(func() {
		if mysqlDB == nil {
			DataSourceName := "tester:123456@tcp(localhost:3306)/post?charset=utf8mb4&parseTime=True&loc=Local"
			db, err := gorm.Open(mysql.Open(DataSourceName))
			if err != nil {
				panic(err)
			}
			mysqlDB = db
		}
	})
	return mysqlDB, CloseMysqlDB
}

// 定期ping，保持连接的活跃
func PingMysqlDB() {
	if mysqlDB != nil {
		sqlDB, _ := mysqlDB.DB()
		sqlDB.Ping()
		slog.Info("ping post db")
	}
}

// 关闭数据库连接
func CloseMysqlDB() {
	if mysqlDB != nil {
		sqlDB, _ := mysqlDB.DB()
		sqlDB.Close()
		log.Println("关闭gorm数据库连接")
	}
}
