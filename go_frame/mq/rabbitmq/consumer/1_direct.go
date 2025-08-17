package main

import (
	"fmt"
	"go/frame/mq/rabbitmq"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	QueueName = "MyQueue1"
)

type Consumer = <-chan amqp.Delivery

func CreateConsumer(ch *amqp.Channel, qName string) Consumer {
	var consumer Consumer
	consumer, err := ch.Consume(
		qName, //queue
		"",    //consumer
		false, //auto-ack。autoAck其实就是noAck，只要server把消息传给consumer，本消息就会被标记为ack，而不管它有没有被consumer成功消费。
		false, //exclusive
		false, //no-local
		false, //no-wait
		nil,   //args
	)

	if err != nil {
		log.Panicf("regist consumer failed: %s", err)
	}
	return consumer
}

//var cnt int32 = 0

func ConsumeFromConsumer(consumer Consumer, flag int) {
	for info := range consumer {
		log.Printf("%d receive message [%s][%s]", flag, info.RoutingKey, info.Body)
		// atomic.AddInt32(&cnt, 1)
		// if cnt%2 == 0 {
		// 	info.Ack(true) //每消费10条就确认一次
		// }
		info.Ack(false) //通知Server此消息已成功消费。Ack参数为true时，此channel里之前未ack的消息会一并被ack（相当于批量ack）。如果没有ack，则下一次启动时还消费到此消息（除非超时30分钟，因为delivery在30分钟后会被强制ack）,因为channel close时，它里没有ack的消息会再次被放入队列的尾部。
		// os.Exit(0)          //后面还有一些分配给该consumer的消息，会丢失
	}
}

func main1() {
	//连接RabbitMQ
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/", rabbitmq.User, rabbitmq.Pass, rabbitmq.Host, rabbitmq.Port))
	if err != nil {
		log.Panicf("connect to RabbitMQ failed: %s", err)
	}
	defer conn.Close()

	//创建Channel
	ch, err := conn.Channel()
	if err != nil {
		log.Panicf("open channel failed: %s", err)
	}
	defer ch.Close()

	// 声明队列  producer consumer 看他们谁先创建
	// _, err = ch.QueueDeclare(
	// 	QueueName, //队列名。为空时Server指定一个随机（且唯一）的队列名
	// 	true,      // durable。durable=true即队列数据会持久化入磁盘，这样RabbitMQ Server重启后，队列数据也不会丢失。
	// 	false,     // delete when unused。若设为true，当所有consumer都退出后队列会自动被删除(这意味着Server要监听Consumer进程是否存活)。
	// 	false,     // exclusive。若设为true，队列由声明它的connection独享，其他连接上的channel不能再使用（包含声明、绑定、消费等等）同名的队列。
	// 	false,     // no-wait。
	// 	nil,       // arguments。
	// )
	//一个队列可以对应多个channel（它们平分这个queue里的数据），一个channel可以有多个consumer（它们平分这个channel里的数据）。broker默认会按轮流(RoundRobin)的方式把各个消息发给所有consumer
	ch2, _ := conn.Channel()
	consumer1 := CreateConsumer(ch, QueueName)
	consumer2 := CreateConsumer(ch, QueueName)
	consumer3 := CreateConsumer(ch2, QueueName)

	go ConsumeFromConsumer(consumer1, 1)
	go ConsumeFromConsumer(consumer2, 2)
	go ConsumeFromConsumer(consumer3, 3)

	select {}
}

// go run .\mq\rabbitmq\consumer\
