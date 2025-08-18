// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// News is the golang structure of table news for DAO operations like Where/Data.
type News struct {
	g.Meta     `orm:"table:news, do:true"`
	Id         interface{} // 新闻id
	UserId     interface{} // 发布者id
	Title      interface{} // 新闻标题
	Article    interface{} // 正文
	CreateTime *gtime.Time // 发布时间
	UpdateTime *gtime.Time // 最后修改时间
	DeleteTime *gtime.Time // 删除时间
}
