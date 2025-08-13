package gorm

import (
	"errors"
	"fmt"
	"log/slog"

	"gorm.io/gorm"
	"gorm.io/hints"
)

// 查询
func Read(db *gorm.DB) {
	user := User{City: "HongKong", Id: 3} //Id会自动放到Where条件里，其他非0字段不会
	tx := db.
		Select("uid,city,gender,keywords").       //参数也可以这样传"uid","city","gender"或者[]string{"uid","city","gender"}。没有Select时默认为select *
		Where("uid>100 and degree='大专'").         //容易发生SQL注入攻击
		Where("city in ?", []string{"北京", "上海"}). //多个Where之间是and关系
		Where("degree like ?", "%科").
		Or("gender=?", "女"). //用?占位，避免发生SQL注入攻击
		Order("id desc, uid").
		Order("city").
		Offset(3).
		Limit(1).
		First(&user) //Find可以传一个结构体，也可以传结构体切片。Take、First、Last查不到结果时会返回gorm.ErrRecordNotFound，但Find不会，Find查无结果时就不去修改结构体
	if tx.Error != nil {
		if !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			slog.Error("读DB失败", "error", tx.Error)
		} else {
			slog.Info("查无结果")
		}
	} else {
		if tx.RowsAffected > 0 {
			fmt.Printf("read结果：%+v\n", user)
		} else {
			slog.Info("查无结果", "user", user)
		}
	}

	var user2 *User //不同于var user2 User，还没申请内存空间
	// 通过反射给user2赋值时发现还没给user2申请好内存空间
	tx = db.Find(user2) //error: invalid value, should be pointer to struct or slice
	if tx.Error != nil {
		if !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			slog.Error("读DB失败", "error", tx.Error)
		} else {
			slog.Info("查无结果")
		}
	}

	var user3 *User = new(User)
	tx = db.Find(user3)
	if tx.Error != nil {
		if !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			slog.Error("读DB失败", "error", tx.Error)
		} else {
			slog.Info("查无结果")
		}
	} else {
		if tx.RowsAffected > 0 {
			fmt.Printf("read结果：%+v\n", user3)
		} else {
			slog.Info("查无结果", "user", user3)
		}
	}

	var users []User
	tx = db.Limit(3).Find(&users) //要修改切片的长度，所以要传切片的指针
	if tx.Error != nil {
		if !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			slog.Error("读DB失败", "error", tx.Error)
		} else {
			slog.Info("查无结果")
		}
	} else {
		if tx.RowsAffected > 0 {
			fmt.Println("多个read结果")
			for _, u := range users {
				fmt.Printf("%+v\n", u)
			}
		} else {
			slog.Info("查无结果")
		}
	}

	user4 := User{Id: 23212} //给主键赋值
	tx = db.Find(&user4)     //主键不为0值时暗含了一个where条件：id=47
	if tx.Error != nil {
		if !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			slog.Error("读DB失败", "error", tx.Error)
		} else {
			slog.Info("查无结果")
		}
	} else {
		if tx.RowsAffected > 0 {
			fmt.Printf("read结果：%+v\n", user4)
		} else {
			slog.Info("查无结果")
		}
	}

	// SELECT * FROM `user` USE INDEX (`id`,`idx_uid`) WHERE uid>0
	db.Where("uid>0").
		Clauses(hints.UseIndex("id", "idx_uid")). //给mysql一个建议的索引范围，这个范围之外的索引mysql就不再考虑了
		Find(&users)
	// SELECT * FROM `user` FORCE INDEX (`idx_uid`) WHERE uid>0
	db.Where("uid>0").
		Clauses(hints.ForceIndex("idx_uid")). //强制mysql使用某个索引
		Find(&users)
}

// 基于统计的查询
func ReadWithStatistics(db *gorm.DB) {
	type Result struct {
		City string
		Mid  float64
	}

	var results []Result
	db.Model(User{}).Select("city,avg(id) as mid").Group("city").Having("mid>0").Find(&results)
	fmt.Println("group by having查询结果：")
	for _, result := range results {
		fmt.Printf("%+v\n", result)
	}

	db.Table("user").Distinct("city").Find(&results)
	fmt.Println("distinct查询结果：")
	for _, result := range results {
		fmt.Printf("%+v\n", result)
	}

	var count int64
	db.Table("user").Where("city=?", "北京").Count(&count)
	fmt.Printf("count=%d\n", count)
}

/*
如果任何钩子回调返回错误，GORM将停止后续的操作并回滚事务。

查询时钩子的执行时机：
// 从 db 中加载数据
// Preloading (eager loading)
AfterFind
*/
