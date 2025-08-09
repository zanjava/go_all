package server

import (
	transport "go/base/socket"
	"log"
	"net"
	"time"
)

func TcpLongConnection() {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", "127.0.0.1:5678")
	transport.CheckError(err)
	listener, err := net.ListenTCP("tcp4", tcpAddr)
	transport.CheckError(err)
	log.Println("waiting for client connection ......")
	conn, err := listener.Accept()
	transport.CheckError(err)
	log.Printf("establish connection to client %s\n", conn.RemoteAddr().String())
	defer conn.Close()

	// time.Sleep(5 * time.Second) //故意多sleep一会儿，让client多发几条消息过来
	request := make([]byte, 256)
	for { //长连接
		conn.SetReadDeadline(time.Now().Add(2 * time.Minute)) //每次都要续命
		n, err := conn.Read(request)                          //如果流里没数据，Read()会阻塞。对方close后，这里会遇到EOF
		transport.CheckError(err)
		log.Printf("receive %s\n", string(request[:n]))
	}
}

func UdpLongConnection() {
	udpAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:5678")
	transport.CheckError(err)
	conn, err := net.ListenUDP("udp", udpAddr)
	transport.CheckError(err)
	log.Println("return conn")
	defer conn.Close()

	time.Sleep(5 * time.Second) //故意多sleep一会儿，让client多发几条消息过来
	request := make([]byte, 256)
	for { //长连接
		conn.SetReadDeadline(time.Now().Add(2 * time.Minute)) //每次都要续命
		n, remoteAddr, err := conn.ReadFromUDP(request)       //对方close后，这里不会有error。但是2分钟之后如果没有数据到来，还是会发生timeout error
		transport.CheckError(err)
		log.Printf("receive request %s from %s\n", string(request[:n]), remoteAddr.String())
	}
}
