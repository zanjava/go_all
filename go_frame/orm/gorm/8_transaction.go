package gorm

import (
	"fmt"
	"math/rand/v2"

	"gorm.io/gorm"
)

func Transaction(db *gorm.DB) error {
	// 为了确保数据一致性，GORM 会在事务里执行写入操作（创建、更新、删除）。如果没有这方面的要求，您可以在初始化时禁用它，这将获得大约 30%+ 性能提升。
	db = db.Session(&gorm.Session{SkipDefaultTransaction: true})

	tx := db.Begin() // 开始事务

	defer func() {
		if err := recover(); err != nil {
			tx.Rollback() //手动回滚
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	user := User{UserId: rand.IntN(100000), Degree: "本科", Gender: "男", City: "上海"}
	fmt.Printf("uid=%d\n", user.UserId)
	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback() //手动回滚
		fmt.Println("第一次Create回滚")
		return err
	}

	user.Id = 0
	if err := tx.Create(&user).Error; err != nil { //第二次会失败，因为uid重复了
		tx.Rollback() //手动回滚
		fmt.Println("第二次Create回滚")
		return err
	}

	return tx.Commit().Error //提交事务。Commit和Rollback只能执行一个，且只能执行一次

	// tx.Commit()

	// user = User{UserId: rand.IntN(100000), Degree: "本科", Gender: "男", City: "上海"}
	// fmt.Printf("uid=%d\n", user.UserId)
	// if err := tx.Create(&user).Error; err != nil {
	// 	tx.Rollback()    //手动回滚
	// 	fmt.Println(err) //sql: transaction has already been committed or rolled back    Commit和Rollback只能执行一个，且只能执行一次
	// 	fmt.Println("第三次Create回滚")
	// 	return err
	// }
	// return tx.Commit().Error
}
