package basic

import "fmt"

func Inherit() {
	fmt.Println(PI)
	type User struct {
		Name string
		Age  int
	}
	type Vedio struct {
		Length int
		Name   string
		User   //匿名成员
	}
	u := User{Name: "张朝阳", Age: 18}
	v := Vedio{
		Length: 120,
		Name:   "go语言教程",
		User:   u, //注意：行尾一定要加逗号
	} //用变量类型来充当变量名称
	fmt.Println(v.Length)
	fmt.Println(v.Name)      //访问自己的Name
	fmt.Println(v.User.Name) //访问“父类”的Name
	fmt.Println(v.Age)       // Vedio从User里“继承”了Age
}
