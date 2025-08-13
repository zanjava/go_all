package gorm

import (
	"fmt"
	"math/rand/v2"

	"gorm.io/gorm"
)

// Save会保存所有的字段，即使字段是零值。主键为0时Save相当于Create
func Save(db *gorm.DB) {
	user := User{UserId: rand.IntN(100000), Degree: "本科", Gender: "男", City: "上海"}
	db.Save(&user) //主键为0值，Save相当于Create

	var user2 User
	db.Last(&user2)
	user2.Degree = "硕士"
	db.Save(&user2) //必须传指针
}

// Update指定需要更新的列
func Update(db *gorm.DB) {
	// 根据map更新
	tx := db.Model(&User{}). //必须传指针
					Where("city=?", "北京").Updates(
		map[string]any{"degree": "硕士", "gender": "男"},
	)
	fmt.Printf("更新了%d行\n", tx.RowsAffected)

	//根据结构体更新，只会更新非0值
	tx = db.Model(&User{}). //必须传指针
				Where("city=?", "北京").Updates(
		User{Degree: "本科", Gender: "男", Id: 1},
	)
	fmt.Printf("更新了%d行\n", tx.RowsAffected)
}

/*
如果任何钩子回调返回错误，GORM将停止后续的操作并回滚事务。

Update时钩子的执行时机：
// 开始事务
BeforeSave
BeforeUpdate
// 关联前的 save
// 更新 db
// 关联后的 save
AfterUpdate
AfterSave
// 提交或回滚事务
*/
