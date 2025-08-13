package xorm

import (
	"fmt"
	"time"

	"xorm.io/xorm"
)

// Update指定需要更新的列。
// 注意：deleted列有值的记录不会被更新
func Update(engine *xorm.Engine) {
	// 根据map更新
	affected, _ := engine.
		Table(User{}). //指定表名
		Where("city=?", "上海").Update(
		//会自动给updated字段赋值
		map[string]any{"degree": "硕士", "gender": "男", "keywords": ""},
	)
	fmt.Printf("更新了%d行\n", affected)

	//根据结构体更新，默认只会更新非0值，除非指定MustCols或AllCols
	affected, _ = engine.
		Where("city=?", "上海").
		// ID(10). // where id=10
		MustCols("keywords").
		// AllCols().
		Update(
			User{Degree: "本科", Gender: "男"}, //传结构体无需指定表名
		)
	fmt.Printf("更新了%d行\n", affected)
}

// 乐观锁更新。
// 在Insert时，version标记的字段将会被设置为1，在Update时，Update的内容必须包含version原来的值才能更新成功（防止并发更新）
func UpdateByVersion(engine *xorm.Engine) {
	var user User
	ok, err := engine.ID(6).Get(&user)
	if err != nil {
		fmt.Printf("Get失败:%s\n", err.Error())
		return
	}
	if !ok {
		fmt.Println("查无结果")
		return
	}
	fmt.Printf("更新前version=%d, id=%d\n", user.Version, user.Id)
	time.Sleep(50 * time.Millisecond) //模拟执行一些操作

	user.UserId += 1 //更新uid
	if affected, err := engine.ID(user.Id).Update(user); err == nil {
		if affected > 0 { //并发情况下，只有一个协程的 affected > 0
			fmt.Printf("更新成功,version=%d, 更新%d行\n", user.Version, affected)
		} else {
			fmt.Println("更新失败")
		}
	} else {
		fmt.Printf("更新失败, error=%s\n", err.Error())
	}
}
