package client

import (
	"log"
	"time"
)

func TcpLongConnection() {
	conn := connect2TcpServer("127.0.0.1:5678")

	for i := 0; i < 3; i++ {
		sendTcpServer(conn)
	}
	conn.Close()
	log.Println("close connection")
}

func UdpLongConnection() {
	conn := connect2UdpServer("127.0.0.1:5678")

	for i := 0; i < 3; i++ {
		sendUdpServer(conn)
	}

	time.Sleep(70 * time.Second)
	sendUdpServer(conn)
	conn.Close()
	log.Println("close connection")
}
