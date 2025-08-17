package main

import (
	"log"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	conn, _, err := websocket.DefaultDialer.Dial("ws://127.0.0.1:5678", nil)
	if err != nil {
		conn = nil
		log.Printf("connect to websocket server failed: %s", err)
		return
	} else {
		// 收到ping后，默认情况就会自动返回pong
		conn.SetPingHandler(func(appData string) error {
			log.Println("receive ping message")
			conn.WriteMessage(websocket.PongMessage, nil)
			log.Println("write pong message")
			return nil
		})
		// 收到CloseMessage后，默认行为就是关闭连接
		conn.SetCloseHandler(func(code int, text string) error {
			log.Println("receive close message")
			conn.Close()
			os.Exit(0)
			return nil
		})
		go conn.ReadMessage() //只需要调一次ReadMessage()，就激活了PingHandler、PongHandler、CloseHandler。具体调用关系是：ReadMessage() -> NextReader() -> advanceFrame()(循环调用) -> 从websocket连接里读数据，根据frameType调用PingHandler、PongHandler或CloseHandler。由于NextReader()只会返回TextMessage或BinaryMessage这2种类型的消息，所以ReadMessage()也只会返回TextMessage或BinaryMessage这2种类型的消息
	}
	for {
		conn.WriteMessage(websocket.TextMessage, []byte("你好"))
		log.Println("write 你好")
		time.Sleep(10 * time.Second)
	}
}

// go run ./im/websocket/client
