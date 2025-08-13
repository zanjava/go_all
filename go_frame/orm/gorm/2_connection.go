package gorm

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/dbresolver"
)

func CreateConnection(host, dbname, user, pass string, port int) *gorm.DB {
	//mb4兼容emoji表情符号。
	// 想要正确的处理time.Time ，您需要带上parseTime参数。
	// loc=Local采用机器本地的时区。
	DataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, pass, host, port, dbname)

	//日志控制
	logFile, _ := os.OpenFile("D:/go_all/go_frame/log/gorm.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
	newLogger := logger.New(
		log.New(logFile, "\r\n", log.LstdFlags), // io writer，可以输出到文件，也可以输出到os.Stdout
		logger.Config{
			SlowThreshold:             500 * time.Millisecond, //耗时超过此值认定为慢查询
			LogLevel:                  logger.Info,            // LogLevel的最低阈值，Silent为不输出日志
			IgnoreRecordNotFoundError: true,                   // 忽略RecordNotFound这种错误日志
			ParameterizedQueries:      true,                   // true代表SQL日志里不包含参数
			Colorful:                  false,                  // 禁用颜色
		},
	)
	db, err := gorm.Open(mysql.Open(DataSourceName), &gorm.Config{
		PrepareStmt:            true,  //执行任何SQL时都会创建一个prepared statement并将其缓存，以提高后续的效率
		SkipDefaultTransaction: false, // 为了确保数据一致性，GORM 会在事务里执行写入操作（创建、更新、删除）。如果没有这方面的要求，您可以在初始化时禁用它，这将获得大约 30%+ 性能提升。
		NamingStrategy: schema.NamingStrategy{ //覆盖默认的NamingStrategy来更改命名约定
			// TablePrefix:   "t_",                              // table name prefix, table for `User` would be `t_users`
			SingularTable: true, //表名映射时不加复数，仅是驼峰-->蛇形
			// NoLowerCase:   true,                              // skip the snake_casing of names
			// NameReplacer:  strings.NewReplacer("CID", "Cid"), // use name replacer to change struct/field name before convert it to db name
		},
		Logger:                   newLogger, //日志控制
		DryRun:                   false,     //true代表生成SQL但不执行，可以用于准备或测试生成的 SQL
		DisableNestedTransaction: true,      //在一个事务中使用Transaction方法，GORM会使用 SavePoint(savedPointName)，RollbackTo(savedPointName) 为你提供嵌套事务支持。如果不需要嵌套事务，可以将其禁用
		DisableAutomaticPing:     false,     //在完成初始化后，GORM 会自动ping数据库以检查数据库的可用性
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
	return db
}

// CreateRWDB 创建一个读写分离的数据库连接
func CreateRWDB(host, dbname, user, pass string, port int) *gorm.DB {
	dsn1 := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, pass, host, port, dbname)
	dsn2 := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, pass, host, port, dbname)
	dsn3 := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, pass, host, port, dbname)

	db, err := gorm.Open(
		mysql.Open(dsn1), //dsn1是主库
		&gorm.Config{},
	)
	if err != nil {
		panic(err)
	}

	replicas := []gorm.Dialector{
		mysql.New(mysql.Config{DSN: dsn2}), //dsn2和dsn3是从库
		mysql.New(mysql.Config{DSN: dsn3}),
	}
	resolver := dbresolver.Register(dbresolver.Config{
		Replicas: replicas,                      //复本，即从库
		Policy:   dbresolver.RoundRobinPolicy(), //或RandomPolicy{}
	})
	db.Use(resolver) //指定从库及其负载均衡策略

	return db
}
