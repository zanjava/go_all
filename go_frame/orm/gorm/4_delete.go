package gorm

import (
	"fmt"

	"gorm.io/gorm"
)

func Delete(db *gorm.DB) {
	tx := db.Where("degree=?", "专科").Delete(User{})
	fmt.Printf("删除%d行\n", tx.RowsAffected)

	var user User = User{Id: 10}
	db.Delete(user) //暗含的Where条件是id=10

	db.Delete(User{}, 1)              //暗含的Where条件是id=1
	db.Delete(User{}, []int{1, 2, 3}) //暗含的Where条件是id IN (1,2,3)
}

/*
如果任何钩子回调返回错误，GORM将停止后续的操作并回滚事务。

Delete时钩子的执行时机：
// 开始事务
BeforeDelete
// 删除 db 中的数据
AfterDelete
// 提交或回滚事务
*/
