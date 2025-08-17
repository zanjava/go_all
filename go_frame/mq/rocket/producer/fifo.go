package main

import (
	"context"
	common "go/frame/mq/rocket"
	"log"
	"strconv"

	rmq_client "github.com/apache/rocketmq-clients/golang/v5" //注意：现在是v5
)

func Fifo() {
	producer := GetProducer()

	// 连发多条消息
	for i := 0; i < 10; i++ {
		// 创建一条消息
		msg := &rmq_client.Message{
			Topic: common.FIFO_TOPIC,                 // 主题
			Body:  []byte("鱼戏莲叶东" + strconv.Itoa(i)), // 消息体
		}
		msg.SetKeys("k1", "k2") // key用于快速检索到消息
		// 消费者订阅消息时按Tag进行过滤
		msg.SetTag("t1") //一条消息只能有一个Tag

		// 创建Topic时必须指定消息类型为FIFO，才能调用SetMessageGroup()函数
		msg.SetMessageGroup("g1") //设置消息的Group，同一个组的消息会被放到同一个队列中，从而实现存储的顺序性。该消息在broker里会被打上“__SHARDINGKEY=g1”

		resp, err := producer.Send(context.Background(), msg) // 同步发送
		if err != nil {
			log.Fatal(err)
		}
		for i := 0; i < len(resp); i++ {
			SendReceipt := resp[i]
			log.Printf("MessageID: %s, TransactionId: %s, Offset: %d, Endpoints: %s\n", SendReceipt.MessageID, SendReceipt.TransactionId, SendReceipt.Offset, SendReceipt.Endpoints.String())
		}
	}
}
