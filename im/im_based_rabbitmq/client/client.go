package main

import (
	"fmt"
	"go/im/im_based_rabbitmq/common"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
)

var (
	host = "127.0.0.1:5678"
	conn *websocket.Conn
)

func RegistUser() (uid int64) {
	resp, err := http.Get("http://" + host + "/regist_user")
	if err != nil {
		log.Printf("regist user failed: %s", err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		bs, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("read resp body failed: %s", err)
			return
		}
		uid, err = strconv.ParseInt(string(bs), 10, 64)
		if err != nil {
			log.Printf("invalid uid [%s]", string(bs))
			return
		}
		return
	} else {
		log.Printf("resp code %d", resp.StatusCode)
		return
	}
}

func JoinGroup(gid, uid int64) (success bool) {
	resp, err := http.Get(fmt.Sprintf("http://%s/join_group/%d/%d", host, gid, uid))
	if err != nil {
		log.Printf("join group failed: %s", err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		success = true
		return
	} else {
		log.Printf("resp code %d", resp.StatusCode)
		return
	}
}

func Send(msg *common.Message) error {
	if conn == nil {
		var err error
		conn, _, err = websocket.DefaultDialer.Dial(fmt.Sprintf("ws://%s/speak", host), nil)
		if err != nil {
			conn = nil
			log.Printf("connect to websocket server failed: %s", err)
			return err
		} else {
			// 收到ping后，默认情况就会自动返回pong
			conn.SetPingHandler(func(appData string) error {
				// log.Println("receive ping message")
				conn.WriteMessage(websocket.PongMessage, nil)
				return nil
			})
			// 收到CloseMessage后，默认行为就是关闭连接
			conn.SetCloseHandler(func(code int, text string) error {
				// log.Println("receive close message")
				conn.Close()
				return nil
			})
			go conn.ReadMessage() //只需要调一次ReadMessage()，就激活了PingHandler、PongHandler、CloseHandler
		}
	}

	err := conn.WriteJSON(msg)
	if err != nil {
		log.Printf("send message failed: %s", err)
		return err
	} else {
		return nil
	}
}

func ReleaseWebSocket() {
	if conn != nil {
		conn.WriteMessage(websocket.CloseMessage, nil)
		conn.Close()
	}
}

// go run ./im/im_based_rabbitmq/client
