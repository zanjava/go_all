package gorm

import (
	"errors"
	"fmt"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

func HandleError(db *gorm.DB) {
	var user User
	db.First(&user)

	tx := db.Create(&user)
	if tx.Error != nil {
		//获得MySQL错误码，每一个code都对应一种特定的错误
		if mysqlErr, ok := tx.Error.(*mysql.MySQLError); ok { //接口的类型断言
			switch mysqlErr.Number { //针对不同的错误码，采取不同的处理方案
			case 1:
				//...
			case 2:
				//...
			default:
				fmt.Println("mysql error", "code", mysqlErr.Number, "msg", mysqlErr.Message)
			}
		} else {
			fmt.Printf("err %#v\n", tx.Error)
		}
	}

	// 或者使用errors.As，跟接口断言是等价的
	if tx.Error != nil {
		var mysqlErr *mysql.MySQLError //必须是指针，因为是指针实现了error接口
		if errors.As(tx.Error, &mysqlErr) {
			switch mysqlErr.Number { //针对不同的错误码，采取不同的处理方案
			case 1:
				//...
			case 2:
				//...
			default:
				fmt.Println("mysql error", "code", mysqlErr.Number, "msg", mysqlErr.Message)
			}
		} else {
			fmt.Printf("err %#v\n", tx.Error)
		}
	}
}
