package server

import (
	transport "go/base/socket"
	"log"
	"net"
	"time"
)

func TcpServer() {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", "127.0.0.1:5678")
	transport.CheckError(err)
	listener, err := net.ListenTCP("tcp4", tcpAddr)
	transport.CheckError(err)
	log.Println("waiting for client connection ......")
	conn, err := listener.Accept()
	transport.CheckError(err)
	log.Printf("establish connection to client %s\n", conn.RemoteAddr().String())
	conn.SetReadDeadline(time.Now().Add(10 * time.Second)) //设置一个读的期限，超过这个期限再调Read()就会发生error。默认是60s内可Read()。
	defer conn.Close()                                     //关闭连接

	request := make([]byte, 256)
	n, err := conn.Read(request) //跟读文件类似。可能会因为超时而导致error(之前设置了ReadDeadline)
	transport.CheckError(err)
	log.Printf("receive %s\n", string(request[:n]))
} // TCP需要先启动server，再启动client，否则client连接不上server
