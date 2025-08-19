package database

import (
	"log"
	"log/slog"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"xorm.io/xorm"
)

var (
	mysqlDB     *xorm.Engine
	mysqlDBOnce sync.Once
)

func GetMysqlDB() (*xorm.Engine, func()) {
	mysqlDBOnce.Do(func() {
		if mysqlDB == nil {
			DataSourceName := "tester:123456@tcp(localhost:3306)/post?charset=utf8mb4&parseTime=True&loc=Local"
			db, err := xorm.NewEngine("mysql", DataSourceName)
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
		mysqlDB.Ping()
		slog.Info("ping post db")
	}
}

// 关闭数据库连接
func CloseMysqlDB() {
	if mysqlDB != nil {
		mysqlDB.Close()
		log.Println("关闭xorm数据库连接")
	}
}
