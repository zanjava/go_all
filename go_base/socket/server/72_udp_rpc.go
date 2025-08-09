package server

import (
	"encoding/json"
	transport "go/base/socket"
	"io"
	"log"
	"net"
)

func UdpRpcServer() {
	udpAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:5678")
	transport.CheckError(err)
	conn, err := net.ListenUDP("udp", udpAddr)
	transport.CheckError(err)
	log.Println("return conn")
	defer conn.Close()

	const P = 1000 //server端开1000个并发处理请求
	for i := 0; i < P; i++ {
		go func() {
			for { //处理完一个请求，紧接着处理下一个请求
				request := make([]byte, 256)
				n, remoteAddr, err := conn.ReadFromUDP(request)
				if err != nil && err != io.EOF {
					log.Printf("read error %v", err)
					continue
				}

				response := handle(request[:n])
				if len(response) > 0 {
					conn.WriteToUDP(response, remoteAddr)

				}
			}
		}()
	}
	select {} //作为服务端，永不退出
}

func handle(request []byte) (response []byte) {
	var AddRequest transport.AddRequest
	var err error
	err = json.Unmarshal(request, &AddRequest)
	if err != nil {
		log.Printf("unmarshal request failed: %s", err)
		return nil
	}
	log.Printf("receive request, id %d a %d b %d", AddRequest.RequestId, AddRequest.A, AddRequest.B)
	AddResponse := transport.AddResponse{
		RequestId: AddRequest.RequestId,
		Sum:       AddRequest.A + AddRequest.B,
	}
	response, err = json.Marshal(AddResponse)
	if err != nil {
		log.Printf("marshal response failed: %s", err)
		return nil
	} else {
		log.Printf("send response, id %d sum %d", AddResponse.RequestId, AddResponse.Sum)
	}
	return
}
