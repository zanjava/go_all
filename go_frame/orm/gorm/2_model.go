package gorm

import (
	"time"
)

// gorm.Model的定义。
// 默认情况下，GORM使用ID作为主键，使用结构体名的蛇形复数作为表名，字段名的蛇形作为列名。
// 不建议使用Migrator(迁移)功能，表的维护由DBA负责，而不是开发人员。
type User struct {
	// ID         int       //名为ID的字段为默认主键
	Id        int       `gorm:"primaryKey;column:id"` //显式指定主键，显式指定表里对应的列名
	UserId    int       `gorm:"column:uid"`           //显式指定列名
	Degree    string    //驼峰转为蛇形就是对应的列名
	Keywords  []string  `gorm:"json"`               //转为json，可以对应DB里的char、varchar或text
	CreatedAt time.Time `gorm:"column:create_time"` //在Create时GORM会自动把当前时间赋给CreatedAt
	UpdatedAt time.Time `gorm:"column:update_time"` //在Update时GORM会自动把当前时间赋给UpdatedAt
	//通过`gorm:"type:date"`可以显式指定对应的不是datetime，而是date
	Gender   string
	City     string
	Province string `gorm:"-"` //表里没有这一列（比如表里没有这个字段，但结构体需要有它）
}

// 显式指定表名
// func (User) TableName() string {
// 	return "user"
// }
