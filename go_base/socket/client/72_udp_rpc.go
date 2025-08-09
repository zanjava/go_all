package client

import (
	"encoding/json"
	transport "go/base/socket"
	"log"
	"math/rand/v2"
	"sync"
)

func UdpRpcClient() {
	const P = 500 // 模拟500个client
	const C = 10  //每个client发起10次请求，然后关闭连接
	wg := sync.WaitGroup{}
	wg.Add(P)
	for i := 0; i < P; i++ {
		go func() {
			defer wg.Done()
			conn := connect2UdpServer("127.0.0.1:5678")
			for j := 0; j < C; j++ {
				request := transport.AddRequest{
					RequestId: rand.Int(),
					A:         int(rand.Int32()) % 100,
					B:         int(rand.Int32()) % 100,
				}
				bs, err := json.Marshal(request)
				if err != nil {
					log.Printf("marshal request failed: %s", err)
					continue
				}
				if _, err := conn.Write(bs); err == nil {
					log.Printf("send request, id %d a %d b %d", request.RequestId, request.A, request.B)
				}
			}

			buffer := make([]byte, 256)
			for j := 0; j < C; j++ {
				if n, err := conn.Read(buffer); err == nil {
					var response transport.AddResponse
					err = json.Unmarshal(buffer[:n], &response)
					if err == nil {
						log.Printf("receive response, id %d sum %d", response.RequestId, response.Sum)
					} else {
						log.Printf("unmarshal response failed: %s", err)
					}
				} else {
					log.Printf("read response failed: %s", err)
				}
			}
			conn.Close() //每个client发起10次请求，然后关闭连接
		}()
	}
	wg.Wait()
}
