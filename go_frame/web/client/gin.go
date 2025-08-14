package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"go/frame/web/idl"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"google.golang.org/protobuf/proto"
	"gopkg.in/yaml.v2"
)

type Student struct {
	Name    string `form:"username" uri:"user" json:"name" xml:"user"`
	Address string `form:"addr" uri:"addr" json:"addr" xml:"addr"`
}

type User struct {
	Name       string `form:"name"`
	Score      int    `form:"score"`
	Enrollment string `form:"enrollment"`
	Graduation string `form:"graduation"`
}

func processResponse(resp *http.Response) {
	defer resp.Body.Close()
	fmt.Println("响应头：")
	for k, v := range resp.Header {
		fmt.Printf("%s=%s\n", k, v[0])
	}
	fmt.Print("响应体：")
	io.Copy(os.Stdout, resp.Body)
	os.Stdout.WriteString("\n")
	if resp.StatusCode != http.StatusOK {
		slog.Error("异常状态码", "http response code", resp.StatusCode)
	}
	os.Stdout.WriteString("\n")
}

func Head(path string) {
	fmt.Println("HEAD " + path)
	resp, err := http.Head("http://127.0.0.1:5678" + path)
	if err != nil {
		slog.Error("http get failed", "error", err)
		return
	}
	processResponse(resp)
}

func Get(path string) {
	fmt.Println("GET " + path)
	resp, err := http.Get("http://127.0.0.1:5678" + path)
	if err != nil {
		slog.Error("http get failed", "error", err)
		return
	}
	processResponse(resp)
}

func GetStudentPB(path string) {
	fmt.Println("GET " + path)
	resp, err := http.Get("http://127.0.0.1:5678" + path)
	if err != nil {
		slog.Error("http get failed", "error", err)
		return
	}
	defer resp.Body.Close()
	fmt.Println("响应头：")
	for k, v := range resp.Header {
		fmt.Printf("%s=%s\n", k, v[0])
	}
	if resp.StatusCode != http.StatusOK {
		slog.Error("异常状态码", "http response code", resp.StatusCode)
		io.Copy(os.Stdout, resp.Body)
	} else {
		bs, _ := io.ReadAll(resp.Body)
		var stu idl.Student
		if err := proto.Unmarshal(bs, &stu); err == nil {
			fmt.Printf("pb反序列化成功, name=%s, addr=%s\n", stu.Name, stu.Address)
		}
	}
	os.Stdout.WriteString("\n\n")
}

func PostForm(path string, stu Student) {
	fmt.Println("post form " + path)
	// PostForm()会自动把请求头的Content-Type设为application/x-www-form-urlencoded，并把url.Values转为URL-encoded参数格式放到请求体里
	if resp, err := http.PostForm("http://127.0.0.1:5678"+path, url.Values{"username": {stu.Name}, "addr": {stu.Address}}); err != nil {
		panic(err)
	} else {
		processResponse(resp)
	}
}

func PostJson(path string, stu Student) {
	fmt.Println("post json " + path)
	if bs, err := json.Marshal(stu); err == nil {
		if resp, err := http.Post("http://127.0.0.1:5678"+path, "application/json", bytes.NewReader(bs)); err != nil {
			panic(err)
		} else {
			processResponse(resp)
		}
	} else {
		slog.Error("json marchal failed", "error", err)
	}
}

func PostXml(path string, stu Student) {
	fmt.Println("post xml " + path)
	if bs, err := xml.Marshal(stu); err == nil {
		if resp, err := http.Post("http://127.0.0.1:5678"+path, "application/xml", bytes.NewReader(bs)); err != nil {
			panic(err)
		} else {
			processResponse(resp)
		}
	} else {
		slog.Error("xml marchal failed", "error", err)
	}
}

func PostYaml(path string, stu Student) {
	fmt.Println("post yaml " + path)
	if bs, err := yaml.Marshal(stu); err == nil {
		if resp, err := http.Post("http://127.0.0.1:5678"+path, "application/yaml", bytes.NewReader(bs)); err != nil {
			panic(err)
		} else {
			processResponse(resp)
		}
	} else {
		slog.Error("yaml marchal failed", "error", err)
	}
}

