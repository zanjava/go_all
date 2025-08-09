package client

import (
	transport "go/base/socket"
	"log"
	"net"
	"sync"
	"time"
)

func connect2UdpServer(serverAddr string) net.Conn {
	//跟tcp_client的唯一区别就是这行代码
	conn, err := net.DialTimeout("udp", serverAddr, 3*time.Minute) //一个conn绑定一个本地端口
	transport.CheckError(err)
	log.Printf("establish connection to server %s myself %s\n", conn.RemoteAddr().String(), conn.LocalAddr().String()) //操作系统会随机给客户端分配一个49152~65535上的端口号
	return conn
}

func sendUdpServer(conn net.Conn) {
	n, err := conn.Write([]byte("hello")) //即使Server还启动，建立连接和发送数据都不会返回error，Server启动后也收不到这个数据
	transport.CheckError(err)
	log.Printf("send %d bytes\n", n)
}

func UdpClient() {
	conn := connect2UdpServer("127.0.0.1:5678")

	sendUdpServer(conn)
	conn.Close()
	log.Println("close connection")
}

// Client端，并发使用udp连接
func UdpConnectionCurrent() {
	conn := connect2UdpServer("127.0.0.1:5678")

	wg := sync.WaitGroup{}
	wg.Add(3)
	for i := 0; i < 3; i++ {
		go func() {
			defer wg.Done()
			sendUdpServer(conn)
		}()
	}
	wg.Wait()
	conn.Close()
}
