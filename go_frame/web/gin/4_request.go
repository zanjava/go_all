package main

import (
	"encoding/json"
	"fmt"
	"go/frame/web/idl"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// 从GET请求的URL中获取参数
func url(engine *gin.Engine) {
	//http://127.0.0.1:5678/student?name=zcy&addr=bj
	engine.GET("/student", func(ctx *gin.Context) {
		a := ctx.Query("name")
		b := ctx.DefaultQuery("addr", "China")     //如果没传addr参数，则默认为China
		ctx.String(http.StatusOK, a+" live in "+b) //http response
	})
}

// 从Restful风格的url中获取参数
// : 对应一级目录  *可以对应多级目录
func restful(engine *gin.Engine) {
	engine.GET("/student/:name/*addr", func(ctx *gin.Context) {
		name := ctx.Param("name")
		addr := ctx.Param("addr") //会包含'/'
		ctx.String(http.StatusOK, name+" live in "+addr)
	})
}

// 从postForm表单中获取参数
func postForm(engine *gin.Engine) {
	engine.POST("/student/form", func(ctx *gin.Context) {
		name := ctx.PostForm("username")
		addr := ctx.DefaultPostForm("addr", "China") //如果没传addr参数，则默认为China
		ctx.String(http.StatusOK, name+" live in "+addr)
	})
}

// 从post请求体中取出参数
func postJson(engine *gin.Engine) {
	engine.POST("/student/json", func(ctx *gin.Context) {
		var stu Student
		bs, _ := io.ReadAll(ctx.Request.Body)
		if err := json.Unmarshal(bs, &stu); err == nil {
			name := stu.Name
			addr := stu.Addr
			ctx.String(http.StatusOK, name+" live in "+addr)
		}
	})
}

// 新版参数绑定
func ShouldBindBodyWith(ctx *gin.Context) error {
	var body []byte
	// 先尝试从ctx里获取request body
	if v, exists := ctx.Get("rbody"); exists {
		if bs, ok := v.([]byte); ok {
			body = bs
		}
	}
	if body == nil {
		// 从流里读取request body
		body, _ = io.ReadAll(ctx.Request.Body)
		// 把request body放入ctx
		ctx.Set("rbody", body)
	}
	// 参数解析、绑定
	var stu Student
	return json.Unmarshal(body, &stu)
}

// 业务Handler
func Handler(ctx *gin.Context) {
	var stu Student
	ctx.BindJSON(&stu) //获取不到参数
	ctx.String(200, stu.Name+stu.Addr)
}

// 上传单个文件
func upload_file(engine *gin.Engine) {
	//限制表单上传大小为8M，默认上限是32M
	engine.MaxMultipartMemory = 8 << 20
	engine.POST("/upload", func(ctx *gin.Context) {
		file, err := ctx.FormFile("file")
		if err != nil {
			fmt.Printf("get file error %v\n", err)
			ctx.String(http.StatusInternalServerError, "upload file failed")
		} else {
			if err = ctx.SaveUploadedFile(file, "../data/"+file.Filename); err == nil { //把用户上传的文件存到data目录下
				ctx.String(http.StatusOK, file.Filename)
			} else {
				fmt.Printf("save file to %s failed: %v\n", "../data/"+file.Filename, err)
			}
		}
	})
}

// 上传多个文件
func upload_multi_file(engine *gin.Engine) {
	engine.POST("/upload_files", func(ctx *gin.Context) {
		form, err := ctx.MultipartForm() //MultipartForm中包含多个文件
		if err != nil {
			ctx.String(http.StatusBadRequest, err.Error())
		} else {
			//从MultipartForm中获取上传的文件
			files := form.File["files"]
			for _, file := range files {
				ctx.SaveUploadedFile(file, "../data/"+file.Filename) //把用户上传的文件存到data目录下
			}
			ctx.String(http.StatusOK, "upload "+strconv.Itoa(len(files))+" files")
		}
	})
}

type Student struct {
	Name     string   `form:"username" json:"name" uri:"user" xml:"user" yaml:"user" binding:"required"` // required:必须上传name参数。form可以绑定formdata和url问号后面的参数
	Addr     string   `form:"addr" json:"addr" uri:"addr" xml:"addr" yaml:"addr" binding:"required"`
	Keywords []string `form:"keywords"` //可以绑定html的checkbox（复选框）
}

func formBind(engine *gin.Engine) {
	engine.POST("/stu/form", func(ctx *gin.Context) {
		var stu Student
		//跟ShouldBind对应的是MustBind。MustBind内部会调用ShouldBind，如果ShouldBind发生error会直接c.AbortWithError(http.StatusBadRequest, err)
		if err := ctx.ShouldBind(&stu); err != nil {
			fmt.Println(err)
			ctx.String(http.StatusBadRequest, "parse paramter failed")
		} else {
			ctx.String(http.StatusOK, stu.Name+" live in "+stu.Addr)
		}
	})
}

func jsonBind(engine *gin.Engine) {
	engine.POST("/stu/json", func(ctx *gin.Context) {
		var stu Student
		if err := ctx.ShouldBindJSON(&stu); err != nil {
			fmt.Println(err)
			ctx.String(http.StatusBadRequest, "parse paramter failed")
		} else {
			ctx.String(http.StatusOK, stu.Name+" live in "+stu.Addr)
		}
	})
}

func uriBind(engine *gin.Engine) {
	//GET请求的参数在uri里
	engine.GET("/stu/uri/:user/*addr", func(ctx *gin.Context) {
		// fmt.Println(ctx.Request.URL)
		var stu Student
		if err := ctx.ShouldBindUri(&stu); err != nil {
			fmt.Println(err)
			ctx.String(http.StatusBadRequest, "parse paramter failed")
		} else {
			ctx.String(http.StatusOK, stu.Name+" live in "+stu.Addr)
		}
	})
}

func xmlBind(engine *gin.Engine) {
	engine.POST("/stu/xml", func(ctx *gin.Context) {
		var stu Student
		if err := ctx.ShouldBindXML(&stu); err != nil {
			fmt.Println(err)
			ctx.String(http.StatusBadRequest, "parse paramter failed")
		} else {
			ctx.String(http.StatusOK, stu.Name+" live in "+stu.Addr)
		}
	})
}

func yamlBind(engine *gin.Engine) {
	engine.POST("/stu/yaml", func(ctx *gin.Context) {
		var stu Student
		if err := ctx.ShouldBindYAML(&stu); err != nil {
			fmt.Println(err)
			ctx.String(http.StatusBadRequest, "parse paramter failed")
		} else {
			ctx.String(http.StatusOK, stu.Name+" live in "+stu.Addr)
		}
	})
}

// 请求体是个流，只能读一次，第二次就读不到了。所以普通的Bind只能绑定一次。
// ctx.ShouldBindBodyWith 会在绑定之前将 body 存储到上下文中, 这会对性能造成轻微影响。
func multiBind(engine *gin.Engine) {
	engine.POST("/stu/multi_type", func(ctx *gin.Context) {
		var stu Student
		var stu2 idl.Student
		if err := ctx.ShouldBindBodyWith(&stu, binding.JSON); err == nil {
			ctx.String(http.StatusOK, stu.Name+" live in "+stu.Addr)
		} else if err := ctx.ShouldBindBodyWith(&stu, binding.XML); err == nil {
			ctx.String(http.StatusOK, stu.Name+" live in "+stu.Addr)
		} else if err := ctx.ShouldBindBodyWith(&stu, binding.YAML); err == nil {
			ctx.String(http.StatusOK, stu.Name+" live in "+stu.Addr)
		} else if err := ctx.ShouldBindBodyWith(&stu2, binding.ProtoBuf); err == nil {
			ctx.String(http.StatusOK, stu2.Name+" live in "+stu2.Address)
		} else {
			ctx.String(http.StatusBadRequest, "不支持的参数类型")
		}

	})
}

func main4() {
	engine := gin.Default()

	// 静态资源
	engine.Static("/css", "web/static/css")                     //在url是访问目录/js相当于访问文件系统中的views/js目录
	engine.StaticFile("/favicon.ico", "web/static/img/dqq.png") //在url中访问文件/favicon.ico，相当于访问文件系统中的views/img/dqq.png文件

	url(engine)
	restful(engine)
	postForm(engine)
	postJson(engine)
	upload_file(engine)       //用postman模拟一个post请求，注意body类型选择form-data，Key名称为file，类型为File，在Value里选择本地文件
	upload_multi_file(engine) //用postman模拟一个post请求，注意body类型选择form-data，Key名称为files，类型为File，在Value里选择多个本地文件

	formBind(engine)
	jsonBind(engine)
	uriBind(engine)
	xmlBind(engine)
	yamlBind(engine)

	multiBind(engine)

	engine.Run("127.0.0.1:5678")
}

// go run ./web/gin
// go build -tags=jsoniter -o gin.exe ./web/gin
