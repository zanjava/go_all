package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type Student struct {
	Name    string `json:"name"`
	Address string `json:"addr"`
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

func Get(path string) {
	fmt.Println("GET " + path)
	resp, err := http.Get("http://127.0.0.1:5678" + path)
	if err != nil {
		slog.Error("http get failed", "error", err)
		return
	}
	processResponse(resp)
}

func PostForm(path string, stu Student) {
	fmt.Println("post form " + path)
	// PostForm()会自动把请求头的Content-Type设为application/x-www-form-urlencoded，并把url.Values转为URL-encoded参数格式放到请求体里
	if resp, err := http.PostForm("http://127.0.0.1:5678"+path, url.Values{"name": {stu.Name}, "addr": {stu.Address}}); err != nil {
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

func main() {
	student := Student{Name: "大乔乔2", Address: "北京海淀"}
	Get("/student?name=大乔乔&addr=北京&age=18")
	Get("/student/大乔乔/北京")
	Get("/student/大乔乔/北京/18")
	PostForm("/student/form", student)
	PostJson("/student/json", student)
	Request("/cookie", http.MethodGet, nil)
}

// go run ./learn/client
