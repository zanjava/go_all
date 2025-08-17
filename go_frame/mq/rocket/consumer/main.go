package main

import (
	"context"
	"errors"
	"fmt"
	common "go/frame/mq/rocket"
	"log"
	"os"
	"time"

	rmq_client "github.com/apache/rocketmq-clients/golang/v5"
	"github.com/apache/rocketmq-clients/golang/v5/credentials"
)

var (
	receiveWaitDuration = 5000 * time.Millisecond
	invisibleDuration   = 12000 * time.Millisecond // 最小为10000ms。SimpleConsumer通过Receive拿到消息后，如果在invisibleDuration时间内没有ACK，RocketMQ会尝试再次投递
)

func consume() {
	//日志默认输出到根目录~/logs/rocketmqlogs下
	// os.Setenv(rmq_client.ENABLE_CONSOLE_APPENDER, "true") //日志输出到终端
	//日志输出到指定目录的指定文件里
	os.Setenv(rmq_client.CLIENT_LOG_ROOT, "D:/go_all/go_frame/mq/rocket/log")
	os.Setenv(rmq_client.CLIENT_LOG_FILENAME, "rocket_consumer.log")
	rmq_client.ResetLogger() //ZapLog

	// 创建Consumer。Consumer和Topic是多对多的关系。
	// SimpleConsumer 是一种接口原子型的消费者类型，消息的获取、消费状态提交以及消费重试都是通过消费者业务逻辑主动发起调用完成。
	simpleConsumer, err := rmq_client.NewSimpleConsumer(
		&rmq_client.Config{
			Endpoint:      common.Endpoint,                   //proxy地址
			ConsumerGroup: "recommend_biz",                   //归属于哪个组
			Credentials:   &credentials.SessionCredentials{}, // 即使空的也得给Credentials赋值，否则会发生空指针异常
		},
		rmq_client.WithAwaitDuration(receiveWaitDuration), //最多等多长时间，超时后Receive函数会返回一个error：CODE=40401
		rmq_client.WithSubscriptionExpressions(map[string]*rmq_client.FilterExpression{ //可以订阅多个主题，所以是map
			common.NORMAL_TOPIC:      rmq_client.NewFilterExpression("t1||t2"), //订阅该主题下的特定Tag
			common.DELAY_TOPIC:       rmq_client.SUB_ALL,                       //订阅该主题下的所有Tag
			common.FIFO_TOPIC:        rmq_client.SUB_ALL,
			common.TRANSACTION_TOPIC: rmq_client.SUB_ALL,
		}),
	)
	if err != nil {
		log.Fatal(err)
	}
	// 启动Consumer
	err = simpleConsumer.Start()
	if err != nil {
		log.Fatal(err)
	}
	// 终止Consumer
	defer simpleConsumer.GracefulStop()

	// 开始消费，这是个无限循环
	ctx := context.Background()
	log.Println("start receive message")
	for {
		megs, err := simpleConsumer.Receive(ctx,
			3,                 //一次接收几条消息。如果还没攒够会阻塞，直到超过AwaitDuration超时返回error。RocketMQ默认的业务场景就是大数据，默认的处理方式就是批量
			invisibleDuration, // SimpleConsumer通过Receive拿到消息后，如果在invisibleDuration时间内没有ACK，RocketMQ会尝试再次投递
		)
		if err != nil {
			var e *rmq_client.ErrRpcStatus
			if errors.As(err, &e) {
				if e.Code != 40401 { // no new message
					log.Printf("error code: %d, error message: %s\n", e.Code, e.Message)
				}
			}
			continue
		}
		for _, mg := range megs {
			// 集群内每条消息的ID全局唯一
			text := fmt.Sprintf("topic:%s, body:%s, tag:%s, keys:%v, offset:%d, msgId:%s", mg.GetTopic(), string(mg.GetBody()), *mg.GetTag(), mg.GetKeys(), mg.GetOffset(), mg.GetMessageId())
			// if rand.Int()%2 == 0 { //如果消费成功
			log.Println(text)
			simpleConsumer.Ack(ctx, mg) // ACK，告诉服务端消费成功了
			// } else { //如果消费成失败
			// 	log.Println(text + " NOT ACK") // invisibleDuration之后，没有ACK的消息，服务端会重新投递，这样消费的顺序性就被扰乱了
			// }

		}
	}
}

func main() {
	//模拟同一个Group里的多个consumer并行消费
	// const C = 1
	// wg := sync.WaitGroup{}
	// wg.Add(C)
	// for i := 0; i < C; i++ {
	// 	defer wg.Done()
	// 	consume()
	// }
	// wg.Done()
	consume()
}

// go run ./mq/rocket/consumer
