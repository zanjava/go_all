package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/client"
)

func main() {
	//Head("/ping")
	clt := client.New()
	//设置json序列化方式
	clt.SetJSONMarshal(sonic.Marshal)
	//设置请求超时
	clt.SetTimeout(time.Minute)

	request := client.AcquireRequest()   //从池子里取出一个Request
	defer client.ReleaseRequest(request) //归还Request
	//设置请求方法
	request.SetMethod(http.MethodPost)
	//设置请求的url
	request.SetURL("http://127.0.0.1:5678/ckp")
	//设置请求头
	request.SetHeader("language", "go")
	request.SetHeader("User-Agent", "client-go")
	//设置Cookie
	request.SetCookie("token", "ye38ry4928")

	//设置请求体
	request.SetFormData("name", "Jerry")
	request.SetJSON(fiber.Map{"name": "zgw", "age": 18})

	// 关联request和client
	request.SetClient(clt)

	// 响应
	resp, err := request.Send()
	if err != nil {
		fmt.Println(err)
	} else {
		if resp.StatusCode() != 200 {
			// 响应码
			fmt.Printf("statusCode=%d\n", resp.StatusCode())
		} else {
			// 响应头
			fmt.Println("响应头")
			for k, v := range resp.Headers() {
				fmt.Printf("%s=%s\n", k, v[0])
			}
			fmt.Println("响应Cookie")
			for _, cookie := range resp.Cookies() {
				fmt.Printf("%s=%s\n", cookie.Key(), cookie.Value())
			}
			// 响应体
			fmt.Println("响应体:", string(resp.Body()))
		}
	}
}
