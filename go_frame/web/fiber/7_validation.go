package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
)

type User struct {
	Name       string    `form:"name" validate:"required"`                          //required:必须上传name参数。form可以绑定formdata和url问号后面的参数
	Score      int       `form:"score" validate:"gt=0,required"`                    //score必须为正数
	Enrollment time.Time `form:"enrollment" validate:"required,before_today"`       //自定义验证before_today。不支持指定时间格式
	Graduation time.Time `form:"graduation" validate:"required,gtfield=Enrollment"` //毕业时间要晚于入学时间
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
		return fmt.Sprintf("invalid error type: %#v", err)
	}
}

type structValidator struct {
	validate *validator.Validate
}

// Validator needs to implement the Validate method
func (v *structValidator) Validate(out any) error {
	return v.validate.Struct(out)
}

func main7() {
	//注册验证器
	myValidator := validator.New()
	myValidator.RegisterValidation("before_today", beforeToday)
	//在fiber里使用自定义的验证器
	app := fiber.New(
		fiber.Config{
			StructValidator: &structValidator{validate: myValidator},
		},
	)

	// 路由
	app.Get("/", func(ctx fiber.Ctx) error {
		var user User
		if err := ctx.Bind().Query(&user); err != nil { //在绑定参数的同时，完成合法性校验
			msg := processErr(err)
			return fiber.NewError(http.StatusBadRequest, "参数绑定失败："+msg) //校验不符合时，返回哪里不符合
		} else {
			return ctx.JSON(user) //校验通过时，返回一个json
		}
	})

	if err := app.Listen("127.0.0.1:5678"); err != nil {
		slog.Error("fiber app start failed", "error", err)
	}
}
