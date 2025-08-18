// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// News is the golang structure for table news.
type News struct {
	Id         int         `json:"id"         orm:"id"          description:"新闻id"`   // 新闻id
	UserId     int         `json:"userId"     orm:"user_id"     description:"发布者id"`  // 发布者id
	Title      string      `json:"title"      orm:"title"       description:"新闻标题"`   // 新闻标题
	Article    string      `json:"article"    orm:"article"     description:"正文"`     // 正文
	CreateTime *gtime.Time `json:"createTime" orm:"create_time" description:"发布时间"`   // 发布时间
	UpdateTime *gtime.Time `json:"updateTime" orm:"update_time" description:"最后修改时间"` // 最后修改时间
	DeleteTime *gtime.Time `json:"deleteTime" orm:"delete_time" description:"删除时间"`   // 删除时间
}
