package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10" //注意要用新版本v10
)

type User struct {
	Name  string `form:"name" binding:"required"`       //required:必须上传name参数。form可以绑定formdata和url问号后面的参数
	Score int    `form:"score" binding:"gt=0,required"` //score必须为正数

	Enrollment time.Time `form:"enrollment" binding:"required,before_today" time_format:"2006-01-02" time_utc:"8"`       //自定义验证before_today，日期格式东8区
	Graduation time.Time `form:"graduation" binding:"required,gtfield=Enrollment" time_format:"2006-01-02" time_utc:"8"` //毕业时间要晚于入学时间
}

// 自定义验证器
var beforeToday validator.Func = func(fl validator.FieldLevel) bool {
	if date, ok := fl.Field().Interface().(time.Time); ok { //通过反射获得结构体Field的值
		today := time.Now()
		if date.Before(today) {
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}

func processErr(err error) string {
	if err == nil {
		return ""
	}

	//ValidationErrors是一个错误切片，它保存了每个字段违反的每个约束信息
	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		msgs := make([]string, 0, 3)
		for _, validationErr := range validationErrs {
			msgs = append(msgs, fmt.Sprintf("字段 [%s] 不满足条件[%s]", validationErr.Field(), validationErr.Tag()))
		}
		return strings.Join(msgs, ";")
	} else {
		return "invalid error"
	}
}

func main5() {
	engine := gin.Default()

	//注册验证器
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("before_today", beforeToday)
	}

	engine.GET("/", func(ctx *gin.Context) {
		var user User
		if err := ctx.ShouldBind(&user); err != nil { //在绑定参数的同时，完成合法性校验
			msg := processErr(err)
			ctx.String(http.StatusBadRequest, "参数绑定失败："+msg) //校验不符合时，返回哪里不符合
		} else {
			ctx.JSON(http.StatusOK, user) //校验通过时，返回一个json
		}
	})

	engine.Run("127.0.0.1:5678")
}

/*
正确：
http://localhost:5678?name=zcy&score=1&enrollment=2025-02-20&graduation=2045-09-23

错误：
name缺失
http://localhost:5678?score=1&enrollment=2025-02-20&graduation=2045-09-23

score=0
http://localhost:5678?name=zcy&score=0&enrollment=2025-02-20&graduation=2045-09-23

enrollment晚于今天
http://localhost:5678?name=zcy&score=1&enrollment=2035-03-22&graduation=2045-09-23

graduation早于enrollment
http://localhost:5678?name=zcy&score=1&enrollment=2025-02-20&graduation=2015-09-23
*/
