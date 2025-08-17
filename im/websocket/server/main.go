package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		HandshakeTimeout: 1 * time.Second,
		ReadBufferSize:   100,
		WriteBufferSize:  100,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func handler(w http.ResponseWriter, r *http.Request) {
	//针对本次请求，创建一个单独的websocket连接
	conn, err := upgrader.Upgrade(w, r, nil) //http升级为websocket协议
	if err != nil {
		log.Printf("upgrade error: %v\n", err)
		return
	}

	conn.SetPongHandler(func(appData string) error {
		log.Println("receive pong")
		return nil
	})

	conn.SetCloseHandler(func(code int, text string) error {
		conn.Close()
		return nil
	})

	for i := 0; i < 5; i++ {
		mtype, message, err := conn.ReadMessage() //连接里没有数据时会阻塞。如果对方主动断开连接或超时，该行会报错，for循环会退出
		if err != nil {
			log.Printf("read error:%v\n", err)
			break
		} else {
			if mtype == websocket.TextMessage {
				log.Printf("read %s\n", string(message))
			} else if mtype == websocket.PongMessage { //ReadMessage()只能读出TextMessage或BinaryMessage，不可能读出PongMessage
				log.Println("read pong")
			}
		}

		err = conn.WriteMessage(websocket.TextMessage, []byte("什么事"))
		if err != nil {
			log.Printf("write error:%v\n", err)
			break
		}

		err = conn.WriteMessage(websocket.PingMessage, []byte("ping"))
		if err != nil {
			log.Printf("send ping error:%v\n", err)
			break
		} else {
			log.Println("send ping")
		}
		deadline := time.Now().Add(5 * time.Second)
		conn.SetReadDeadline(deadline) // 下一次调用ReadMessage如果超出了5秒，会返回error
		log.Printf("must read from websocket connection before %s\n", deadline.Format("2006-01-02 15:04:05"))
		// conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
	}
	// conn.Close()
	conn.WriteMessage(websocket.CloseMessage, nil) //关闭连接
	log.Println("send close message")
}

func main1() {
	//启动http server
	http.HandleFunc("/", handler)
	if err := http.ListenAndServe("127.0.0.1:5678", nil); err != nil {
		log.Printf("start http service error: %s\n", err)
	}
}

// go run ./im/websocket/server
