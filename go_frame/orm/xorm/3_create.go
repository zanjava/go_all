package xorm

import (
	"fmt"
	"log/slog"
	"math/rand/v2"

	"xorm.io/xorm"
)

// 根据struct创建
func Create(engine *xorm.Engine) {
	//插入一条记录
	user := User{
		// Id:       10, //Id是自增的，可以不赋值，也可以显式地给Id赋值
		UserId:   rand.IntN(100000),
		Degree:   "本科",
		Gender:   "男",
		City:     "上海",
		Keywords: []string{"编程", "golang"},
		//会自动用当前时刻给created和updated字段赋值
		//在Insert时，version标记的字段将会被设置为1
	}
	affected, err := engine.Insert(&user) //如果Field Id没有xorm tag，则Insert会给user的Id赋值，否则不会。插入一条记录也可以使用InsertOne
	if err != nil {
		slog.Error("插入记录失败", "error", err)
	}
	fmt.Printf("after insert user is %#v\n", user)
	fmt.Printf("影响%d行\n", affected)

	//一次性插入多条
	user1 := user                    //发生拷贝
	user1.Id = 0                     //把主键置为0
	user1.UserId = rand.IntN(100000) //UserId上有唯一性约束
	user2 := user                    //发生拷贝
	user2.Id = 0                     //把主键置为0
	user2.UserId = rand.IntN(100000)
	users := []User{user1, user2}        //切片里的元素是或不是指针都行
	affected, err = engine.Insert(users) //一条SQL插入所有数据
	if err != nil {
		slog.Error("插入记录失败", "error", err)
	}
	fmt.Printf("影响%d行\n", affected)

	//API暂不支持批量插入。各个数据库对SQL语句有长度限制，因此这样的语句有一个最大的记录数，根据经验测算在150条左右。大于150条后，生成的sql语句将太长可能导致执行失败。因此在插入大量数据时，目前需要自行分割成每150条插入一次。
}
