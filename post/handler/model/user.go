package model

type User struct {
	Name     string `form:"name" binding:"required,gte=2" validate:"required,gte=2"`   //长度>=2
	PassWord string `form:"pass" binding:"required,len=32" validate:"required,len=32"` //长度=32, pass是md5之后的密码
}

type ModifyPassRequest struct {
	//Uid     int    `form:"uid" binding:"required,gte=1"`
	OldPass string `form:"old_pass" binding:"required,len=32" validate:"required,len=32"`
	NewPass string `form:"new_pass" binding:"required,len=32" validate:"required,len=32"`
}
