package main

import (
	common "go/frame/mq/rocket"
	"log"
	"sync"

	rmq_client "github.com/apache/rocketmq-clients/golang/v5" //注意：现在是v5
	"github.com/apache/rocketmq-clients/golang/v5/credentials"
)

var (
	producer rmq_client.Producer
	once     sync.Once
	//localTransactionResult = sync.Map{}
)

// 单例模式（考虑并发调用）
// Producer和Topic是多对多的关系
func GetProducer() rmq_client.Producer {
	once.Do(func() {
		if producer == nil {
			var err error
			// 创建Producer
			producer, err = rmq_client.NewProducer(
				&rmq_client.Config{
					Endpoint:    common.Endpoint,
					Credentials: &credentials.SessionCredentials{}, // 即使空的也得给Credentials赋值，否则会发生空指针异常
				},
				rmq_client.WithMaxAttempts(1), // 如果发送失败，再试1次

				//若Rocket服务端未收到发送者提交的二次确认结果，经过固定时间后，服务端将对消息生产者即生产者集群中任一生产者实例发起消息回查
				rmq_client.WithTransactionChecker(&rmq_client.TransactionChecker{
					Check: func(msg *rmq_client.MessageView) rmq_client.TransactionResolution {
						// Producer收到服务端的回查信息后需要决断一下（查一下本地事务的执行结果）：该消息到底应该ROLLBACK还是COMMIT。
						if _, exists := localTransactionResult.Load(msg.GetMessageId()); exists {
							return rmq_client.COMMIT
						} else {
							return rmq_client.ROLLBACK
						}
					},
				}),
			)
			if err != nil {
				log.Fatal(err)
			}
			// 启动Producer
			err = producer.Start()
			if err != nil {
				log.Fatal(err)
			}
		}
	})
	return producer
}
