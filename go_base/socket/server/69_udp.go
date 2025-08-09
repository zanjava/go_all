package server

import (
	transport "go/base/socket"
	"log"
	"net"
	"sync"
	"time"
)

func UdpServer() {
	udpAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:5678")
	transport.CheckError(err)
	conn, err := net.ListenUDP("udp", udpAddr) //UDP不需要创建连接，所以不需要像TCP那样通过Accept()创建连接，这里的conn是个假连接，不需要阻塞
	transport.CheckError(err)
	log.Println("return conn")
	conn.SetReadDeadline(time.Now().Add(30 * time.Second)) //超时到来之前，client必须发来数据
	defer conn.Close()

	request := make([]byte, 256)
	n, remoteAddr, err := conn.ReadFromUDP(request) //ReadFromUDP会返回remoteAddr。由于多人client共享同一个conn，所以server需要知道这个数据是从哪个client发过来的。可能会因为超时而导致error(之前设置了ReadDeadline)
	transport.CheckError(err)
	log.Printf("receive request %s from %s\n", string(request[:n]), remoteAddr.String())
	conn.WriteToUDP([]byte("hello client"), remoteAddr) //回复数据给client
}

// Server端，并发使用udp连接
func UdpConnectionCurrent() {
	udpAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:5678")
	transport.CheckError(err)
	conn, err := net.ListenUDP("udp", udpAddr)
	transport.CheckError(err)
	log.Println("return conn")
	defer conn.Close()

	wg := sync.WaitGroup{}
	wg.Add(3)
	for i := 0; i < 3; i++ {
		go func() {
			defer wg.Done()
			request := make([]byte, 256)
			for { //长连接
				conn.SetReadDeadline(time.Now().Add(2 * time.Minute)) //每次都要续命
				n, remoteAddr, err := conn.ReadFromUDP(request)       //对方close后，这里不会有error。但是2分钟之后如果没有数据到来，还是会发生timeout error
				transport.CheckError(err)
				log.Printf("%d receive request %s from %s\n", i, string(request[:n]), remoteAddr.String())
			}
		}()
	}
	wg.Wait()
}
