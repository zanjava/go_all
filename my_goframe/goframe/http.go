package main

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcfg"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/gvalid"
)

// 从GET请求的URL中获取参数
func url(r *ghttp.Request) {
	name := gconv.String(r.GetQuery("name"))
	addr := gconv.String(r.GetQuery("addr", "china")) //如果没传addr参数，则默认为China
	r.Response.WriteString(name + " live in " + addr) //返回 text/plain格式
}

// 从Restful风格的url中获取参数
func restful(r *ghttp.Request) {
	name := gconv.String(r.GetRouter("name"))
	addr := gconv.String(r.GetRouter("addr", "china"))
	r.Response.WriteString(name + " live in " + addr)
}

// 从postForm表单中获取参数
func postForm(r *ghttp.Request) {
	name := gconv.String(r.GetForm("name"))
	addr := gconv.String(r.GetForm("addr", "china"))
	r.Response.WriteXmlExit(g.Map{"name": name, "address": addr}) //返回xml格式
}

// 对参数的来源不作限制，可以是query、form、path、json body等
type Student struct {
	Name string `p:"name" v:"required|length:4,10"` //用p指定参数名称
	Addr string `p:"addr" v:"required|length:3,20"`
}

// 从post请求体中取出参数
func postJson(r *ghttp.Request) {
	var req Student
	if err := r.Parse(&req); err != nil { //参数绑定并校验
		if v, ok := err.(gvalid.Error); ok {
			r.Response.WriteString(strings.Join(v.Strings(), ";"))
		} else {
			r.Response.WriteString(err.Error())
		}
		return
	} else {
		// r.Response.WriteJson(req)   // 返回json格式
		// return
		r.Response.WriteJsonExit(req) // 带Exit会立即结束当前handler，就不用手动写return了
	}
}

// 上传单个文件
func upload_file(r *ghttp.Request) {
	file := r.GetUploadFile("file") //获取上传的文件
	if file == nil {
		r.Response.Write("empty file")
		return
	}
	name, err := file.Save("./data/") // 将上传的文件保存到data目录下
	if err != nil {
		r.Response.Write(err)
		return
	}
	r.Response.Write(name)
}

// 上传多个文件
func upload_multi_file(r *ghttp.Request) {
	files := r.GetUploadFiles("files")  //获取上传的多个文件
	names, err := files.Save("./data/") //把多个文件都保存到data目录下
	if err != nil {
		r.Response.WriteExit(err)
	} else {
		r.Response.WriteExit("成功上传" + strconv.Itoa(len(names)) + "个文件")
	}
}

func cookie(r *ghttp.Request) {
	token := gconv.String(r.Cookie.Get("token")) //从请求头里取出指定的Cookie

	//返回cookie
	name := "language"
	value := token
	maxAge := 7 * 24 * time.Hour //cookie的有效时间，时间单位秒。如果不设置过期时间，默认情况下关闭浏览器后cookie被删除
	path := "/"                  //cookie存放目录
	domain := "www.baidu.com"    //cookie从属的域名,不区分协议和端口。如果不指定domain则默认为本host(如b.a.com)，如果指定的domain是一级域名(如a.com)，则二级域名(b.a.com)下也可以访问
	r.Cookie.SetCookie(name, value, domain, path, maxAge)
}

// 返回html页面
func html(r *ghttp.Request) {
	r.Response.WriteTpl("learn/static/template.html", g.Map{"title": "用户信息", "name": "zcy", "addr": "bj"}) //模板渲染
}

// 接口计时中间件
func TimeMW(r *ghttp.Request) {
	begin := time.Now()
	r.Middleware.Next()
	g.Log().Infof(context.Background(), "interaface %s use time %d ms", r.RequestURI, time.Since(begin).Milliseconds())
}

func Web() {
	// 指定logger配置文件在哪儿
	g.Cfg().GetAdapter().(*gcfg.AdapterFile).SetFileName("learn/conf/learn.yaml")

	server := g.Server()
	server.SetAddr("127.0.0.1:5678") //指定监听端口

	// 静态资源
	server.AddStaticPath("/css", "learn/static/css")

	server.Use(TimeMW) //添加若干个全局中间件

	// 路由，GET和POST不限
	server.BindHandler("/student", url)
	server.BindHandler("/student/:name/*addr", restful)
	server.BindHandler("/student/form", postForm)
	server.BindHandler("/student/json", postJson)
	server.BindHandler("/upload", upload_file)
	server.BindHandler("/upload_files", upload_multi_file)
	server.BindHandler("/cookie", cookie)
	server.BindHandler("/user/html", html)
	server.BindHandler("/user/old_page", func(r *ghttp.Request) {
		r.Response.RedirectTo("/user/html", http.StatusMovedPermanently)
	})

	// 启动server
	server.Run()
}