func PostPb(path string, stu Student) {
	fmt.Println("post pb " + path)
	inst := idl.Student{Name: stu.Name, Address: stu.Address}
	if bs, err := proto.Marshal(&inst); err == nil {
		if resp, err := http.Post("http://127.0.0.1:5678"+path, "application/x-protobuf", bytes.NewReader(bs)); err != nil {
			panic(err)
		} else {
			processResponse(resp)
		}
	} else {
		slog.Error("yaml marchal failed", "error", err)
	}
}

func PostAll(path string, stu Student) {
	PostForm(path, stu)
	PostJson(path, stu)
	PostXml(path, stu)
	PostYaml(path, stu)
	PostPb(path, stu)
}

func Request(path, method string, body []byte) {
	request, err := http.NewRequest(method, "http://127.0.0.1:5678"+path, bytes.NewReader(body))
	if err != nil {
		panic(err)
	}
	//  可添加多个Cookie
	request.AddCookie(
		&http.Cookie{
			Name:  "token",
			Value: "ye38ry4928---",
		},
	)
	//所有的cookie都会放到一个http request header中。Cookie: [auth=pass; money=100dollar]
	//设置请求超时
	client := &http.Client{
		Timeout: 500 * time.Millisecond,
	}
	//发起请求
	if resp, err := client.Do(request); err != nil {
		fmt.Println(err)
	} else {
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			io.Copy(os.Stdout, resp.Body)
			return
		}
		fmt.Println("response header:")
		// 其实可以直接通过resp.Cookies()获得*http.Cookie，没必要自己解析
		if values, exists := resp.Header["Set-Cookie"]; exists {
			fmt.Println(values[0])
			cookie, _ := http.ParseSetCookie(values[0])
			fmt.Println("Name:", cookie.Name)
			fmt.Println("Value:", cookie.Value)
			fmt.Println("Domain:", cookie.Domain)
			fmt.Println("MaxAge:", cookie.MaxAge)
			fmt.Println(strings.Repeat("-", 50))
		}
		os.Stdout.WriteString("\n\n")
	}
}

func main_gin() {
	//student := Student{Name: "大乔乔", Address: "北京海淀"}

	// Head("/ping")
	// Get("/home")
	//PostForm("/home", student)

	// Get("/")
	// Get("/1")
	// Get("/2")
	// Get("/3")
	// Get("/4")
	// for i := 0; i < 6; i++ {
	// 	Get("/lt")
	// }

	// Get("/student?name=大乔乔&addr=北京&age=18")
	// Get("/student2?name=大乔乔&addr=北京&age=18")
	// Get("/student/大乔乔/北京")
	// Get("/student/大乔乔/北京/18")
	// Get("/school/北京大学")
	// Get("/location/60.90")
	// Get("/size/3-29")
	// PostForm("/student/form", student)
	// PostJson("/student/json", student)
	//PostPb("/student/pb", student)

	// PostForm("/stu/form", student)
	// PostJson("/stu/json", student)
	// Get("/stu/uri/大乔乔/北京")
	// Get("/stu/uri/大乔乔/北京/海淀")
	// PostXml("/stu/xml", student)
	// PostYaml("/stu/yaml", student)
	//PostAll("/stu/multi_type", student)

	// Get("/user/text")
	// Get("/user/json0")
	// Get("/user/json1")
	// Get("/user/json2")
	// Get("/user/jsonp?callback=yyds") //使用JSONP可以向不同域的服务器请求数据
	// Get("/user/jsonp")               //使用JSONP可以向不同域的服务器请求数据
	// Get("/user/xml")
	// Get("/user/yaml")
	// Get("/user/html")
	// Get("/user/old_page")

	GetStudentPB("/user/pb")

	// ago := url.QueryEscape(time.Now().Add(-86400 * time.Second).Format(time.RFC3339))
	// future := url.QueryEscape(time.Now().Add(86400 * time.Second).Format(time.RFC3339))
	// Get("?name=zcy&score=1&enrollment=" + ago + "&graduation=" + future) //正确
	// Get("?score=1&enrollment=" + ago + "&graduation=" + future)          //name缺失
	// Get("?name=zcy&score=0&enrollment=" + ago + "&graduation=" + future) //score=0
	// Get("?name=zcy&score=1&enrollment=" + future + "&graduation=" + ago) //enrollment晚于今天,graduation早于enrollment

	//Request("/ck", http.MethodGet, nil)
}

// go run ./web/client
