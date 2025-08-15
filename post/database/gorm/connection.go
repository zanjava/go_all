package database

import (
	"fmt"
	"go/post/util"
	"log"
	"log/slog"
	"os"
	"path"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var (
	PostDB *gorm.DB
)

func ConnectPostDB(confDir, confFile, fileType, logDir string) {
	viper := util.InitViper(confDir, confFile, fileType)
	user := viper.GetString("post.user")
	pass := viper.GetString("post.pass")
	host := viper.GetString("post.host")
	port := viper.GetInt("post.port")
	dbname := "post"
	logFileName := viper.GetString("post.log")
	DataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, pass, host, port, dbname)

	//日志控制
	logFile, _ := os.OpenFile(path.Join(logDir, logFileName), os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
	newLogger := logger.New(
		log.New(logFile, "\r\n", log.LstdFlags), // io writer，可以输出到文件，也可以输出到os.Stdout
		logger.Config{
			SlowThreshold:             100 * time.Millisecond, //耗时超过此值认定为慢查询
			LogLevel:                  logger.Info,            // LogLevel的最低阈值，Silent为不输出日志
			IgnoreRecordNotFoundError: true,                   // 忽略RecordNotFound这种错误日志
			Colorful:                  false,                  // 禁用颜色
			ParameterizedQueries:      true,                   // 启用参数化查询
		},
	)
	db, err := gorm.Open(mysql.Open(DataSourceName), &gorm.Config{
		PrepareStmt:            true, //执行任何SQL时都会创建一个prepared statement并将其缓存，以提高后续的效率
		SkipDefaultTransaction: true, // 为了确保数据一致性，GORM 会在事务里执行写入操作（创建、更新、删除）。如果没有这方面的要求，您可以在初始化时禁用它，这将获得大约 30%+ 性能提升。
		NamingStrategy: schema.NamingStrategy{ //覆盖默认的NamingStrategy来更改命名约定
			SingularTable: true, //表名映射时不加复数，仅是驼峰-->蛇形
		},
		Logger: newLogger, //日志控制
	})
	if err != nil {
		panic(err)
	}

	//连接池控制参数
	sqlDB, _ := db.DB()
	//池子里空闲连接的数量上限（超出此上限就把相应的连接关闭掉）
	sqlDB.SetMaxIdleConns(10)
	//最多开这么多连接
	sqlDB.SetMaxOpenConns(100)
	//一个连接最多可使用这么长时间，超时后连接会自动关闭（因为数据库本身可能也对NoActive连接设置了超时时间，我们的应对办法：定期ping，或者SetConnMaxLifetime）
	sqlDB.SetConnMaxLifetime(time.Hour)
	PostDB = db
}

// 定期ping，保持连接的活跃
func PingPostDB() {
	if PostDB != nil {
		sqlDB, _ := PostDB.DB()
		sqlDB.Ping()
		slog.Info("ping post db")
	}
}

// 关闭数据库连接
func ClosePostDB() {
	if PostDB != nil {
		sqlDB, _ := PostDB.DB()
		sqlDB.Close()
	}
}
