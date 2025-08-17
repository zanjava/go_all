package main

import (
	"fmt"
	"go/frame/mq/rabbitmq"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func Receive(ch *amqp.Channel, qName string, flag int) {
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

	for info := range consumer {
		log.Printf("%d receive message [%s][%s]", flag, info.RoutingKey, info.Body)
		time.Sleep(2 * time.Second)
		info.Ack(false) //通知Server此消息已成功消费。Ack参数为true时，此channel里之前未ack的消息会一并被ack（相当于批量ack）。如果没有ack，则下一次启动时还消费到此消息（除非超时30分钟，因为delivery在30分钟后会被强制ack）,因为channel close时，它里没有ack的消息会再次被放入队列的尾部。
	}
}

func main2() {
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

	//RoundRobin不一定是好的负载均衡方式，因为有的消息消费起来需要更多的时间。
	err = ch.Qos( //quality of service
		2,     //prefetch count。一个消费方最多能有多少条未ack的消息。如果consumer是以noAck（autoAck）启动的则server会忽略该参数。<=0时忽略该参数。该值越小，负载越均衡，但是单个消费方的吞吐也越低。
		0,     //prefetch size。Server端至少攒够这么多字节，才发给消费方。<=0时忽略该参数。
		false, //global
	)
	if err != nil {
		log.Panicf("set Qos  failed: %s", err)
	}

	go Receive(ch, QueueName, 1)
	go Receive(ch, QueueName, 2)
	select {}
}

// go run .\mq\rabbitmq\consumer\
