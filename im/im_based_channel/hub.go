package main

import "sync/atomic"

type Hub struct {
	broadcast  chan []byte          //需要广播的消息
	clients    map[*Client]struct{} //维护所有Server
	register   chan *Client         //Server注册请求通过管道来接收
	unregister chan *Client         //Server注销请求通过管道来接收
	state      int32
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[*Client]struct{}),
		//以下都是无缓冲管道，确保server和hub的交互是同步的，即server的请求hub处理完之后server才执行下一步工作
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan []byte),
		state:      0,
	}
}

func (hub *Hub) Run() {
	if atomic.CompareAndSwapInt32(&hub.state, 0, 1) { //只能Run()一次，Run()不支持并发
		for {
			select {
			case client := <-hub.register:
				hub.clients[client] = struct{}{} //注册server
			case client := <-hub.unregister:
				if _, ok := hub.clients[client]; ok { //防止重复注销，一个channel只能被close一次
					delete(hub.clients, client) //注销server
					close(client.send)          //hub从此以后不需要再向该server广播消息了
				}
			case msg := <-hub.broadcast:
				for client := range hub.clients {
					client.send <- msg
				}
			}
		}
	}
}
