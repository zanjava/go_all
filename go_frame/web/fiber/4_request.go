package main

import (
	"encoding/json"
	"fmt"
	"go/frame/web/idl"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"google.golang.org/protobuf/proto"
)

func uri(app *fiber.App) {
	//http://127.0.0.1:5678/student?name=zcy&addr=bj&age=18
	app.Get("/student", func(ctx fiber.Ctx) error {
		a := ctx.Query("name")          //如果参数不存在将返回0值
		b := ctx.Query("addr", "China") //如果没传addr参数，则默认为China
		c := ctx.Query("age")
		return ctx.SendString(a + " live in " + b + " age is " + c)
	})

	app.Get("/student2", func(ctx fiber.Ctx) error {
		args := ctx.Queries() // ctx.Queries()直接把所有参数存储到一个map里
		return ctx.SendString(args["name"] + " live in " + args["addr"] + " age is " + args["age"])
	})
}

func restful(app *fiber.App) {
	app.Get("/student/:name/+/*", func(ctx fiber.Ctx) error {
		name, _ := url.QueryUnescape(ctx.Params("name")) //url反转义
		addr, _ := url.QueryUnescape(ctx.Params("+"))    //必选参数
		age, _ := url.QueryUnescape(ctx.Params("*"))     //可选参数
		return ctx.SendString(name + " live in " + addr + " age is " + age)
	})

	app.Get("/school/:school", func(ctx fiber.Ctx) error {
		a, _ := url.QueryUnescape(ctx.Params("school", "THU")) //可选参数
		return ctx.SendString("school is " + a)
	})

	app.Get("/location/:longitude.:latitude", func(ctx fiber.Ctx) error {
		longitude, _ := url.QueryUnescape(ctx.Params("longitude"))
		latitude, _ := url.QueryUnescape(ctx.Params("latitude"))
		return ctx.SendString("longitude is " + longitude + " latitude is " + latitude)
	})

	app.Get("/size/:min-:max", func(ctx fiber.Ctx) error {
		min, _ := url.QueryUnescape(ctx.Params("min"))
		max, _ := url.QueryUnescape(ctx.Params("max"))
		return ctx.SendString("size is " + min + " to " + max)
	})
}

func postForm(app *fiber.App) {
	app.Post("/student/form", func(ctx fiber.Ctx) error {
		name := ctx.FormValue("username")
		addr := ctx.FormValue("addr", "China") //如果没传addr参数，则默认为China
		return ctx.SendString(name + " live in " + addr)
	})
}

// 内容类型        						标记
//
// application/x-www-form-urlencoded	form
//
// multipart/form-data					form
//
// application/json						json
//
// application/xml						xml
//
// text/xml								xml
type Student struct {
	Name    string `form:"username" uri:"user" json:"name" xml:"user" validate:"required"` // required:必须上传name参数。form可以绑定formdata和url问号后面的参数
	Address string `form:"addr" uri:"addr" json:"addr" xml:"addr" validate:"required"`
}

func postJson(app *fiber.App) {
	app.Post("/student/json", func(ctx fiber.Ctx) error {
		var stu Student
		if err := json.Unmarshal(ctx.Request().Body(), &stu); err == nil {
			name := stu.Name
			addr := stu.Address
			return ctx.SendString(name + " live in " + addr)
		} else {
			return err
		}
	})
}

func postProtobuf(app *fiber.App) {
	app.Post("/student/pb", func(ctx fiber.Ctx) error {
		var stu idl.Student
		if err := proto.Unmarshal(ctx.Request().Body(), &stu); err == nil {
			name := stu.Name
			addr := stu.Address
			return ctx.SendString(name + " live in " + addr)
		} else {
			return err
		}
	})
}

// 上传单个文件
func upload_file(app *fiber.App) {
	app.Post("/upload", func(ctx fiber.Ctx) error {
		file, err := ctx.FormFile("file") //获得第一个文件
		if err != nil {
			fmt.Printf("get file error %v\n", err)
			return fiber.NewError(http.StatusInternalServerError, "upload file failed")
		}
		//把用户上传的文件存到data目录下
		if err = ctx.SaveFile(file, fmt.Sprintf("./data/%s", file.Filename)); err != nil {
			return fiber.NewError(http.StatusInternalServerError, fmt.Sprintf("save file to %s failed: %v\n", "./data/"+file.Filename, err))
		}
		return ctx.SendString(file.Filename)
	})
}

