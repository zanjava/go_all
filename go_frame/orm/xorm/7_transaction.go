package xorm

import (
	"fmt"
	"math/rand/v2"

	"xorm.io/xorm"
)

func Transaction(engine *xorm.Engine) error {
	session := engine.NewSession() //创建 Session 对象
	defer session.Close()

	//事务开始
	if err := session.Begin(); err != nil {
		return err
	}

	user := User{UserId: rand.IntN(100000), Degree: "本科", Gender: "男", City: "上海"}
	fmt.Printf("uid=%d\n", user.UserId)
	if _, err := session.Insert(&user); err != nil { //通过Session调用Insert时必须传指针
		session.Rollback() //手动回滚
		fmt.Println("第一次Insert回滚")
		return err
	}

	user.Id = 0
	if _, err := session.Insert(&user); err != nil { //第二次会失败，因为uid重复了
		session.Rollback() //手动回滚
		fmt.Println("第二次Insert回滚")
		return err
	}

	//提交事务
	// session.Commit()

	// user = User{UserId: rand.IntN(100000), Degree: "本科", Gender: "男", City: "上海"}
	// fmt.Printf("uid=%d\n", user.UserId)
	// if _, err := session.Insert(&user); err != nil { //通过Session调用Insert时必须传指针
	// 	session.Rollback() //手动回滚
	// 	fmt.Println("第三次Insert回滚")
	// 	return err
	// }

	//提交事务
	return session.Commit() // xorm的Commit可以执行多次，但gorm的不行，mysql原生也是支持多次commit的
}
