package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
)

// 数据校验。详细规则参见 https://goframe.org/docs/core/gvalid-rules
func Validation() {
	type Argument struct {
		Id        int    `v:"required|min:1"`
		Name      string `v:"required|length:4,10"`
		Password  string `v:"required|length:6,6"`
		Password2 string `v:"required|same:Password"`
		Email     string `v:"regex:^[\\w-\\.]+@([\\w-]+\\.)+[\\w-]{2,4}$"`
	}

	ctx := context.Background()
	arg1 := Argument{
		Id:        1,
		Name:      "orisun",
		Password:  "123456",
		Password2: "123456",
		Email:     "zbc.fds@sohu.com",
	}
	if err := g.Validator().Data(arg1).Run(ctx); err != nil {
		fmt.Println("arg1 invalid")
		for field, msgs := range err.Maps() {
			fmt.Println(field)
			for k, v := range msgs { //一个Field可能违反多项约束，所以这里是map
				fmt.Println(k, v)
			}
		}
		fmt.Println(strings.Repeat("-", 50))
	}
	arg2 := Argument{
		Id:        0,
		Name:      "zcy",
		Password:  "12345",
		Password2: "123456",
		Email:     "zbc.fds@sohucom",
	}
	if err := g.Validator().Data(arg2).Run(ctx); err != nil {
		fmt.Println("arg2 invalid")
		for field, msgs := range err.Maps() {
			fmt.Println(field)
			for k, v := range msgs { //一个Field可能违反多项约束，所以这里是map
				fmt.Println(k, v)
			}
		}
		fmt.Println(strings.Repeat("-", 50))
	}
}
