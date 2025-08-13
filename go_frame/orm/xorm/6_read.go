package xorm

import (
	"fmt"
	"log/slog"

	"xorm.io/xorm"
)

func Read(engine *xorm.Engine) {
	/**
	Get返回一个查询结果
	*/
	user := User{City: "上海"} //非0值的Field将默认作为查询的where条件
	ok, err := engine.
		Select("*").
		Get(&user) // 虽然查询结果有多条，但是Get只取一条
	if err != nil {
		slog.Error("读DB失败", "error", err)
	} else {
		if !ok { // 查询返回false表示没有结果
			slog.Info("查无结果")
		} else {
			slog.Info("读库成功", "user", user)
		}
	}

	//启用PrepareStmt
	session := engine.Prepare()
	session.Select("*").Get(&user)

	var user2 User
	ok, err = engine.
		Cols("uid", "city", "gender", "keywords").
		// Select("uid,city,gender,keywords"). //跟上一行是等价的
		IndexHint("force", "", "idx_uid"). //强制mysql使用某个索引
		// IndexHint("use", "", "id,idx_uid"). //给mysql一个建议的索引范围，这个范围之外的索引mysql就不再考虑了
		Where("uid>100 and degree='大专'").
		Where("degree='大专'"). //多个where之间是And关系
		// Where("city in ?", []string{"北京", "上海"}).//错误，必须用In()函数
		In("city", []string{"北京", "上海"}).
		And("degree like ?", "%科").
		Or("gender=?", "女").
		OrderBy("id desc"). //排序
		Desc("uid").        //降序
		Asc("city").        //升序
		Limit(1, 3).        //LIMIT 1 OFFSET 3
		Get(&user2)         //注意：要传结构体指针；同时结构体的非0字段将作为where条件
	if err != nil {
		slog.Error("读DB失败", "error", err)
	} else {
		if !ok {
			slog.Info("查无结果")
		} else {
			slog.Info("读库成功", "user", user2)
		}
	}

	// 错误的写法
	var user3 *User //不同于var user3 User，还没申请内存空间。通过反射给user3赋值时发现还没给user3申请好内存空间
	_, err = engine.Get(user3)
	if err != nil {
		slog.Error("读DB失败", "error", err)
	}

	// 正确的写法
	var user4 = new(User)
	_, err = engine.Get(user4)
	if err != nil {
		slog.Error("读DB失败", "error", err)
	} else {
		slog.Info("读库成功", "user", user4)
	}

	/**
	Find返回多个查询结果
	*/
	var users []User
	err = engine.Limit(3).Find(&users) //要修改切片的长度，所以要传切片的指针
	if err != nil {
		slog.Error("读DB失败", "error", err)
	} else {
		fmt.Println("多个read结果")
		for _, u := range users {
			fmt.Printf("%+v\n", u)
		}
	}

	/**
	用Rows查询多个结果
	*/
	rows, err := engine.Limit(3).Rows(&User{}) //需要结构体传指针
	if err == nil {
		defer rows.Close()
		var u User
		fmt.Println("多个rows结果")
		for rows.Next() {
			rows.Scan(&u)
			fmt.Printf("user %+v\n", u)
		}
	} else {
		slog.Error("Rows failed", "error", err)
	}
}

// 基于统计的查询
func ReadWithStatistics(engine *xorm.Engine) {
	type Result struct {
		Degree string
		City   string
		Mid    float64
	}
	var results []Result

	// group by having
	err := engine.Table(User{}).Select("city, avg(id) as mid").GroupBy("city").Having("mid>0").Find(&results)
	if err == nil {
		fmt.Println("city  mid")
		for _, result := range results {
			fmt.Printf("%s %.2f\n", result.City, result.Mid)
		}
	}

	// distinct
	err = engine.Distinct("degree,city").Table(User{}).Find(&results)
	if err == nil {
		fmt.Println("degree  city")
		for _, result := range results {
			fmt.Printf("%s  %s\n", result.Degree, result.City)
		}
	}

	// count
	count, err := engine.Where("degree=?", "本科").Count(User{})
	if err == nil {
		fmt.Println("总数", count)
	}
}
