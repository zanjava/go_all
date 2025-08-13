package xorm

import (
	"fmt"

	"xorm.io/xorm"
)

// 执行原生的select语句
func RawSelect(engine *xorm.Engine) {
	results, err := engine.QueryInterface("select id,uid,city from user where id>? and uid>? limit 3", 2, 4)
	if err == nil {
		for _, result := range results {
			fmt.Printf("%+v\n", result) //result是一个map[string]any
		}
	}
	fmt.Println()
}

// 执行原生的update、insert、delete语句
func RawExec(engine *xorm.Engine) {
	result, err := engine.Exec("update user set degree=? where id=?", "大专", 24)
	if err == nil {
		affected, _ := result.RowsAffected()
		fmt.Printf("更新了%d行\n", affected)
	}
}
