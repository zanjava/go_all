package model

import "time"

type News struct {
	Id           int        `gorm:"primaryKey" xorm:"pk autoincr"`
	UserId       int        //发布者id
	UserName     string     `gorm:"-" xorm:"-"` //数据库里没有这一列
	Title        string     //新闻标题
	Content      string     `gorm:"column:article" xorm:"article"`                 //正文
	PostTime     *time.Time `gorm:"column:create_time" xorm:"create_time created"` //发布时间
	DeleteTime   *time.Time `xorm:"deleted"`                                       //删除时间
	ViewPostTime string     `gorm:"-" xorm:"-"`
}

func (News) TableName() string {
	return "news"
}
