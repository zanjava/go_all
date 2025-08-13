package gorm

import (
	"fmt"

	"gorm.io/gorm"
)

// 执行原生的select语句
func RawSelect(db *gorm.DB) {
	var users []User
	db.Raw("select id,uid,city from user where id>? and uid>? limit 3", 2, 4).Scan(&users)
	for _, user := range users {
		fmt.Printf("uid %d city %s\n", user.UserId, user.City)
	}

	fmt.Println()
	rows, err := db.Raw("select id,uid,city from user where id>? and uid>? limit 3", 2, 4).Rows()
	if err == nil {
		defer rows.Close()
		var id, uid int
		var city string
		for rows.Next() {
			rows.Scan(&id, &uid, &city)
			fmt.Printf("uid %d city %s\n", uid, city)
		}
	}
}

// 执行原生的update、insert、delete语句
func RawExec(db *gorm.DB) {
	tx := db.Exec("update user set degree=? where id=?", "大专", 30)
	fmt.Printf("更新了%d行\n", tx.RowsAffected)
}
