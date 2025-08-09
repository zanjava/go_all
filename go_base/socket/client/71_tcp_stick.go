package client

import (
	"log"
)

// 测试TCP粘包问题
func TcpStick() {
	conn := connect2TcpServer("127.0.0.1:5678")

	for i := 0; i < 13; i++ {
		sendTcp1Server(conn)
	}
	conn.Close()
	log.Println("close connection")
}
