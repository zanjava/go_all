package database

import (
	"errors"
	"fmt"
	"go/post/database/model"
	"log/slog"

	"github.com/go-sql-driver/mysql"
)

// 注册新用户。password是md5之后的密码
func RegistUser(name, password string) error {
	user := &model.User{Name: name, PassWord: password}
	_, err := PostDB.Insert(user)
	if err != nil {
		var mysqlErr *mysql.MySQLError //必须是指针，因为是指针实现了error接口
		if errors.As(err, &mysqlErr) {
			if mysqlErr.Number == 1062 { //违反uniq key
				return fmt.Errorf("用户名[%s]已存在", name)
			}
		}
		slog.Error("用户注册失败", "name", name, "error", err)
		return errors.New("用户注册失败，请稍后重试")
	}
	return nil
}

// 注销用户
func LogOffUser(uid int) error {
	affected, err := PostDB.ID(uid).Delete(model.User{})
	if err != nil {
		slog.Error("注销用户失败", "uid", uid, "error", err)
		return errors.New("用户注销失败，请稍后重试")
	}
	if affected == 0 {
		return fmt.Errorf("用户注销失败，uid %d不存在", uid)
	}
	return nil
}

func GetUserById(uid int) *model.User {
	user := &model.User{Id: uid}
	ok, err := PostDB.Get(user) //隐含的where条件是id
	if err != nil {             //发生异常
		slog.Error("GetUserById failed", "uid", uid, "error", err)
		return nil
	}
	if !ok { //查无结果
		return nil
	}
	return user
}

func GetUserByName(name string) *model.User {
	var user model.User
	ok, err := PostDB.Where("name=?", name).Get(&user)
	if err != nil { //发生异常
		slog.Error("GetUserByName failed", "name", name, "error", err)
		return nil
	}
	if !ok { //查无结果
		return nil
	}
	return &user
}

func UpdateUserName(uid int, name string) error {
	affected, err := PostDB.Where("id=?", uid).Update(model.User{Name: name})
	if err != nil {
		slog.Error("UpdateUserName failed", "uid", uid, "new name", name, "error", err)
		return errors.New("用户名修改失败，请稍后重试")
	}
	if affected <= 0 {
		return fmt.Errorf("用户id[%d]不存在", uid)
	}
	return nil
}

// newPass和oldPass都是md5之后的密码
func UpdatePassword(uid int, newPass, oldPass string) error {
	affected, err := PostDB.Where("id=? and password=?", uid, oldPass).Update(model.User{PassWord: newPass})
	if err != nil {
		slog.Error("UpdatePassword failed", "uid", uid, "error", err)
		return errors.New("用户名修改失败，请稍后重试")
	} else {
		if affected <= 0 {
			return errors.New("旧密码不对")
		} else {
			return nil
		}
	}
}