// 上传多个文件
func upload_multi_file(app *fiber.App) {
	app.Post("/upload_files", func(ctx fiber.Ctx) error {
		form, err := ctx.MultipartForm()
		if err != nil {
			return fiber.NewError(http.StatusBadRequest, err.Error())
		}
		//从MultipartForm中获取上传的文件
		files := form.File["files"]
		for _, file := range files {
			ctx.SaveFile(file, "./data/"+file.Filename) //把用户上传的文件存到data目录下
		}
		return ctx.SendString("upload " + strconv.Itoa(len(files)) + " files")
	})
}

func formBind(app *fiber.App) {
	app.Post("/stu/form", func(ctx fiber.Ctx) error {
		var stu Student
		if err := ctx.Bind().Form(&stu); err != nil {
			fmt.Println(err)
			return fiber.NewError(http.StatusBadRequest, "parse paramter failed")
		} else {
			return ctx.SendString(stu.Name + " live in " + stu.Address)
		}
	})
}

func jsonBind(app *fiber.App) {
	app.Post("/stu/json", func(ctx fiber.Ctx) error {
		var stu Student
		if err := ctx.Bind().JSON(&stu); err != nil {
			fmt.Println(err)
			return fiber.NewError(http.StatusBadRequest, "parse paramter failed")
		} else {
			return ctx.SendString(stu.Name + " live in " + stu.Address)
		}
	})
}

func uriBind(app *fiber.App) {
	app.Get("/stu/uri/:user/:addr", func(ctx fiber.Ctx) error {
		var stu Student
		if err := ctx.Bind().URI(&stu); err != nil {
			fmt.Println(err)
			return fiber.NewError(http.StatusBadRequest, "parse paramter failed")
		} else {
			name, _ := url.QueryUnescape(stu.Name)
			addr, _ := url.QueryUnescape(stu.Address)
			return ctx.SendString(name + " live in " + addr)
		}
	})
}

func xmlBind(app *fiber.App) {
	app.Post("/stu/xml", func(ctx fiber.Ctx) error {
		var stu Student
		if err := ctx.Bind().XML(&stu); err != nil {
			fmt.Println(err)
			return fiber.NewError(http.StatusBadRequest, "parse paramter failed")
		} else {
			return ctx.SendString(stu.Name + " live in " + stu.Address)
		}
	})
}

func multiBind(app *fiber.App) {
	app.Post("/stu/multi_type", func(ctx fiber.Ctx) error {
		var stu Student
		if err := ctx.Bind().JSON(&stu); err == nil {
			return ctx.SendString(stu.Name + " live in " + stu.Address)
		} else if err := ctx.Bind().XML(&stu); err == nil { //可以多次执行Bind()，数据一直能取到
			return ctx.SendString(stu.Name + " live in " + stu.Address)
		} else if err := ctx.Bind().Form(&stu); err == nil {
			return ctx.SendString(stu.Name + " live in " + stu.Address)
		} else {
			fmt.Println(err)
			return fiber.NewError(http.StatusBadRequest, "parse paramter failed")
		}
	})
}

func main4() {
	app := fiber.New()
	app.Use(logger.New())

	uri(app)
	restful(app)
	postForm(app)
	postJson(app)
	postProtobuf(app)
	upload_file(app)       //用postman模拟一个post请求，注意body类型选择form-data，Key名称为file，类型为File，在Value里选择本地文件
	upload_multi_file(app) //用postman模拟一个post请求，注意body类型选择form-data，Key名称为files，类型为File，在Value里选择多个本地文件

	formBind(app)
	jsonBind(app)
	uriBind(app)
	xmlBind(app)
	multiBind(app)

	if err := app.Listen("127.0.0.1:5678"); err != nil {
		slog.Error("fiber app start failed", "error", err)
	}
}
