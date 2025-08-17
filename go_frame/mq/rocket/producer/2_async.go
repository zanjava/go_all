package main

import (
	"context"
	common "go/frame/mq/rocket"
	"log"
	"math/rand/v2"
	"strconv"
	"time"

	rmq_client "github.com/apache/rocketmq-clients/golang/v5" //注意：现在是v5
)

func Async() {
	producer := GetProducer()

	// 连发多条消息
	for i := 0; i < 10; i++ {
		// 创建一条消息
		msg := &rmq_client.Message{
			Topic: common.NORMAL_TOPIC,               // 主题
			Body:  []byte("鱼戏莲叶东" + strconv.Itoa(i)), // 消息体
		}
		msg.SetKeys("k1", "k2") // key用于快速检索到消息
		// 消费者订阅消息时按Tag进行过滤
		if rand.Int()%2 == 0 {
			msg.SetTag("t1") //一条消息只能有一个Tag
		} else {
			msg.SetTag("t2")
		}

		//异步发送
		producer.SendAsync(context.Background(), msg, func(ctx context.Context, resp []*rmq_client.SendReceipt, err error) {
			if err != nil {
				log.Fatal(err)
			}
			for i := 0; i < len(resp); i++ {
				SendReceipt := resp[i]
				log.Printf("MessageID: %s, TransactionId: %s, Offset: %d, Endpoints: %s\n", SendReceipt.MessageID, SendReceipt.TransactionId, SendReceipt.Offset, SendReceipt.Endpoints.String())
			}
		})
	}
	time.Sleep(3 * time.Second)
}
