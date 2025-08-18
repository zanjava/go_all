package main

import (
	"fmt"
	"io"
	"time"

	"github.com/gogf/gf/v2/net/gtcp"
	"github.com/gogf/gf/v2/net/gudp"
)

func Udp() {
	const ADDR = "127.0.0.1:5678"
	var serve = func(conn *gudp.ServerConn) {
		const BUFFER_SIZE = 1024
		for {
			data, remoteAddr, err := conn.Recv(BUFFER_SIZE)
			if err == nil {
				fmt.Printf("receive datagrame: [%s] from %s\n", data, remoteAddr.String())
				conn.WriteToUDP(data, remoteAddr)
			} else {
				if err == io.EOF {
					time.Sleep(time.Second)
					continue
				} else {
					fmt.Println(err)
					break
				}
			}
		}
	}
	server := gudp.NewServer(ADDR, serve)
	go server.Run()
	time.Sleep(time.Second)

	conn, _ := gudp.NewClientConn(ADDR)
	defer conn.Close()
	conn.Send([]byte("123"))
	conn.Send([]byte("45678"))
	conn.Send([]byte("9"))
	time.Sleep(time.Second)
}

func Tcp() {
	const ADDR = "127.0.0.1:5678"
	var serve = func(conn *gtcp.Conn) {
		for {
			data, err := conn.RecvPkg()
			if err == nil {
				fmt.Printf("receive datagrame: [%s] \n", data)
			} else {
				if err == io.EOF {
					time.Sleep(time.Second)
					continue
				} else {
					fmt.Println(err)
					break
				}
			}
		}
	}
	server := gtcp.NewServer(ADDR, serve)
	go server.Run()
	time.Sleep(time.Second)

	conn, _ := gtcp.NewConn(ADDR)
	defer conn.Close()
	conn.SendPkg([]byte("123"))
	conn.SendPkg([]byte("45678"))
	conn.SendPkg([]byte("9"))
	time.Sleep(time.Second)
}
