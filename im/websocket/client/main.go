package main

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

func request() {
	// 发起一次http请求（相当于websocket握手），从而创建websocket连接
	conn, resp, err := websocket.DefaultDialer.Dial("ws://127.0.0.1:5678/", nil)
	if err != nil {
		panic(err)
	}

	// 收到ping后，默认情况就会自动返回pong。得调一次ReadMessage()才能激活PingHandler、PongHandler、CloseHandler
	conn.SetPingHandler(func(appData string) error { //SetPingHandler 定义收到ping后做何反应
		log.Println("receive ping message")
		return conn.WriteMessage(websocket.PongMessage, []byte("pong"))
	})
	// 收到CloseMessage后，默认行为就是关闭连接
	conn.SetCloseHandler(func(code int, text string) error {
		log.Println("receive close message")
		return conn.Close()
	})

	// http请求的响应
	log.Println("http响应码", resp.Status) //101 Switching Protocols
	log.Println("http响应头")
	for k, v := range resp.Header {
		log.Printf("%s : %s\n", k, v[0])
	}
	// Upgrade : websocket
	// Connection : Upgrade
	// Sec-Websocket-Accept : Php+5HMrQdbsHZfJ8E9FeBY3lA0=

	//基于websocket连接收发消息
	for i := 0; i < 5; i++ {
		err := conn.WriteMessage(websocket.TextMessage, []byte("你好"))
		if err != nil {
			log.Printf("write error:%v\n", err)
			break
		}
		time.Sleep(10 * time.Second)              //故意停顿一下，验证Server端的SetReadDeadline是否生效
		mtype, message, err := conn.ReadMessage() //如果对方主动断开连接或超时，该行会报错，for循环会退出
		if err != nil {
			log.Printf("read error:%v\n", err)
			break
		} else {
			log.Println("receive message type", mtype) //PingMessage不会走到这边来，会走到PingHandler里面去
			if mtype == websocket.TextMessage {        //ReadMessage()只能读出TextMessage或BinaryMessage
				log.Printf("read TextMessage %s\n", string(message))
			}
		}
	}
}

func main1() {
	request()
}

// go run ./im/websocket/client
