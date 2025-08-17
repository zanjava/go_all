package main

import (
	"context"
	common "go/frame/mq/rocket"
	"log"
	"math/rand/v2"
	"strconv"
	"sync"
	"time"

	rmq_client "github.com/apache/rocketmq-clients/golang/v5" //注意：现在是v5
)

var (
	localTransactionResult = sync.Map{}
)

// RocketMQ事务机制要实现的目标：执行本地事务和发消息这2件事要绑定在一起，如果其中1件事失败，另1件就不做或回滚
func Transaction() {
	producer := GetProducer()

	// 连发多条消息
	for i := 0; i < 10; i++ {
		// 事务开始
		transaction := producer.BeginTransaction()

		// 发送第1条消息
		// 创建一条消息
		msg := &rmq_client.Message{
			Topic: common.TRANSACTION_TOPIC,          // 主题
			Body:  []byte("鱼戏莲叶东" + strconv.Itoa(i)), // 消息体
		}
		msg.SetKeys("k1", "k2") // key用于快速检索到消息
		// 消费者订阅消息时按Tag进行过滤
		msg.SetTag("t1") //一条消息只能有一个Tag

		// 第一次提交消息，Rocket服务端将该消息标记为"暂不能投递"，该消息称为“半消息”
		resp, err := producer.SendWithTransaction(context.Background(), msg, transaction)
		if err != nil {
			log.Printf("senf half message failed: %s", err)
			continue //不执行本地事务
		}
		var msgID string
		SendReceipt := resp[0]
		log.Printf("MessageID: %s, TransactionId: %s, Offset: %d, Endpoints: %s\n", SendReceipt.MessageID, SendReceipt.TransactionId, SendReceipt.Offset, SendReceipt.Endpoints.String())
		msgID = SendReceipt.MessageID

		time.Sleep(1 * time.Second) // 模拟执行本地事务（比如修改数据库）

		if rand.Float32() < 0.5 { //如果本地事务执行成功
			localTransactionResult.Store(msgID, true) //存储本地事务的执行结果
			err = transaction.Commit()                //通知服务端，之前的半消息可以投递了
			if err != nil {
				log.Fatal(err)
			}
		} else { //如果本地事务执行失败
			err = transaction.RollBack() //通知服务端，把之前的半消息给删掉
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
