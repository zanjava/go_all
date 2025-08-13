package xorm

import (
	"errors"
	"fmt"

	"github.com/go-sql-driver/mysql"
	"xorm.io/xorm"
)

func HandleError(engine *xorm.Engine) {
	var user User
	engine.Get(&user)

	_, err := engine.Insert(user)
	if err != nil {
		//获得MySQL错误码，每一个code都对应一种特定的错误
		if mysqlErr, ok := err.(*mysql.MySQLError); ok { //接口的类型断言
			switch mysqlErr.Number { //针对不同的错误码，采取不同的处理方案
			case 1:
				//...
			case 2:
				//...
			default:
				fmt.Println("mysql error", "code", mysqlErr.Number, "msg", mysqlErr.Message)
			}
		} else {
			fmt.Printf("err %#v\n", err)
		}
	}

	// 或者使用errors.As，跟接口断言是等价的
	if err != nil {
		var mysqlErr *mysql.MySQLError //必须是指针，因为是指针实现了error接口
		if errors.As(err, &mysqlErr) {
			switch mysqlErr.Number { //针对不同的错误码，采取不同的处理方案
			case 1:
				//...
			case 2:
				//...
			default:
				fmt.Println("mysql error", "code", mysqlErr.Number, "msg", mysqlErr.Message)
			}
		} else {
			fmt.Printf("err %#v\n", err)
		}
	}
}
