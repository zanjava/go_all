package xorm

import (
	"fmt"
	"os"
	"time"

	"xorm.io/xorm"
	"xorm.io/xorm/log"
	"xorm.io/xorm/names"
)

func CreateEngine(host, dbname, user, pass string, port int) *xorm.Engine {
	// mb4兼容emoji表情符号。
	// 想要正确的处理time.Time ，您需要带上parseTime参数。
	// loc=Local采用机器本地的时区。
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, pass, host, port, dbname)
	engine, err := xorm.NewEngine("mysql", dsn) // data source name
	if err != nil {
		panic(err)
	}

	//SnakeMapper是默认的命名方式，驼峰转蛇形。
	//SameMapper不转。
	//GonicMapper 和SnakeMapper很类似，但是对于特定词支持更好，比如ID会翻译成id而不是i_d。
	engine.SetMapper(names.GonicMapper{})

	// // 表名统一带前缀。当然也可以使用NewSuffixMapper指定后缀
	// tableMapper := names.NewPrefixMapper(names.GonicMapper{}, "t_")
	// // 分别给表名和列名指定不同的名称映射方式
	// engine.SetTableMapper(tableMapper)
	// engine.SetColumnMapper(names.GonicMapper{})

	//日志控制
	logFile, _ := os.OpenFile("D:/go_all/go_frame/log/xorm.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
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

	return engine
}

// 读写分离，Engine分组
func CreateEngineGroup(host, dbname, user, pass string, port int) *xorm.EngineGroup {
	dsn1 := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, pass, host, port, dbname)
	dsn2 := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, pass, host, port, dbname)
	dsn3 := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, pass, host, port, dbname)

	eg, err := xorm.NewEngineGroup("mysql", []string{
		dsn1, //第一个是master。写操作命中master
		dsn2, //其余的都是slave。读操作命中slave(具体命中哪台取决于负载均衡策略)
		dsn3,
	})
	if err != nil {
		panic(err)
	}

	//创建完成engine之后，并没有立即连接数据库，此时可以通过Ping()来测试是否可以连接到数据库(会依次测试组里的每一个连接)
	if err = eg.Ping(); err != nil {
		panic(err)
	}

	// 设置负载均衡策略。Random、WeightRandom、RoundRobin、WeightRoundRobin等
	eg.SetPolicy(xorm.LeastConnPolicy()) //最小连接数

	return eg
}
