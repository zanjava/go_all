package database

import (
	"errors"
	"fmt"
	"go/post/database/model"
	"log/slog"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

// 注册新用户。password是md5之后的密码
func RegistUser(name, password string) error {
	user := new(model.User)
	user.Name = name
	user.PassWord = password
	err := PostDB.Create(user).Error
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
	user := model.User{Id: uid}
	tx := PostDB.
		// Session(&gorm.Session{DryRun: true}). // 测试阶段，先不要真删
		Delete(user)
	if tx.Error != nil {
		slog.Error("注销用户失败", "uid", uid, "error", tx.Error)
		return errors.New("用户注销失败，请稍后重试")
	}
	if tx.RowsAffected == 0 {
		return fmt.Errorf("用户注销失败，uid %d不存在", uid)
	}
	return nil
}

func GetUserById(uid int) *model.User {
	user := &model.User{Id: uid}
	tx := PostDB.Select("*").First(user) //隐含的where条件是id。注意：Find不会返回ErrRecordNotFound
	if tx.Error != nil {
		if !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			slog.Error("GetUserById failed", "uid", uid, "error", tx.Error)
		}
		return nil
	}
	return user
}

func GetUserByName(name string) *model.User {
	user := &model.User{}
	tx := PostDB.Select("*").Where("name=?", name).First(user) //注意：Find不会返回ErrRecordNotFound
	if tx.Error != nil {
		if !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			slog.Error("GetUserByName failed", "name", name, "error", tx.Error)
		}
		return nil
	}
	return user
}

func UpdateUserName(uid int, name string) error {
	tx := PostDB.Model(&model.User{}).Where("id=?", uid).Update("name", name)
	if tx.Error != nil {
		slog.Error("UpdateUserName failed", "uid", uid, "new name", name, "error", tx.Error)
		return errors.New("用户名修改失败，请稍后重试")
	} else {
		if tx.RowsAffected <= 0 {
			return fmt.Errorf("用户id[%d]不存在", uid)
		} else {
			return nil
		}
	}
}

// newPass和oldPass都是md5之后的密码
func UpdatePassword(uid int, newPass, oldPass string) error {
	tx := PostDB.Model(&model.User{}).Where("id=? and password=?", uid, oldPass).Update("password", newPass)
	if tx.Error != nil {
		slog.Error("UpdatePassword failed", "uid", uid, "error", tx.Error)
		return errors.New("密码修改失败，请稍后重试")
	} else {
		if tx.RowsAffected <= 0 {
			return errors.New("旧密码不对")
		} else {
			return nil
		}
	}
}
