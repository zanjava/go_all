package main

import (
	"go/frame/web/idl"
	"io"
	"net/http"
	"os"

	"github.com/bytedance/sonic"
	"github.com/gin-gonic/gin"
)

func text(engine *gin.Engine) {
	engine.GET("/user/text", func(c *gin.Context) {
		c.String(http.StatusOK, "hi boy") //response Content-Type:text/plain
	})
}

func json0(engine *gin.Engine) {
	engine.GET("/user/json0", func(c *gin.Context) {
		var stu struct { //匿名结构体
			Name    string `json:"name"`
			Address string `json:"addr"`
		}
		stu.Name = "zcy"
		stu.Address = "bj"
		s, _ := sonic.MarshalString(stu)
		c.Request.Header.Add("Content-Type", "application/json")
		c.String(http.StatusOK, s)
	})
}

func json1(engine *gin.Engine) {
	engine.GET("/user/json1", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"name": "zcy", "addr": "bj"}) //response Content-Type:application/json
		//log.Println("json1 response:", gin.H{"name": "zcy", "addr": "bj"})
	})
}

func json2(engine *gin.Engine) {
	var stu struct { //匿名结构体
		Name    string `json:"name"`
		Address string `json:"addr"`
	}
	stu.Name = "zcy"
	stu.Address = "bj"
	engine.GET("/user/json2", func(c *gin.Context) {
		c.JSON(http.StatusOK, stu) //response Content-Type:application/json
	})
}

// 使用 JSONP 可以向不同域的服务器请求数据
func jsonp(engine *gin.Engine) {
	var stu struct { //匿名结构体
		Name    string `json:"name"`
		Address string `json:"addr"`
	}
	stu.Name = "zcy"
	stu.Address = "bj"
	engine.GET("/user/jsonp", func(ctx *gin.Context) {
		//如果请求参数里有callback=xxx，则response Content-Type为application/javascript，否则response Content-Type为application/json
		ctx.JSONP(http.StatusOK, stu)
	})
}

func xml(engine *gin.Engine) {
	type Stu struct {
		Name    string `xml:"name"`
		Address string `xml:"addr"`
	}
	var stu Stu // xml不支持匿名结构体
	stu.Name = "zcy"
	stu.Address = "bj"
	engine.GET("/user/xml", func(c *gin.Context) {
		c.XML(http.StatusOK, stu) //response Content-Type:application/xml
	})
}

func yaml(engine *gin.Engine) {
	var stu struct {
		Name string
		Addr string
	}
	stu.Name = "zcy"
	stu.Addr = "bj"
	engine.GET("/user/yaml", func(c *gin.Context) {
		c.YAML(http.StatusOK, stu) //response Content-Type:application/yaml
	})
}

func protoBuf(engine *gin.Engine) {
	stu := &idl.Student{
		Name:    "zcy",
		Address: "bj",
	}
	engine.GET("/user/pb", func(ctx *gin.Context) {
		ctx.ProtoBuf(http.StatusOK, stu)
	})
}

func html(engine *gin.Engine) {
	// engine.LoadHTMLGlob("web/static/")
	engine.LoadHTMLFiles("..\\static\\template.html", "..\\static\\student.html") //后面可以在go代码直接用template.html和student.html
	engine.GET("/user/html", func(c *gin.Context) {
		//使用go语言标准库的html template。通过json往前端页面上传值
		c.HTML(http.StatusOK, "template.html", gin.H{"title": "用户信息", "name": "zcy", "addr": "bj"})
	})
}

func redirect(engine *gin.Engine) {
	engine.GET("/user/old_page", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/user/html")
	})
}

func main() {
	// gin.SetMode(gin.ReleaseMode)   //GIN线上发布模式
	fout, _ := os.OpenFile("../data/log/gin.log", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	gin.DefaultWriter = io.MultiWriter(os.Stdout, fout) //GIN日志输出到多个地方
	// gin.DefaultWriter = io.Discard //禁用GIN日志

	engine := gin.Default()

	// 修改静态资源不需要重启GIN，刷新页面即可
	// 静态资源
	// http://127.0.0.1:5678/css/dqq.css
	engine.Static("/css", "../static/css")                     //在url里访问目录/css相当于访问文件系统中的web/static/css目录
	engine.StaticFile("/favicon.ico", "../static/img/dqq.png") //在url中访问文件/favicon.ico，相当于访问文件系统中的views/img/dqq.png文件

	text(engine)
	json0(engine)
	json1(engine)
	json2(engine)
	jsonp(engine)
	xml(engine)
	yaml(engine)
	protoBuf(engine)
	html(engine)
	redirect(engine)
	engine.Run("127.0.0.1:5678")
}

// go run ./web/gin
