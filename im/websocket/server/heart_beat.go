package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var (
	pongWait   = 5 * time.Second //等待pong的超时时间
	pingPeriod = 3 * time.Second //发送ping的周期，必须短于pongWait
)

func heartBeat(conn *websocket.Conn) {
	conn.SetPongHandler(func(appData string) error {
		log.Println("receive pong")
		deadline := time.Now().Add(pongWait)
		conn.SetReadDeadline(deadline)
		log.Printf("must read before %s", deadline.Format("2006-01-02 15:04:05"))
		return nil
	})

	// err := conn.WriteMessage(websocket.PingMessage, nil)
	// if err != nil {
	// 	log.Printf("write ping error:%v\n", err)
	// 	conn.WriteMessage(websocket.CloseMessage, nil)
	// }

	ticker := time.NewTicker(pingPeriod)
LOOP:
	for {
		select { //通过select确保每次Ping的间隔是准确的2秒
		case <-ticker.C:
			err := conn.WriteMessage(websocket.PingMessage, nil)
			if err != nil {
				log.Printf("write ping error:%v\n", err)
				conn.WriteMessage(websocket.CloseMessage, nil)
				break LOOP
			}
			log.Println("write ping")
		}
	}
}

func hdl(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil) //http升级为websocket协议
	if err != nil {
		log.Printf("upgrade error: %v\n", err)
		return
	}
	//结束时关闭websocket连接
	defer func() {
		log.Println("close websocket connection")
		conn.WriteMessage(websocket.CloseMessage, nil)
	}()

	go heartBeat(conn) //心跳保持

	for {

		_, body, err := conn.ReadMessage()
		if err != nil {
			log.Printf("read error:%v\n", err)
			break
		} else {
			log.Printf("read %s", string(body))
		}
	}
}

func main() {
	//启动http server
	http.HandleFunc("/", hdl)
	if err := http.ListenAndServe("127.0.0.1:5678", nil); err != nil {
		log.Printf("start http service error: %s\n", err)
	}
}

// go run ./im/websocket/server
