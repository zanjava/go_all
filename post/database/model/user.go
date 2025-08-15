package model

type User struct {
	Id       int    `gorm:"primaryKey" xorm:"pk autoincr"`   //用户id
	Name     string `gorm:"column:name" xorm:"name"`         //用户名
	PassWord string `gorm:"column:password" xorm:"password"` //md5之后的密码
}

func (User) TableName() string {
	return "user"
}
