package database

import (
	"fmt"
	"go/post/util"
	"log/slog"
	"os"
	"path"
	"time"

	"xorm.io/xorm"
	"xorm.io/xorm/log"
	"xorm.io/xorm/names"
)

var (
	PostDB *xorm.Engine
)

func ConnectPostDB(confDir, confFile, fileType, logDir string) {
	viper := util.InitViper(confDir, confFile, fileType)
	user := viper.GetString("post.user")
	pass := viper.GetString("post.pass")
	host := viper.GetString("post.host")
	port := viper.GetInt("post.port")
	dbname := "post"
	logFileName := viper.GetString("post.log")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, pass, host, port, dbname)
	engine, err := xorm.NewEngine("mysql", dsn) // data source name
	if err != nil {
		panic(err)
	}

	engine.SetMapper(names.GonicMapper{})

	//日志控制
	logFile, _ := os.OpenFile(path.Join(logDir, logFileName), os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
	logger := log.NewSimpleLogger(logFile)
	logger.ShowSQL(true)             //日志中显示SQL语句
	engine.SetLogger(logger)         //日志输出到指定文件
	engine.SetLogLevel(log.LOG_INFO) //日志级别

	//连接池控制参数
	engine.SetMaxIdleConns(10)           //池子里空闲连接的数量上限（超出此上限就把相应的连接关闭掉）
	engine.SetMaxOpenConns(100)          //最多开这么多连接
	engine.SetConnMaxLifetime(time.Hour) //一个连接最多可使用这么长时间，超时后连接会自动关闭（因为数据库本身可能也对NoActive连接设置了超时时间，我们的应对办法：定期ping，或者SetConnMaxLifetime）

	//创建完成engine之后，并没有立即连接数据库，此时可以通过Ping()来测试是否可以连接到数据库
	if err = engine.Ping(); err != nil {
		panic(err)
	}

	PostDB = engine
}

// 定期ping，保持连接的活跃
func PingPostDB() {
	if PostDB != nil {
		PostDB.Ping()
		slog.Info("ping post db")
	}
}

// 关闭数据库连接
func ClosePostDB() {
	if PostDB != nil {
		PostDB.Close()
	}
}
