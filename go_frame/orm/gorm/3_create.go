package gorm

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand/v2"
	"time"

	"gorm.io/gorm"
)

// 根据struct创建
func Create(db *gorm.DB) {
	//插入一条记录
	user := User{UserId: rand.IntN(100000), Degree: "本科", Gender: "男", City: "上海", Keywords: []string{"编程", "golang"}}
	result := db.Create(&user) //必须传指针，因为要给User的主键赋值。主键为0值时Create会自动给主键赋值
	if result.Error != nil {
		slog.Error("插入记录失败", "error", result.Error)
	}
	fmt.Printf("record id is %d\n", user.Id)
	fmt.Printf("影响%d行\n", result.RowsAffected)

	//会话模式
	tx := db.Session(&gorm.Session{SkipHooks: true}) //不执行钩子Hook
	// db := db.Session(&gorm.Session{DryRun: true}) //生成SQL，但不执行
	//一次性插入多条
	user1 := user                    //发生拷贝
	user1.Id = 0                     //把主键置为0
	user1.UserId = rand.IntN(100000) //UserId上有唯一性约束
	user2 := user                    //发生拷贝
	user2.Id = 0                     //把主键置为0
	user2.UserId = rand.IntN(100000)
	users := []*User{&user1, &user2} //切片里的元素也可以不是指针
	result = tx.Create(users)        //一条SQL插入所有数据
	fmt.Printf("影响%d行\n", result.RowsAffected)

	//量太大时分批插入（SQL语句的长度是有上限的，同时避免长时间阻塞）
	batchSize := 1 //通常为几百
	user3 := user
	user3.Id = 0
	user3.UserId = rand.IntN(100000)
	user4 := user3
	user4.Id = 0
	db.CreateInBatches([]*User{&user3, &user4}, batchSize) //一个批次一条SQL。且所有批次被放到一个事务中来执行。由于user4插不进去，所以user3也会回滚。但如果设置了SkipDefaultTransaction就没有事务
}

// 根据map创建
func CreateByMap(db *gorm.DB) {
	//插入一条记录
	db.Model(User{}).Create(map[string]any{
		"uid": rand.IntN(100000), "degree": "本科", "gender": "男", "city": "上海",
		"keywords": []string{"编程", "golang"}, "create_time": time.Now(), "update_time": time.Now(),
	})

	//一次性插入多条
	db.Model(User{}).Create([]map[string]any{
		{"uid": rand.IntN(100000), "degree": "本科", "gender": "男", "city": "北京", "create_time": time.Now(), "update_time": time.Now()},
		{"uid": rand.IntN(100000), "degree": "本科", "gender": "男", "city": "深圳", "create_time": time.Now(), "update_time": time.Now()},
	})
}

/*
如果任何钩子回调返回错误，GORM将停止后续的操作并回滚事务。

Create时钩子的执行时机：
// 开始事务
BeforeSave
BeforeCreate
// 关联前的 save
// 插入记录至 db
// 关联后的 save
AfterCreate
AfterSave
// 提交或回滚事务
*/

func (u *User) BeforeSave(db *gorm.DB) (err error) {
	db.Logger.Info(context.Background(), "exec hook BeforeSave")
	return nil
}

func (u *User) BeforeCreate(db *gorm.DB) (err error) {
	db.Logger.Info(context.Background(), "exec hook BeforeCreate")
	return nil
}

func (u *User) AfterCreate(db *gorm.DB) (err error) {
	db.Logger.Info(context.Background(), "exec hook AfterCreate")
	return nil
}

func (u *User) AfterSave(db *gorm.DB) (err error) {
	db.Logger.Info(context.Background(), "exec hook AfterSave")
	return nil
}
