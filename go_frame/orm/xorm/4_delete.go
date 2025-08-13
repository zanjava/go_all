package xorm

import (
	"fmt"
	"log/slog"

	"xorm.io/xorm"
)

func Delete(engine *xorm.Engine) {
	//根据主键删除
	affected, err := engine.ID(10).Delete(User{})
	if err != nil {
		slog.Error("删除记录失败", "error", err)
	}
	fmt.Printf("删除%d行\n", affected)

	//根据普通Where条件删除
	affected, err = engine.Where("degree=?", "本科").Delete(User{})
	if err != nil {
		slog.Error("删除记录失败", "error", err)
	}
	fmt.Printf("删除%d行\n", affected)

	/**
	如果model中有一个Field是time.Time/int/int64类型，且被标记为`xorm:"deleted"`，则engine.Delete()是软删除，即只是给deleted列赋值了，并没有真正删除，背后对应的SQL语句是update而不是delete
	*/
}
