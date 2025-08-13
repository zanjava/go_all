package orm_test

import (
	"fmt"
	"math/rand/v2"
	"testing"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"xorm.io/xorm"
	"xorm.io/xorm/log"
)

var (
	host   = "localhost"
	port   = 3306
	dbname = "test"
	user   = "tester"
	pass   = "123456"
)

type User struct {
	Id        int `gorm:"primaryKey;column:id" xorm:"pk"`
	UserId    int `gorm:"column:uid" xorm:"uid"`
	Degree    string
	CreatedAt time.Time `gorm:"column:create_time" xorm:"created create_time"`
	UpdatedAt time.Time `gorm:"column:update_time" xorm:"updated update_time"`
	Gender    string
	City      string
	Province  string `gorm:"-" xorm:"-"`
}

func (User) TableName() string {
	return "user"
}

func InitGorm() *gorm.DB {
	DataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, pass, host, port, dbname)

	newLogger := logger.New(nil, logger.Config{LogLevel: logger.Silent}) //Silent禁用Logger
	db, err := gorm.Open(mysql.Open(DataSourceName), &gorm.Config{
		PrepareStmt:            true,                                       //启用 prepare statment
		SkipDefaultTransaction: true,                                       //写操作不用事务
		NamingStrategy:         schema.NamingStrategy{SingularTable: true}, //表名用单数
		Logger:                 newLogger,                                  //禁用Logger
	})
	if err != nil {
		panic(err)
	}

	//连接池控制参数
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
	return db
}

func InitXorm() *xorm.Engine {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, pass, host, port, dbname)
	engine, err := xorm.NewEngine("mysql", dsn)
	if err != nil {
		panic(err)
	}
	engine.SetLogLevel(log.LOG_OFF) //禁用Logger
	//连接池控制参数
	engine.SetMaxIdleConns(10)
	engine.SetMaxOpenConns(100)
	engine.SetConnMaxLifetime(time.Hour)
	//先Ping一下
	if err = engine.Ping(); err != nil {
		panic(err)
	}
	return engine
}

func TestGorm(t *testing.T) {
	db := InitGorm()
	db.Where("1=1").Delete(User{}) //先把表清空
	const C = 100
	for i := 0; i < C; i++ {
		user := User{UserId: rand.IntN(1000000000), Degree: "本科", Gender: "男", City: "上海"}
		db = db.Create(&user)
		if db.Error != nil {
			t.Fatal(db.Error)
		}
		if db.RowsAffected != 1 {
			t.Fatalf("插入了%d行", db.RowsAffected)
		}
	}

	var users []User
	db = db.Where("degree in ?", []string{"本科", "硕士"}).Where("gender=?", "男").Where("city=?", "上海").Limit(C).Find(&users)
	if db.Error != nil {
		t.Fatal(db.Error)
	}
	if db.RowsAffected != C {
		t.Fatalf("检索到%d行", db.RowsAffected)
	}
	fmt.Println(users[0])
}

func TestXorm(t *testing.T) {
	engine := InitXorm()
	engine.Where("1=1").Delete(User{}) //先把表清空
	const C = 100
	for i := 0; i < C; i++ {
		user := User{UserId: rand.IntN(1000000000), Degree: "本科", Gender: "男", City: "上海"}
		affected, err := engine.Insert(user)
		if err != nil {
			t.Fatal(err)
		}
		if affected != 1 {
			t.Fatalf("插入了%d行", affected)
		}
	}

	session := engine.Prepare()
	var users []User
	err := session.In("degree", []string{"本科", "硕士"}).And("gender=?", "男").Where("city=?", "上海").Limit(C).Find(&users)
	if err != nil {
		t.Fatal(err)
	}
	if len(users) != C {
		t.Fatalf("检索到%d行", len(users))
	}
	fmt.Println(users[0])
}

func BenchmarkCreateGorm(b *testing.B) {
	db := InitGorm()
	db.Where("1=1").Delete(User{}) //先把表清空
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		user := User{UserId: rand.IntN(1000000000), Degree: "本科", Gender: "男", City: "上海"}
		db.Create(&user)
	}
}

func BenchmarkCreateXorm(b *testing.B) {
	engine := InitXorm()
	engine.Where("1=1").Delete(User{}) //先把表清空
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		user := User{UserId: rand.IntN(1000000000), Degree: "本科", Gender: "男", City: "上海"}
		engine.Insert(user)
	}
}

func BenchmarkReadGorm(b *testing.B) {
	db := InitGorm()
	const C = 100
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var users []User
		db.Where("degree in ?", []string{"本科", "硕士"}).Where("gender=?", "男").Where("city=?", "上海").Limit(C).Find(&users)
	}
}

func BenchmarkReadXorm(b *testing.B) {
	engine := InitXorm()
	const C = 100
	session := engine.Prepare() //启用PrepareStmt
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var users []User
		session.In("degree", []string{"本科", "硕士"}).Where("gender=?", "男").Where("city=?", "上海").Limit(C).Find(&users)
	}
}

// go test -v ./orm -run=^TestGorm$ -count=1
// go test -v ./orm -run=^TestXorm$ -count=1

// go test -v ./orm -bench=^BenchmarkCreate -run=^$ -benchtime=10s -count=1
// go test -v ./orm -bench=^BenchmarkRead -run=^$ -benchtime=10s -count=1

/**
goos: windows
goarch: amd64
pkg: dqq/go/frame/orm
cpu: 11th Gen Intel(R) Core(TM) i7-1165G7 @ 2.80GHz

单条插入性能（gorm略慢）
BenchmarkCreateGorm-8               1416           7524891 ns/op
BenchmarkCreateXorm-8               1812           6579417 ns/op

启用PrepareStmt，批量读取性能（gorm优势明显）
BenchmarkReadGorm-8        33554            347008 ns/op
BenchmarkReadXorm-8        20244            589800 ns/op

关闭PrepareStmt，批量读取性能（gorm优势明显）
BenchmarkReadGorm-8        29090            425082 ns/op
BenchmarkReadXorm-8        19848            625568 ns/op
*/
