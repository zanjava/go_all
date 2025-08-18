// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// User is the golang structure of table user for DAO operations like Where/Data.
type User struct {
	g.Meta     `orm:"table:user, do:true"`
	Id         interface{} // 用户id，自增
	Name       interface{} // 用户名
	Password   interface{} // 密码的md5
	CreateTime *gtime.Time // 用户注册时间
	UpdateTime *gtime.Time // 最后修改时间
}
